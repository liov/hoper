package jwti

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
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
	*jwt.RegisteredClaims
}

type CustomClaims struct {
	CustomInfo interface{} `json:"customInfo"`
	*jwt.RegisteredClaims
}

func NewStandardClaims(maxAge int64, sign string) *jwt.RegisteredClaims {
	now := time.Now()
	exp := now.Add(time.Duration(maxAge))
	return &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: exp},
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    sign,
	}
}

func GenerateToken(claims jwt.Claims, secret interface{}) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func ParseToken(claims jwt.Claims, token string, secret []byte) (*jwt.Token, error) {
	return Parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func ParseTokenWithKeyFunc(claims jwt.Claims, token string, f func(token *jwt.Token) (interface{}, error)) (*jwt.Token, error) {

	return Parser.ParseWithClaims(token, claims, f)
}
