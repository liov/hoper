package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/liov/hoper/server/go/lib/utils/strings"
)

var (
	Parser = jwt.NewParser()
	err    = errors.New("无效的token")
)

// 如果只存一个id，jwt的意义在哪呢，跟session_id有什么区别
// jwt应该存放一些用户不能更改的信息，所以不能全存在jwt里
// 或者说用户每更改一次信息就刷新token（貌似可行）
// 有泛型这里多好写
type Claims struct {
	UserId uint64 `json:"userId"`
	*jwt.StandardClaims
}

func (claims *Claims) Valid(helper *jwt.ValidationHelper) error {
	return helper.ValidateExpiresAt(claims.ExpiresAt)
}

type CustomClaims struct {
	CustomInfo interface{} `json:"customInfo"`
	*jwt.StandardClaims
}

func (claims *CustomClaims) Valid(helper *jwt.ValidationHelper) error {
	return helper.ValidateExpiresAt(claims.ExpiresAt)
}

func NewStandardClaims(maxAge int64, sign string) *jwt.StandardClaims {
	now := time.Now()
	exp := now.Add(time.Duration(maxAge))
	return &jwt.StandardClaims{
		ExpiresAt: &jwt.Time{Time: exp},
		IssuedAt:  &jwt.Time{Time: now},
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
		return err
	}
	tokenClaims, err := Parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return stringsi.ToBytes(secret), nil
	})
	if err != nil {
		return err
	}
	if tokenClaims != nil && tokenClaims.Valid {
		return tokenClaims.Claims.Valid(jwt.DefaultValidationHelper)
	}

	return err
}

func ParseTokenWithKeyFunc(claims jwt.Claims, token string, f func(token *jwt.Token) (interface{}, error)) error {
	tokenClaims, _ := Parser.ParseWithClaims(token, claims, f)

	if tokenClaims != nil && tokenClaims.Valid {
		return tokenClaims.Claims.Valid(Parser.ValidationHelper)
	}

	return err
}
