package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/strings"
)

var Parser = &jwt.Parser{SkipClaimsValidation: true}
//如果只存一个id，jwt的意义在哪呢，跟session_id有什么区别
//jwt应该存放一些用户不能更改的信息，所以不能全存在jwt里
//或者说用户每更改一次信息就刷新token（貌似可行）
//有泛型这里多好写
type Claims struct {
	UserId uint64 `json:"userId"`
	*jwt.StandardClaims
}

func (claims *Claims) Valid() error {
	if claims.VerifyExpiresAt(claims.IssuedAt, false) == false {
		return errorcode.TimeTooMuch
	}
	return nil
}

type CustomClaims struct {
	CustomInfo interface{} `json:"customInfo"`
	*jwt.StandardClaims
}

func (claims *CustomClaims) Valid() error {
	if claims.VerifyExpiresAt(claims.IssuedAt, false) == false {
		return errorcode.TimeTooMuch
	}
	return nil
}

func NewStandardClaims(maxAge int64,sign string) *jwt.StandardClaims {
	now := time.Now().Unix()
	return &jwt.StandardClaims{
		ExpiresAt: now + maxAge,
		IssuedAt:  now,
		Issuer:    sign,
	}
}

func GenerateToken(claims jwt.Claims, secret interface{}) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func ParseToken(claims jwt.Claims, token, secret string) error {
	if token == "" {
		return model.UserErr_NoLogin
	}
	tokenClaims, _ := Parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return stringsi.ToBytes(secret), nil
		})

	if tokenClaims != nil && tokenClaims.Valid {
		return tokenClaims.Claims.Valid()
	}

	return model.UserErr_LoginError
}

func ParseTokenWithKeyFunc(claims jwt.Claims, token string, f func(token *jwt.Token) (interface{}, error)) error {
	tokenClaims, _ := Parser.ParseWithClaims(token, claims, f)

	if tokenClaims != nil && tokenClaims.Valid {
		return tokenClaims.Claims.Valid()
	}

	return model.UserErr_LoginError
}
