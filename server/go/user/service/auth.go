package service

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/cherry/context/httpctx"
	stringsi "github.com/hopeio/cherry/utils/strings"
	jwti "github.com/hopeio/cherry/utils/validation/auth/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"
	"strings"
	"time"

	"github.com/liov/hoper/server/go/user/confdao"
	"github.com/liov/hoper/server/go/user/data"
)

var ExportAuth = auth

type Authorization struct {
	*user.AuthInfo `json:"auth"`
	jwt.RegisteredClaims
	AuthInfoRaw string `json:"-"`
}

func (x *Authorization) Validate() error {
	return nil
}

func (x *Authorization) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func (x *Authorization) ParseToken(token string, secret []byte) error {
	_, err := jwti.ParseToken(x, token, secret)
	if err != nil {
		return err
	}
	x.ID = x.AuthInfo.IdStr()
	authBytes, _ := json.Marshal(x.AuthInfo)
	x.AuthInfoRaw = stringsi.BytesToString(authBytes)
	return nil
}

func auth(ctx *httpctx.Context, update bool) (*user.AuthInfo, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := confdao.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*Authorization)
		err := cache.Validate()
		if err != nil {
			return nil, err
		}
		authInfo := cache.AuthInfo
		return authInfo, nil
	}

	authorization := Authorization{AuthInfo: &user.AuthInfo{}}
	if err := authorization.ParseToken(ctx.Token, confdao.Conf.Customize.TokenSecretBytes); err != nil {
		return nil, user.UserErrNoLogin
	}

	authInfo := authorization.AuthInfo
	ctx.AuthID = authInfo.IdStr()
	ctx.AuthInfo = authInfo
	ctx.AuthInfoRaw = authorization.AuthInfoRaw

	if update {
		userDao := data.GetRedisDao(ctx, confdao.Dao.Redis.Client)
		err := userDao.EfficientUserHashFromRedis()
		if err != nil {
			return nil, err
		}
	}
	if !ok {
		confdao.Dao.Cache.SetWithTTL(signature, authorization, 0, 5*time.Second)
	}
	return authInfo, nil
}
