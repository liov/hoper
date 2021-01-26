package dao

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"github.com/liov/hoper/go/v2/utils/encoding/hash"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserRedis struct {
	*redis.Conn
}

func (conn *UserRedis) Close() {
	conn.Conn.Close()
}

func NewUserRedis() *UserRedis {
	conn := Dao.Redis.Conn(Dao.Redis.Context())
	return &UserRedis{conn}
}

// UserToRedis 将用户信息存到redis
func (conn *UserRedis) UserToRedis(ctx *model.Ctx, user *model.UserAuthInfo) error {

	UserString, err := json.Standard.MarshalToString(user)
	if err != nil {
		return err
	}

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.SetEX(ctx, loginUserKey, UserString, conf.Conf.Customize.TokenMaxAge)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (conn *UserRedis) UserFromRedis(ctx *model.Ctx, userID uint64) (*model.AuthInfo, error) {

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)
	var cmd *redis.StringCmd
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		cmd = pipe.Get(ctx, loginUser)
		return nil
	}); redisErr != nil {
		return nil, redisErr
	}
	userString, err := redisi.String(cmd.Result())
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

func (conn *UserRedis) EditRedisUser(ctx *model.Ctx, user *model.AuthInfo) error {

	UserString, err := json.Standard.MarshalToString(user)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.Set(ctx, loginUserKey, UserString, 0)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func (conn *UserRedis) UserHashToRedis(ctx *model.Ctx, user *model.AuthInfo) error {

	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	redisArgs = append(redisArgs, redisi.HMSET, loginUserKey)
	redisArgs = append(redisArgs, hash.Marshal(user)...)
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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
func (conn *UserRedis) UserHashFromRedis(ctx *model.Ctx, user *model.AuthInfo) error {

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	var cmd *redis.Cmd
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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
	hash.UnMarshal(user, userArgs)
	return nil
}

func (conn *UserRedis) EfficientUserHashToRedis(ctx *model.Ctx) error {
	user := ctx.AuthInfo
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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
func (conn *UserRedis) EfficientUserHashFromRedis(ctx *model.Ctx, user *model.AuthInfo) error {
	loginUser := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	var cmd *redis.Cmd
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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
	user.Name = userArgs[1]
	n, err := strconv.ParseUint(userArgs[3], 10, 32)
	user.Role = model.Role(n)
	n, err = strconv.ParseUint(userArgs[5], 10, 8)
	user.Status = model.UserStatus(n)
	return nil
}

func (conn *UserRedis) UserLastActiveTime(ctx *model.Ctx, userID uint64, now time.Time) error {

	id := strconv.FormatUint(userID, 10)
	loginUser := modelconst.LoginUserKey + id
	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.CronIndex)
		//有序集合存一份，遍历长时间未活跃用户用
		pipe.ZAdd(ctx, modelconst.LoginUserKey+"ActiveTime",
			&redis.Z{Score: float64(now.Unix()), Member: id})
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.HSet(ctx, loginUser, "LastActiveAt")
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}

func (conn *UserRedis) RedisUserInfoEdit(ctx *model.Ctx, userID uint64, key string, value interface{}) error {

	id := strconv.FormatUint(userID, 10)
	loginUser := modelconst.LoginUserKey + id

	if _, redisErr := conn.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Select(ctx, modelconst.UserIndex)
		pipe.HSet(ctx, loginUser, key, value)
		return nil
	}); redisErr != nil {
		return redisErr
	}
	return nil
}
