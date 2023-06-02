package service

import (
	"context"
	"github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/actliboy/hoper/server/go/user/confdao"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	goauth "github.com/hopeio/pandora/protobuf/oauth"
	"github.com/hopeio/pandora/protobuf/response"
	jwti "github.com/hopeio/pandora/utils/verification/auth/jwt"
	"github.com/hopeio/pandora/utils/verification/auth/oauth"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
	"time"
)

func GetOauthService() *OauthService {
	if oauthSvc != nil {
		return oauthSvc
	}
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", confdao.Conf.Customize.TokenSecret, jwt.SigningMethodHS512))

	clientStore := oauth.NewClientStore(confdao.Dao.GORMDB.DB)

	manager.MapClientStorage(clientStore)

	srv := oauth.NewServer(server.NewConfig(), manager)

	srv.UserAuthorizationHandler = func(token string) (userID string, err error) {
		if token == "" {
			return "", errors.ErrInvalidAccessToken
		}
		claims := new(jwti.Claims)
		if _, err := jwti.ParseToken(claims, token, confdao.Conf.Customize.TokenSecret); err != nil {
			return "", err
		}
		return strconv.FormatUint(claims.UserId, 10), nil
	}

	srv.InternalErrorHandler = func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	}

	srv.ResponseErrorHandler = func(re *errors.Response) {
		log.Println("HttpResponse Error:", re.Error.Error())
	}
	oauthSvc = &OauthService{Server: srv, ClientStore: clientStore}
	return oauthSvc
}

type OauthService struct {
	Server      *oauth.Server
	ClientStore *oauth.ClientStore
	user.UnimplementedOauthServiceServer
}

func (u *OauthService) OauthAuthorize(ctx context.Context, req *goauth.OauthReq) (*response.HttpResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("auth")
	tokens = append(tokens, "")
	req.AccessTokenExp = int64(24 * time.Hour)
	req.LoginURI = "/oauth/login"
	res := u.Server.HandleAuthorizeRequest(ctx, req, tokens[0])
	return res, nil
}

func (u *OauthService) OauthToken(ctx context.Context, req *goauth.OauthReq) (*response.HttpResponse, error) {
	req.GrantType = string(oauth2.AuthorizationCode)
	res, _ := u.Server.HandleTokenRequest(ctx, req)
	return res, nil
}
