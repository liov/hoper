package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v5"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/protobuf/response"
	jwtx "github.com/hopeio/scaffold/jwt"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/user"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

func GetOauthService() *OauthService {
	if oauthSvc != nil {
		return oauthSvc
	}
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", global.Conf.User.TokenSecretBytes, jwt.SigningMethodHS512))

	clientStore := NewClientStore(global.Dao.GORMDB.DB)
	manager.MapClientStorage(clientStore)

	srv := NewServer(server.NewConfig(), manager)
	srv.UserAuthorizationHandler = func(token string) (userID string, err error) {
		if token == "" {
			return "", errors.ErrInvalidAccessToken
		}
		claims := new(jwtx.Claims[uint64])
		if _, err := jwtx.ParseToken(claims, token, global.Conf.User.TokenSecretBytes); err != nil {
			return "", err
		}
		return strconv.FormatUint(claims.Auth, 10), nil
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
	Server      *Server
	ClientStore *ClientStore
	user.UnimplementedOauthServiceServer
}

func NewClientStore(db *gorm.DB) *ClientStore {
	return (*ClientStore)(db)
}

type ClientStore gorm.DB

func (cs *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	db := (*gorm.DB)(cs)
	var client models.Client
	if err := db.Table("oauth_client").Find(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (cs *ClientStore) Set(cli oauth2.ClientInfo) (err error) {
	db := (*gorm.DB)(cs)
	db.Table("oauth_client").Create(cli)
	return
}

func oauthReqFromPB(req *user.OauthReq) *Reuqest {
	return &Reuqest{
		ResponseType: req.ResponseType, ClientID: req.ClientID, Scope: req.Scope,
		RedirectURI: req.RedirectURI, State: req.State, UserID: req.UserID,
		AccessTokenExp: req.AccessTokenExp, ClientSecret: req.ClientSecret, Code: req.Code,
		RefreshToken: req.RefreshToken, GrantType: req.GrantType, AccessType: req.AccessType, LoginURI: req.LoginURI,
	}
}

func (u *OauthService) OauthAuthorize(ctx context.Context, req *user.OauthReq) (*response.HttpResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("auth")
	tokens = append(tokens, "")
	oauthReq := oauthReqFromPB(req)
	oauthReq.AccessTokenExp = int64(24 * time.Hour)
	oauthReq.LoginURI = "/oauth/login"
	var res httpx.Recorder
	u.Server.HandleAuthorizeRequest(ctx, oauthReq, tokens[0], &res)
	headers := make(map[string]string)
	for k, v := range res.ResponseRecorder.Header() {
		headers[k] = v[0]
	}
	return &response.HttpResponse{Body: res.ResponseRecorder.Body.Bytes(), Status: int32(res.StatusCode), Headers: headers}, nil
}

func (u *OauthService) OauthToken(ctx context.Context, req *user.OauthReq) (*response.HttpResponse, error) {
	oauthReq := oauthReqFromPB(req)
	oauthReq.GrantType = string(oauth2.AuthorizationCode)
	var res httpx.Recorder
	if err := u.Server.HandleTokenRequest(ctx, oauthReq, &res); err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	for k, v := range res.Header() {
		headers[k] = v[0]
	}
	return &response.HttpResponse{Body: res.ResponseRecorder.Body.Bytes(), Status: int32(res.StatusCode), Headers: headers}, nil
}
