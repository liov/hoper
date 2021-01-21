package service

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func init() {
	pick.RegisterService(&UserService{}, &TestService{}, &StaticService{})
}

type Claims struct {
	User *user.UserMainInfo
	jwt.StandardClaims
}

func (claims *Claims) GenerateToken() (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + int64(24*time.Hour),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "hoper",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString("secret")

	return token, err
}

func (claims *Claims) ParseToken(req *http.Request) error {
	var token string
	cookie, err := req.Cookie("token")
	if err == nil {
		token, _ = url.QueryUnescape(cookie.Value)
	}
	if token == "" {
		token = req.Header.Get("HeaderAuthorization")
	}
	if token == "" {
		return errors.New("未登录")
	}
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return "secret", nil
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
