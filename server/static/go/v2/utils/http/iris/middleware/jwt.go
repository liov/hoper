package middleware

import (
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/http/iris/response"
)

//中间件的两种方式
func GetUser(fullInfo bool, secret string) iris.Handler {
	return func(ctx iris.Context) {
		var user *model.Auth
		var userID uint64
		var err error
		var code errorcode.ErrCode
		if fullInfo {
			user, err = getUser(ctx, secret)
		} else {
			userID, err = getUserID(ctx, secret)
		}

		if err != nil && code != errorcode.SUCCESS {
			//ctx.StatusCode(iris.StatusUnauthorized)
			ctx.SetCookie(&http.Cookie{
				Name:     "token",
				Value:    "del",
				Path:     "/",
				Domain:   ctx.Host(),
				Expires:  time.Now().Add(-1),
				MaxAge:   -1,
				Secure:   false,
				HttpOnly: true,
			})
			response.Res(ctx, int(errorcode.ERROR), err.Error(), nil)
			return
		}

		if fullInfo {
			ctx.Values().Set("user", user) //指针
		}
		ctx.Values().Set("userID", userID)
		ctx.Next()
	}
}

func GetUserId(ctx iris.Context, secret string) {
	if userID, err := getUserID(ctx, secret); err == nil {
		ctx.Values().Set("userID", userID)
	} else {
		ctx.Values().Set("userID", uint64(0))
	}
	ctx.Next()
}

func getUser(ctx iris.Context, secret string) (*model.Auth, error) {

	tokenString := ctx.GetCookie("token")
	if len(tokenString) == 0 {
		tokenString = ctx.GetHeader("Authorization")
	}
	if len(tokenString) == 0 {
		return nil, errorcode.LoginError
	}

	claims, err := token.ParseToken(tokenString, secret)

	if err != nil {
		return nil, err
	}
	user, err := token.UserFromRedis(claims.UserID)

	return user, nil
}

func getUserID(ctx iris.Context, secret string) (uint64, error) {
	tokenString := ctx.GetCookie("token")
	if len(tokenString) == 0 {
		tokenString = ctx.GetHeader("Authorization")
	}
	if len(tokenString) == 0 {
		return 0, errorcode.LoginError
	}

	claims, err := token.ParseToken(tokenString, secret)

	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

//Config全局变量太大了
//var jwtSecret = utils.ToBytes(initialize.Config.Server.JwtSecret))
