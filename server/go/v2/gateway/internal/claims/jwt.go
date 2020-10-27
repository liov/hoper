package claims

import (
	"errors"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/hero"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/protobuf/user"
)

type Claims struct {
	User *user.UserMainInfo
	jwt.StandardClaims
}

func (claims *Claims) GenerateToken() (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + config.Conf.Customize.TokenMaxAge,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "hoper",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.Conf.Customize.TokenSecret)

	return token, err
}

func (claims *Claims) ParseToken(token string) error {
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Customize.TokenSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return errors.New("登录超时")
			}
			return nil
		}
	}
	return errors.New("未登录")
}

var SESSION = hero.Dependency{
	Handle: func(ctx *context.Context, input *hero.Input) (reflect.Value, error) {
		var claims = new(Claims)
		return reflect.ValueOf(claims), nil
	},
}
