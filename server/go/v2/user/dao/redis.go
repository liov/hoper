package dao

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/protobuf/common"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"github.com/liov/hoper/go/v2/utils/encoding/hash"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserRedis struct {
}

// UserToRedis 将用户信息存到redis
func (*UserRedis) UserToRedis(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	UserString, err := json.Standard.MarshalToString(ctxi.AuthInfo)
	if err != nil {
		return err
	}

	loginUserKey := modelconst.LoginUserKey + ctxi.IdStr
	if redisErr := Dao.Redis.SetEX(ctx, loginUserKey, UserString, conf.Conf.Customize.TokenMaxAge).Err(); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (*UserRedis) UserFromRedis(ctxi *model.Ctx) (*model.AuthInfo, error) {
	ctx:=ctxi.Context
	loginUser := modelconst.LoginUserKey + ctxi.IdStr

	userString, err := redisi.String(Dao.Redis.Get(ctx, loginUser).Result())
	if err != nil {
		return nil, err
	}
	var user model.AuthInfo
	err = json.Standard.UnmarshalFromString(userString, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserRedis) EditRedisUser(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	UserString, err := json.Standard.MarshalToString(ctxi.AuthInfo)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + ctxi.IdStr
	return Dao.Redis.Set(ctx, loginUserKey, UserString, 0).Err()
}

// UserToRedis 将用户信息存到redis
func (*UserRedis) UserHashToRedis(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + ctxi.IdStr
	redisArgs = append(redisArgs, redisi.HMSET, loginUserKey)
	redisArgs = append(redisArgs, hash.Marshal(ctxi.AuthInfo)...)
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Do(ctx, redisArgs...)
		pipe.Expire(ctx, loginUserKey, conf.Conf.Customize.TokenMaxAge)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (*UserRedis) UserHashFromRedis(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	loginUser := modelconst.LoginUserKey + ctxi.IdStr

	userArgs, err := redisi.Strings(Dao.Redis.Do(ctx, redisi.HGETALL, loginUser).Result())
	if err != nil {
		return err
	}
	log.Debug(userArgs)
	if len(userArgs) == 0 {
		return model.UserErrInvalidToken
	}
	hash.UnMarshal(ctxi.AuthInfo, userArgs)
	return nil
}

func (*UserRedis) EfficientUserHashToRedis(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	user := ctxi.AuthInfo
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, loginUserKey, "Name", user.Name,
			"Role", uint32(user.Role),
			"Status", uint8(user.Status),
			"LastActiveAt", user.LastActiveAt)
		pipe.Expire(ctx, loginUserKey, conf.Conf.Customize.TokenMaxAge)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

/*
创建空白哈希表时， 程序默认使用 REDIS_ENCODING_ZIPLIST 编码， 当以下任何一个条件被满足时， 程序将编码从 REDIS_ENCODING_ZIPLIST 切换为 REDIS_ENCODING_HT ：

哈希表中某个键或某个值的长度大于 server.hash_max_ziplist_value （默认值为 64 ）。
压缩列表中的节点数量大于 server.hash_max_ziplist_entries （默认值为 512 ）。
*/
func (*UserRedis) EfficientUserHashFromRedis(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	loginUser := modelconst.LoginUserKey + ctxi.IdStr

	userArgs, err := redisi.Strings(Dao.Redis.Do(ctx, redisi.HGETALL, loginUser).Result())
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(userArgs) == 0 {
		return model.UserErrInvalidToken
	}
	ctxi.AuthInfo.Name = userArgs[1]
	n, err := strconv.ParseUint(userArgs[3], 10, 32)
	ctxi.AuthInfo.Role = model.Role(n)
	n, err = strconv.ParseUint(userArgs[5], 10, 8)
	ctxi.AuthInfo.Status = model.UserStatus(n)
	return nil
}

func (*UserRedis) UserLastActiveTime(ctxi *model.Ctx) error {
	ctx:=ctxi.Context
	loginUser := modelconst.LoginUserKey + ctxi.IdStr
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, common.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, modelconst.LoginUserKey+"ActiveTime",
			&redis.Z{Score: float64(ctxi.RequestUnix), Member: ctxi.IdStr})
		pipe.Select(ctx, conf.Conf.Redis.Index)
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

func (*UserRedis) RedisUserInfoEdit(ctxi *model.Ctx, key string, value interface{}) error {
	ctx:=ctxi.Context
	loginUser := modelconst.LoginUserKey + ctxi.IdStr

	return Dao.Redis.HSet(ctx, loginUser, key, value).Err()
}
