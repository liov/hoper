package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

//如果只存一个id，jwt的意义在哪呢，跟session_id有什么区别
//jwt应该存放一些用户不能更改的信息，所以不能全存在jwt里
//或者说用户每更改一次信息就刷新token（貌似可行）
//有泛型这里多好写
type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(UserID uint64, now, maxAge int64, secret string) (string, error) {
	claims := Claims{
		UserID: UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + maxAge,
			IssuedAt:  now,
			Issuer:    "hoper",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	token, err := tokenClaims.SignedString(strings2.ToBytes(secret))

	return token, err
}

func ParseToken(token, secret string) (*Claims, error) {
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return strings2.ToBytes(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return nil, model.UserErr_LoginTimeout
			}
			return claims, nil
		}
	}

	return nil, model.UserErr_LoginError
}

