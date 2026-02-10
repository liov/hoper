package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hopeio/gox/database/redis/hash"
	"github.com/hopeio/gox/encoding/json"
	"github.com/hopeio/scaffold/errcode"
	"go.uber.org/zap"

	"github.com/hopeio/gox/log"
	"github.com/liov/hoper/server/go/protobuf/common"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/global"
	modelconst "github.com/liov/hoper/server/go/user/model"
)

type UserDao struct {
	*redis.Client
}

func GetUserDao(ctx context.Context, client *redis.Client) *UserDao {
	return &UserDao{client.WithContext(ctx)}
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserToRedis(user *model.AuthInfo) error {
	ctx := d.Context()

	UserString, err := json.MarshalToString(user)
	if err != nil {
		log.Errorw("UserToRedis.MarshalToString", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if redisErr := d.SetEX(ctx, loginUserKey, UserString, global.Conf.User.TokenMaxAge).Err(); redisErr != nil {
		log.Errorw("UserToRedis.SetEX", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserFromRedis(userId uint64) (*model.AuthInfo, error) {
	ctx := d.Context()

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userId, 10)

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

func (d *UserDao) EditRedisUser(user *model.AuthInfo) error {
	ctx := d.Context()

	UserString, err := json.MarshalToString(user)
	if err != nil {
		log.Errorw("UserToRedis.MarshalToString", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	err = d.Client.Set(ctx, loginUserKey, UserString, 0).Err()
	if err != nil {
		log.Errorw("UserToRedis.SetEX", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserHashToRedis(user *model.AuthInfo) error {
	ctx := d.Context()

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, loginUserKey, hash.Marshal(user)...)
		pipe.Expire(ctx, loginUserKey, global.Conf.User.TokenMaxAge)
		return nil
	}); err != nil {
		log.Errorw("UserHashToRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserHashFromRedis(user *model.AuthInfo) error {
	ctx := d.Context()

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	userArgs, err := d.HGetAll(ctx, loginUser).Result()
	if err != nil {
		log.Errorw("UserHashFromRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	log.Debug(userArgs)
	if len(userArgs) == 0 {
		return model.UserErrInvalidToken
	}
	hash.Unmarshal(user, userArgs)
	return nil
}

func (d *UserDao) EfficientUserHashToRedis(user *model.AuthInfo) error {
	ctx := d.Context()
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
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
func (d *UserDao) EfficientUserHashFromRedis(user *model.AuthInfo) error {
	ctx := d.Context()

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	userArgs, err := d.HGetAll(ctx, loginUser).Result()
	log.Debug(userArgs)
	if err != nil {
		log.Errorw("EfficientUserHashFromRedis", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	if len(userArgs) == 0 {
		return model.UserErrLoginTimeout
	}
	user.Name = userArgs["Name"]
	user.Name = userArgs["Name"]
	n, err := strconv.ParseUint(userArgs["Role"], 10, 32)
	user.Role = model.Role(n)
	n, err = strconv.ParseUint(userArgs["Status"], 10, 8)
	return nil
}

func (d *UserDao) UserLastActiveTime(userId uint64) error {
	ctx := d.Context()
	userIdStr := strconv.FormatUint(userId, 10)
	loginUser := modelconst.LoginUserKey + userIdStr
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, common.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, modelconst.LoginUserKey+"ActiveTime",
			&redis.Z{Score: float64(time.Now().UnixMilli()), Member: userIdStr})
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); err != nil {
		log.Errorw("UserLastActiveTime", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

func (d *UserDao) RedisUserInfoEdit(field string, user *model.AuthInfo) error {
	ctx := d.Context()

	key := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	err := d.HSet(ctx, key, field, user).Err()
	if err != nil {
		log.Errorw("RedisUserInfoEdit", zap.Error(err))
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

func (d *UserDao) GetUserExtRedis(userId uint64) (*model.UserExt, error) {
	ctx := d.Context()

	key := modelconst.UserExtKey + strconv.FormatUint(userId, 10)

	userExt, err := d.HGetAll(ctx, key).Result()
	if err != nil {
		log.Errorw("GetUserExtRedis", zap.Error(err))
		return nil, errcode.RedisErr.Wrap(err)
	}
	if len(userExt) > 3 {
		followCount, _ := strconv.ParseUint(userExt["Follow"], 10, 64)
		followedCount, _ := strconv.ParseUint(userExt["Followed"], 10, 64)
		return &model.UserExt{
			Follow:   followCount,
			Followed: followedCount,
		}, nil
	}
	return nil, nil
}
