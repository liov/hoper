package dao

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"github.com/liov/hoper/go/v2/utils/encoding/hash"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/koko"
)

type UserRedis struct {
}


// UserToRedis 将用户信息存到redis
func (*UserRedis) UserHashToRedis(ctxi *koko.Ctx) error {
	ctx:=ctxi.Context
	user := ctxi.AuthInfo.(*model.AuthInfo)
	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + user.IdStr
	redisArgs = append(redisArgs, redisi.HMSET, loginUserKey)
	redisArgs = append(redisArgs, hash.Marshal(ctxi.AuthInfo)...)
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.Do(ctx, redisArgs...)
		pipe.Expire(ctx, loginUserKey, conf.Conf.Customize.TokenMaxAge)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (*UserRedis) UserHashFromRedis(ctxi *koko.Ctx) error {
	ctx:=ctxi.Context
	user := ctxi.AuthInfo.(*model.AuthInfo)
	loginUser := modelconst.LoginUserKey + user.IdStr
	var cmd *redis.Cmd
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		cmd = pipe.Do(ctx, redisi.HGETALL, loginUser)
		return nil
	}); redisErr != nil {
		return redisErr
	}

	userArgs, err := redisi.Strings(cmd.Result())
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(userArgs) == 0 {
		return model.UserErr_InvalidToken
	}
	hash.UnMarshal(ctxi.AuthInfo, userArgs)
	return nil
}

func (*UserRedis) EfficientUserHashToRedis(ctxi *koko.Ctx) error {
	ctx:=ctxi.Context
	user := ctxi.AuthInfo.(*model.AuthInfo)
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
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
func (*UserRedis) EfficientUserHashFromRedis(ctxi *koko.Ctx) error {
	ctx:=ctxi.Context
	authInfo := ctxi.AuthInfo.(*model.AuthInfo)
	loginUser := modelconst.LoginUserKey + authInfo.IdStr
	var cmd *redis.Cmd
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		cmd = pipe.Do(ctx, redisi.HGETALL, loginUser)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	userArgs, err := redisi.Strings(cmd.Result())
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(userArgs) == 0 {
		return model.UserErr_InvalidToken
	}
	authInfo.Name = userArgs[1]
	n, err := strconv.ParseUint(userArgs[3], 10, 32)
	authInfo.Role = model.Role(n)
	n, err = strconv.ParseUint(userArgs[5], 10, 8)
	authInfo.Status = model.UserStatus(n)
	return nil
}

func (*UserRedis) UserLastActiveTime(ctxi *koko.Ctx) error {
	ctx:=ctxi.Context
	authInfo := ctxi.AuthInfo.(*model.AuthInfo)
	loginUser := modelconst.LoginUserKey + authInfo.IdStr
	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, modelconst.LoginUserKey+"ActiveTime",
			&redis.Z{Score: float64(ctxi.RequestUnix), Member: authInfo.IdStr})
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

func (*UserRedis) RedisUserInfoEdit(ctxi *koko.Ctx, key string, value interface{}) error {
	ctx:=ctxi.Context
	authInfo := ctxi.AuthInfo.(*model.AuthInfo)
	loginUser := modelconst.LoginUserKey + authInfo.IdStr

	if _, redisErr := Dao.Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.HSet(ctx, loginUser, key, value)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}
