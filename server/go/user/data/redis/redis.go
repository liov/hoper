package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/utils/encoding/json"
	"strconv"

	"github.com/hopeio/cherry/protobuf/errcode"
	redisi "github.com/hopeio/cherry/utils/dao/redis"
	"github.com/hopeio/cherry/utils/encoding/hash"
	"github.com/hopeio/cherry/utils/log"
	"github.com/liov/hoper/server/go/protobuf/common"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/confdao"
	modelconst "github.com/liov/hoper/server/go/user/model"
)

type UserDao struct {
	*httpctx.Context
	*redis.Client
}

func GetUserDao(ctx *httpctx.Context, client *redis.Client) *UserDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &UserDao{ctx, client}
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserToRedis() error {
	ctxi := d
	ctx := ctxi.BaseContext()
	UserString, err := json.MarshalToString(ctxi.AuthInfo)
	if err != nil {
		return d.ErrorLog(errcode.RedisErr, err, "UserToRedis.MarshalToString")
	}

	loginUserKey := modelconst.LoginUserKey + ctxi.AuthID
	if redisErr := d.SetEX(ctx, loginUserKey, UserString, confdao.Conf.Customize.TokenMaxAge).Err(); redisErr != nil {
		return d.ErrorLog(errcode.RedisErr, err, "UserToRedis.SetEX")
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserFromRedis() (*model.AuthInfo, error) {
	ctxi := d
	ctx := ctxi.BaseContext()
	loginUser := modelconst.LoginUserKey + ctxi.AuthID

	userString, err := redisi.String(d.Get(ctx, loginUser).Result())
	if err != nil {
		return nil, d.ErrorLog(errcode.RedisErr, err, "UserFromRedis.Get")
	}

	var user model.AuthInfo
	err = json.UnmarshalFromString(userString, &user)
	if err != nil {
		return nil, d.ErrorLog(errcode.RedisErr, err, "UserFromRedis.UnmarshalFromString")
	}
	return &user, nil
}

func (d *UserDao) EditRedisUser() error {
	ctx := d.BaseContext()
	UserString, err := json.MarshalToString(d.AuthInfo)
	if err != nil {
		return d.ErrorLog(errcode.RedisErr, err, "EditRedisUser.MarshalToString")
	}
	loginUserKey := modelconst.LoginUserKey + d.AuthID
	err = d.Set(ctx, loginUserKey, UserString, 0).Err()
	if err != nil {
		return d.ErrorLog(errcode.RedisErr, err, "EditRedisUser.MarshalToString")
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func (d *UserDao) UserHashToRedis() error {
	ctxi := d
	ctx := d.BaseContext()
	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + ctxi.AuthID
	redisArgs = append(redisArgs, redisi.CommandHMSET, loginUserKey)
	redisArgs = append(redisArgs, hash.Marshal(ctxi.AuthInfo)...)
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Do(ctx, redisArgs...)
		pipe.Expire(ctx, loginUserKey, confdao.Conf.Customize.TokenMaxAge)
		return nil
	}); err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "UserHashToRedis")
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (d *UserDao) UserHashFromRedis() error {
	ctxi := d
	ctx := ctxi.BaseContext()
	loginUser := modelconst.LoginUserKey + ctxi.AuthID

	userArgs, err := redisi.Strings(d.Do(ctx, redisi.CommandHGETALL, loginUser).Result())
	if err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, redisi.CommandHGETALL)
	}
	log.Debug(userArgs)
	if len(userArgs) == 0 {
		return model.UserErrInvalidToken
	}
	hash.Unmarshal(ctxi.AuthInfo, userArgs)
	return nil
}

func (d *UserDao) EfficientUserHashToRedis() error {
	ctxi := d
	ctx := ctxi.Context
	user := ctxi.AuthInfo.(*model.AuthInfo)
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, err := d.Pipelined(ctx.BaseContext(), func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx.BaseContext(), loginUserKey, "Name", user.Name,
			"Role", uint32(user.Role),
			"Status", uint8(user.Status),
			"LastActiveAt", ctxi.TimeStamp)
		pipe.Expire(ctx.BaseContext(), loginUserKey, confdao.Conf.Customize.TokenMaxAge)
		return nil
	}); err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "EfficientUserHashToRedis")
	}
	return nil
}

/*
创建空白哈希表时， 程序默认使用 REDIS_ENCODING_ZIPLIST 编码， 当以下任何一个条件被满足时， 程序将编码从 REDIS_ENCODING_ZIPLIST 切换为 REDIS_ENCODING_HT ：

哈希表中某个键或某个值的长度大于 server.hash_max_ziplist_value （默认值为 64 ）。
压缩列表中的节点数量大于 server.hash_max_ziplist_entries （默认值为 512 ）。
*/
func (d *UserDao) EfficientUserHashFromRedis() error {
	defer d.StartSpanEnd("EfficientUserHashFromRedis")
	ctxi := d
	ctx := ctxi.BaseContext()
	loginUser := modelconst.LoginUserKey + ctxi.AuthID

	userArgs, err := redisi.Strings(d.Do(ctx, redisi.CommandHGETALL, loginUser).Result())
	log.Debug(userArgs)
	if err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "EfficientUserHashFromRedis")
	}
	if len(userArgs) == 0 {
		return model.UserErrLoginTimeout
	}
	user := ctxi.AuthInfo.(*model.AuthInfo)
	user.Name = userArgs[1]
	n, err := strconv.ParseUint(userArgs[3], 10, 32)
	user.Role = model.Role(n)
	n, err = strconv.ParseUint(userArgs[5], 10, 8)
	user.Status = model.UserStatus(n)
	return nil
}

func (d *UserDao) UserLastActiveTime() error {
	ctxi := d
	ctx := ctxi.BaseContext()
	loginUser := modelconst.LoginUserKey + ctxi.AuthID
	if _, err := d.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, common.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, modelconst.LoginUserKey+"ActiveTime",
			&redis.Z{Score: float64(ctxi.TimeStamp), Member: ctxi.AuthID})
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "UserLastActiveTime")
	}
	return nil
}

func (d *UserDao) RedisUserInfoEdit(field string, value interface{}) error {
	ctxi := d
	ctx := ctxi.BaseContext()
	key := modelconst.LoginUserKey + ctxi.AuthID

	err := d.HSet(ctx, key, field, value).Err()
	if err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "RedisUserInfoEdit")
	}
	return nil
}

func (d *UserDao) GetUserExtRedis() (*model.UserExt, error) {
	ctxi := d
	ctx := ctxi.BaseContext()
	key := modelconst.UserExtKey + ctxi.AuthID

	userExt, err := redisi.Strings(d.Do(ctx, redisi.CommandHGETALL, key).Result())
	if err != nil {
		return nil, ctxi.ErrorLog(errcode.RedisErr, err, "GetUserExtRedis")
	}
	if len(userExt) > 3 {
		followCount, _ := strconv.ParseUint(userExt[1], 10, 64)
		followedCount, _ := strconv.ParseUint(userExt[3], 10, 64)
		return &model.UserExt{
			FollowCount:   followCount,
			FollowedCount: followedCount,
		}, nil
	}
	return nil, nil
}
