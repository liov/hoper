package service

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/utils/net/http/auth"
	ijwt "github.com/liov/hoper/go/v2/utils/net/http/auth/jwt"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(r *http.Request) (*model.UserMainInfo, error) {
	token := auth.GetToken(r)
	if token == "" {
		return nil, model.UserErr_NoLogin
	}
	claims, err := ijwt.ParseToken(token, conf.Conf.Customize.TokenSecret)
	if err != nil {
		return nil, model.UserErr_InvalidToken
	}
	conn := dao.NewUserRedis()
	defer conn.Close()
	user, err := conn.EfficientUserHashFromRedis(claims.UserId)
	if err != nil {
		log.Error(err)
		return nil, model.UserErr_InvalidToken
	}
	return user, nil
}

func (*UserService) AuthClaims(ctx context.Context) (*ijwt.Claims, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("auth")
	if len(tokens) == 0 || tokens[0] == "" {
		return nil, model.UserErr_NoLogin
	}
	return ijwt.ParseToken(tokens[0], conf.Conf.Customize.TokenSecret)
}

func (u *UserService) AuthMainInfo(ctx context.Context) (*model.UserMainInfo, error) {
	claims, err := u.AuthClaims(ctx)
	if err != nil {
		return nil, err
	}
	//redis
	conn := dao.NewUserRedis()
	defer conn.Close()
	user, err := conn.EfficientUserHashFromRedis(claims.UserId)
	if err != nil {
		return nil, model.UserErr_InvalidToken
	}

	if user.LoginTime != claims.IssuedAt {
		return nil, model.UserErr_InvalidToken
	}
	return user, nil
}

func (*UserService) Device(ctx context.Context) *model.UserDeviceInfo {
	var info model.UserDeviceInfo
	md, _ := metadata.FromIncomingContext(ctx)
	//Device-Info:device-osInfo-appCode-appVersion
	deviceInfo := md.Get("device-info")
	if len(deviceInfo) > 0 && deviceInfo[0] != "" {
		infos := strings.Split(deviceInfo[0], "-")
		if len(infos) == 4 {
			info.Device = infos[0]
			info.Os = infos[1]
			info.AppCode = infos[2]
			info.AppVersion = infos[3]
		}
	}
	//area:xxx
	//location:1.23456,2.123456
	location := md.Get("location")
	if len(location) > 0 && location[0] != "" {
		info.Area, _ = url.PathUnescape(location[0])
	}

	if len(location) > 1 && location[1] != "" {
		infos := strings.Split(location[1], ",")
		if len(infos) == 2 {
			info.Lng = infos[0]
			info.Lat = infos[1]
		}
	}

	userAgent := md.Get("user-agent")
	if len(userAgent) > 0 && userAgent[0] != "" {
		info.UserAgent = userAgent[0]
	}
	ip := md.Get("x-forwarded-for")
	if len(ip) > 0 && ip[0] != "" {
		info.IP = ip[0]
	}
	return &info
}

type authKey struct{}

// AuthContext returns a new Context that carries value u.
func AuthContext(r *http.Request) context.Context {
	user, _ := Auth(r)
	return context.WithValue(r.Context(), authKey{}, user)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*model.UserMainInfo, bool) {
	user, ok := ctx.Value(authKey{}).(*model.UserMainInfo)
	return user, ok
}

type Claims struct {
	User *model.UserMainInfo
	jwt.StandardClaims
}

func (claims *Claims) GenerateToken() (string, error) {
	now := time.Now().Unix()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: now + int64(24*time.Hour),
		IssuedAt:  now,
		Issuer:    "hoper",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString("secret")

	return token, err
}

func (claims *Claims) ParseToken(req *http.Request) error {
	var token string
	cookie, err := req.Cookie("token")
	if err == nil {
		token, _ = url.QueryUnescape(cookie.Value)
	}
	if token == "" {
		token = req.Header.Get("Authorization")
	}
	if token == "" {
		return errors.New("未登录")
	}
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return "secret", nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return errors.New("登录超时")
			}
			return nil
		}
	}
	return errors.New("未登录")
}

type FasthttpClaims struct {
	User *model.UserMainInfo
	jwt.StandardClaims
}

func (claims *FasthttpClaims) GenerateToken() (string, error) {
	now := time.Now().Unix()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: now + int64(24*time.Hour),
		IssuedAt:  now,
		Issuer:    "hoper",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString("secret")

	return token, err
}

func (claims *FasthttpClaims) ParseToken(req *fasthttp.Request) error {
	var token string
	cookie := stringsi.ToString(req.Header.Cookie("token"))
	if len(cookie) > 0 {
		token, _ = url.QueryUnescape(cookie)
	}
	if token == "" {
		token = stringsi.ToString(req.Header.Peek("Authorization"))
	}
	if token == "" {
		return errors.New("未登录")
	}
	tokenClaims, _ := (&jwt.Parser{SkipClaimsValidation: true}).
		ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return "secret", nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if claims.VerifyExpiresAt(now, false) == false {
				return errors.New("登录超时")
			}
			return nil
		}
	}
	return errors.New("未登录")
}
