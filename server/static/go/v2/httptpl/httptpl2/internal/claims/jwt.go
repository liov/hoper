package claims

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/liov/hoper/go/v2/httptpl/httptpl2/internal/config"
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

func (claims *Claims) ParseToken(req *http.Request) error {
	claims.Id = req.URL.Path
	return errors.New("未登录")
}
