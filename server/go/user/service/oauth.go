package service

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	oauth3 "github.com/hopeio/cherry/oauth"
	goauth "github.com/hopeio/protobuf/oauth"
	"github.com/hopeio/protobuf/response"
	"github.com/hopeio/utils/net/http/oauth"
	jwti "github.com/hopeio/utils/validation/auth/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/confdao"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
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
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", confdao.Conf.Customize.TokenSecretBytes, jwt.SigningMethodHS512))

	clientStore := NewClientStore(confdao.Dao.GORMDB.DB)

	manager.MapClientStorage(clientStore)

	srv := oauth3.NewServer(server.NewConfig(), manager)

	srv.UserAuthorizationHandler = func(token string) (userID string, err error) {
		if token == "" {
			return "", errors.ErrInvalidAccessToken
		}
		claims := new(jwti.Claims)
		if _, err := jwti.ParseToken(claims, token, confdao.Conf.Customize.TokenSecretBytes); err != nil {
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
	Server      *oauth3.Server
	ClientStore oauth.ClientStore
	user.UnimplementedOauthServiceServer
}

// NewClientStore create client store
func NewClientStore(db *gorm.DB) *ClientStore {
	return (*ClientStore)(db)
}

// ClientStore client information store
type ClientStore gorm.DB

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	db := (*gorm.DB)(cs)
	var client models.Client
	if err := db.Table("oauth_client").Find(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

// Set set client information
func (cs *ClientStore) Set(cli oauth2.ClientInfo) (err error) {
	db := (*gorm.DB)(cs)
	db.Table("oauth_client").Create(cli)
	return
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
