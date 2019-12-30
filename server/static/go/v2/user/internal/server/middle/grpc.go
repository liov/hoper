package middle

import (
	"context"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	modelconst "github.com/liov/hoper/go/v2/user/internal/model"
	"github.com/liov/hoper/go/v2/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/strings2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("authorization")
	authErr := status.Error(codes.Code(errorcode.Auth), errorcode.Auth.Error())
	if len(tokens) == 0 || tokens[0] == "" {
		err = authErr
		return
	}
	token, err := ParseToken(tokens[0])
	if err != nil {
		err = authErr
		return
	}
	user, err := UserFromRedis(token.UserID)
	if err != nil {
		err = authErr
		return
	}
	ctx = context.WithValue(ctx, "auth", user)
	return handler(ctx, req)
}

type Claims struct {
	UserID   uint64 `json:"user_id"`
	UserRole uint32 `json:"user_role"`
	jwt.StandardClaims
}

func GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.Id,
		UserRole: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + config.Conf.Server.TokenMaxAge,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "hoper",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(strings2.ToBytes(config.Conf.Server.TokenSecret))

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return strings2.ToBytes(config.Conf.Server.TokenSecret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return nil, errorcode.LoginTimeout
			}
			return claims, nil
		}
	}

	return nil, errorcode.LoginError
}

// UserToRedis 将用户信息存到redis
func UserToRedis(user *model.User) error {
	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString, "EX", config.Conf.Server.TokenMaxAge); redisErr != nil {
		return redisErr
	}
	return nil
}

// UserFromRedis 从redis中取出用户信息
func UserFromRedis(userID uint64) (*model.User, error) {
	loginUser := modelconst.LoginUserKey + strconv.FormatUint(userID, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	userString, err := redis.String(conn.Do("GET", loginUser))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var user model.User
	err = json.Json.UnmarshalFromString(userString, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserLastActiveTime(userID uint64) error {
	conn := dao.Dao.Redis.Get()
	defer conn.Close()

	err := conn.Send("SELECT", modelconst.CronIndex)
	_, err = conn.Do("ZADD", modelconst.LoginUserKey+"ActiveTime",
		time.Now().Unix(), strconv.FormatUint(userID, 10))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func EditUserRedis(user *model.User) error {
	UserString, err := json.Json.MarshalToString(user)
	if err != nil {
		return err
	}
	loginUserKey := modelconst.LoginUserKey + strconv.FormatUint(user.Id, 10)

	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	conn.Send("SELECT", modelconst.UserIndex)
	if _, redisErr := conn.Do("SET", loginUserKey, UserString); redisErr != nil {
		return redisErr
	}
	return nil
}
