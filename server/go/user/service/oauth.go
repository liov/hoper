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
	goauth "github.com/hopeio/protobuf/oauth"
	"github.com/hopeio/protobuf/response"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/oauth"
	oauth3 "github.com/hopeio/utils/net/http/oauth"
	"github.com/hopeio/utils/types/param"
	jwti "github.com/hopeio/utils/validation/auth/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/global"
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
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", global.Conf.Customize.TokenSecretBytes, jwt.SigningMethodHS512))

	clientStore := NewClientStore(global.Dao.GORMDB.DB)

	manager.MapClientStorage(clientStore)

	srv := oauth3.NewServer(server.NewConfig(), manager)

	srv.UserAuthorizationHandler = func(token string) (userID string, err error) {
		if token == "" {
			return "", errors.ErrInvalidAccessToken
		}
		claims := new(jwti.Claims[uint64])
		if _, err := jwti.ParseToken(claims, token, global.Conf.Customize.TokenSecretBytes); err != nil {
			return "", err
		}
		return strconv.FormatUint(claims.Data, 10), nil
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
	var res httpi.ResponseRecorder
	u.Server.HandleAuthorizeRequest(ctx, &param.OauthReq{}, tokens[0], &res)
	headers := make(map[string]string)
	for k, v := range res.Headers {
		headers[k] = v[0]
	}
	return &response.HttpResponse{Body: res.Body.Bytes(), Status: uint32(res.Code), Headers: headers}, nil
}

func (u *OauthService) OauthToken(ctx context.Context, req *goauth.OauthReq) (*response.HttpResponse, error) {
	req.GrantType = string(oauth2.AuthorizationCode)
	var res httpi.ResponseRecorder
	u.Server.HandleTokenRequest(ctx, &param.OauthReq{
		ResponseType:   req.ResponseType,
		ClientID:       req.ClientID,
		Scope:          req.Scope,
		RedirectURI:    req.RedirectURI,
		State:          req.State,
		UserID:         req.UserID,
		AccessTokenExp: req.AccessTokenExp,
		ClientSecret:   req.ClientSecret,
		Code:           req.Code,
		RefreshToken:   req.RefreshToken,
		GrantType:      req.GrantType,
		AccessType:     req.AccessType,
		LoginURI:       req.LoginURI,
	}, &res)
	headers := make(map[string]string)
	for k, v := range res.Headers {
		headers[k] = v[0]
	}
	return &response.HttpResponse{Body: res.Body.Bytes(), Status: uint32(res.Code), Headers: headers}, nil
}
