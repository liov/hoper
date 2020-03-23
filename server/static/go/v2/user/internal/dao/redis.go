package dao

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	"github.com/liov/hoper/go/v2/utils/dao/redis/hash"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserRedis struct {
	redis.Conn
}

func (conn *UserRedis) Close() {
	conn.Conn.Close()
}

func NewUserRedis() *UserRedis {
	conn := Dao.Redis.Get()
	return &UserRedis{conn}
}

// UserToRedis 将用户信息存到redis
func (conn *UserRedis) UserToRedis(user *model.UserMainInfo) error {

	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString, "EX", config.Conf.Customize.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (conn *UserRedis) UserFromRedis(userID uint64) (*model.UserMainInfo, error) {

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn.Send("SELECT", modelconst.UserIndex)
	userString, err := redis.String(conn.Do("GET", loginUser))
	if err != nil {
		return nil, err
	}
	var user model.UserMainInfo
	err = json.Json.UnmarshalFromString(userString, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (conn *UserRedis) EditRedisUser(user *model.UserMainInfo) error {

	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserToRedis 将用户信息存到redis
func (conn *UserRedis) UserHashToRedis(user *model.UserMainInfo) error {

	var redisArgs []interface{}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)
	redisArgs = append(redisArgs, loginUserKey)
	redisArgs = append(redisArgs, hash.Marshal(user)...)

	conn.Send("SELECT", modelconst.UserIndex)
	conn.Send("HMSET", redisArgs...)
	if _, redisErr := conn.Do("EXPIRE", loginUserKey, config.Conf.Customize.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func (conn *UserRedis) UserHashFromRedis(userID uint64) (*model.UserMainInfo, error) {

	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn.Send("SELECT", modelconst.UserIndex)
	userArgs, err := redis.Strings(conn.Do("HGETALL", loginUser))
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(userArgs) == 0 {
		return nil, model.UserErr_InvalidToken
	}
	var user model.UserMainInfo
	hash.UnMarshal(&user, userArgs)
	return &user, nil
}

func (conn *UserRedis) EfficientUserHashToRedis(user *model.UserMainInfo) error {

	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn.Send("SELECT", modelconst.UserIndex)
	conn.Send("HMSET", loginUserKey,
		"Id", user.Id,
		"Role", uint32(user.Role),
		"LastActiveAt", user.LastActiveAt,
		"Score", user.Score,
		"Status", uint8(user.Status),
		"LoginTime", user.LoginTime)
	if _, redisErr := conn.Do("EXPIRE", loginUserKey, config.Conf.Customize.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

/*
创建空白哈希表时， 程序默认使用 REDIS_ENCODING_ZIPLIST 编码， 当以下任何一个条件被满足时， 程序将编码从 REDIS_ENCODING_ZIPLIST 切换为 REDIS_ENCODING_HT ：

哈希表中某个键或某个值的长度大于 server.hash_max_ziplist_value （默认值为 64 ）。
压缩列表中的节点数量大于 server.hash_max_ziplist_entries （默认值为 512 ）。
*/
func (conn *UserRedis) EfficientUserHashFromRedis(userID uint64) (*model.UserMainInfo, error) {
	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn.Send("SELECT", modelconst.UserIndex)
	userArgs, err := redis.Strings(conn.Do("HGETALL", loginUser))
	log.Debug(userArgs)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(userArgs) == 0 {
		return nil, model.UserErr_InvalidToken
	}
	var user model.UserMainInfo
	user.Id, err = strconv.ParseUint(userArgs[1], 10, 64)
	n, err := strconv.ParseUint(userArgs[3], 10, 32)
	user.Role = model.Role(n)
	user.LastActiveAt, err = strconv.ParseInt(userArgs[5], 10, 64)
	user.Score, err = strconv.ParseUint(userArgs[7], 10, 64)
	n, err = strconv.ParseUint(userArgs[9], 10, 8)
	user.Status = model.UserStatus(n)
	user.LoginTime, err = strconv.ParseInt(userArgs[11], 10, 64)
	return &user, nil
}

func (conn *UserRedis) UserLastActiveTime(userID uint64, now time.Time) error {

	id := strconv.FormatUint(userID, 10)
	loginUser := modelconst.LoginUserKey + id
	err := conn.Send("SELECT", modelconst.CronIndex)
	//有序集合存一份，遍历长时间未活跃用户用
	_, err = conn.Do("ZADD", modelconst.LoginUserKey+"ActiveTime",
		now.Unix(), id)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", loginUser, "LastActiveAt",
		now.Unix())
	if err != nil {
		return err
	}
	return nil
}

func (conn *UserRedis) RedisUserInfoEdit(userID uint64, key string, value interface{}) error {

	id := strconv.FormatUint(userID, 10)
	loginUser := modelconst.LoginUserKey + id

	err := conn.Send("SELECT", modelconst.CronIndex)
	_, err = conn.Do("HSET", loginUser, key, value)
	if err != nil {
		return err
	}
	return nil
}
