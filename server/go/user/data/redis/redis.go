package redis

import (
	"context"
	"strconv"
	"time"

	redisx "github.com/hopeio/gox/database/redis"
	"github.com/hopeio/gox/encoding/json"
	"github.com/hopeio/scaffold/errcode"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/hopeio/gox/log"

	"github.com/liov/hoper/server/go/user/model"
	"github.com/liov/hoper/server/go/global"
	userpb "github.com/liov/hoper/server/go/protobuf/user"
)

type UserDao struct {
	*redis.Client
}

func GetUserDao(client *redis.Client) *UserDao {
	return &UserDao{client}
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserToRedis(ctx context.Context, user *model.AuthInfo) error {

	UserString, err := json.MarshalToString(user)
	if err != nil {
		log.Errorw("UserToRedis.MarshalToString", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}

	loginUserKey := model.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if redisErr := d.Set(ctx, loginUserKey, UserString, global.Conf.User.TokenMaxAge).Err(); redisErr != nil {
		log.Errorw("UserToRedis.Set", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserFromRedis(ctx context.Context, userId uint64) (*model.AuthInfo, error) {

	loginUser := model.LoginUserKey + strconv.FormatUint(userId, 10)

	userString, err := d.Get(ctx, loginUser).Result()
	if err != nil {
		log.Errorw("UserFromRedis.Get", zap.Error(err))
		return nil, errcode.RedisErr.Wrap(err)
	}

	var user model.AuthInfo
	err = json.UnmarshalFromString(userString, &user)
	if err != nil {
		log.Errorw("UserFromRedis.UnmarshalFromString", zap.Error(err))
		return nil, errcode.RedisErr.Wrap(err)
	}
	return &user, nil
}

func (d *UserDao) EditRedisUser(ctx context.Context, user *model.AuthInfo) error {

	UserString, err := json.MarshalToString(user)
	if err != nil {
		log.Errorw("UserToRedis.MarshalToString", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	loginUserKey := model.LoginUserKey + strconv.FormatUint(user.Id, 10)
	err = d.Client.Set(ctx, loginUserKey, UserString, 0).Err()
	if err != nil {
		log.Errorw("UserToRedis.Set", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserHashToRedis(ctx context.Context, user *model.AuthInfo) error {

	loginUserKey := model.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, loginUserKey, redisx.HashEncode(user)...)
		pipe.Expire(ctx, loginUserKey, global.Conf.User.TokenMaxAge)
		return nil
	}); err != nil {
		log.Errorw("UserHashToRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserHashFromRedis(ctx context.Context, user *model.AuthInfo) error {

	loginUser := model.LoginUserKey + strconv.FormatUint(user.Id, 10)

	userArgs, err := d.HGetAll(ctx, loginUser).Result()
	if err != nil {
		log.Errorw("UserHashFromRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	log.Debug(userArgs)
	if len(userArgs) == 0 {
		return userpb.UserErrInvalidToken
	}
	redisx.HashDecode(user, userArgs)
	return nil
}

func (d *UserDao) EfficientUserHashToRedis(ctx context.Context, user *model.AuthInfo) error {
	loginUserKey := model.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, loginUserKey, "Name", user.Name,
			"Role", uint32(user.Role),
			"LastActiveAt", time.Now().UnixMilli())
		pipe.Expire(ctx, loginUserKey, global.Conf.User.TokenMaxAge)
		return nil
	}); err != nil {
		log.Errorw("EfficientUserHashToRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

/*
创建空白哈希表时， 程序默认使用 REDIS_ENCODING_ZIPLIST 编码， 当以下任何一个条件被满足时， 程序将编码从 REDIS_ENCODING_ZIPLIST 切换为 REDIS_ENCODING_HT ：

哈希表中某个键或某个值的长度大于 server.hash_max_ziplist_value （默认值为 64 ）。
压缩列表中的节点数量大于 server.hash_max_ziplist_entries （默认值为 512 ）。
*/
func (d *UserDao) EfficientUserHashFromRedis(ctx context.Context, user *model.AuthInfo) error {

	loginUser := model.LoginUserKey + strconv.FormatUint(user.Id, 10)

	userArgs, err := d.HGetAll(ctx, loginUser).Result()
	log.Debug(userArgs)
	if err != nil {
		log.Errorw("EfficientUserHashFromRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	if len(userArgs) == 0 {
		return userpb.UserErrLoginTimeout
	}
	user.Name = userArgs["Name"]
	user.Name = userArgs["Name"]
	n, err := strconv.ParseUint(userArgs["Role"], 10, 32)
	user.Role = userpb.Role(n)
	n, err = strconv.ParseUint(userArgs["Status"], 10, 8)
	return nil
}

func (d *UserDao) UserLastActiveTime(ctx context.Context, userId uint64) error {
	userIdStr := strconv.FormatUint(userId, 10)
	loginUser := model.LoginUserKey + userIdStr
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, global.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, model.LoginUserKey+"ActiveTime",
			redis.Z{Score: float64(time.Now().UnixMilli()), Member: userIdStr})
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); err != nil {
		log.Errorw("UserLastActiveTime", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

func (d *UserDao) RedisUserInfoEdit(ctx context.Context, field string, user *model.AuthInfo) error {

	key := model.LoginUserKey + strconv.FormatUint(user.Id, 10)

	err := d.HSet(ctx, key, field, user).Err()
	if err != nil {
		log.Errorw("RedisUserInfoEdit", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

func (d *UserDao) GetUserExtRedis(ctx context.Context, userId uint64) (*userpb.UserExt, error) {

	key := model.UserExtKey + strconv.FormatUint(userId, 10)

	userExt, err := d.HGetAll(ctx, key).Result()
	if err != nil {
		log.Errorw("GetUserExtRedis", zap.Error(err))
		return nil, errcode.RedisErr.Wrap(err)
	}
	if len(userExt) > 3 {
		followCount, _ := strconv.ParseUint(userExt["Follow"], 10, 64)
		followedCount, _ := strconv.ParseUint(userExt["Followed"], 10, 64)
		return &userpb.UserExt{
			Follow:   followCount,
			Followed: followedCount,
		}, nil
	}
	return nil, nil
}
