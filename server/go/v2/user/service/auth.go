package service

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/utils/net/http/auth"
	"github.com/liov/hoper/go/v2/utils/net/http/auth/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(r *http.Request) (*model.UserMainInfo, error) {
	token := auth.GetToken(r)
	if token == "" {
		return nil, model.UserErr_NoLogin
	}
	claims, err := jwt.ParseToken(token, conf.Conf.Customize.TokenSecret)
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

func (*UserService) AuthClaims(ctx context.Context) (*jwt.Claims, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("auth")
	if len(tokens) == 0 || tokens[0] == "" {
		return nil, model.UserErr_NoLogin
	}
	return jwt.ParseToken(tokens[0], conf.Conf.Customize.TokenSecret)
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
	if len(deviceInfo)>0 && deviceInfo[0] != "" {
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
	if len(location)>0 && location[0] != "" {
		info.Area, _ = url.PathUnescape(location[0])
	}

	if len(location)>1 && location[1] != "" {
		infos := strings.Split(location[1], ",")
		if len(infos) == 2 {
			info.Lng = infos[0]
			info.Lat = infos[1]
		}
	}

	userAgent := md.Get("user-agent")
	if len(userAgent)>0 && userAgent[0] != "" {
		info.UserAgent = userAgent[0]
	}
	ip := md.Get("x-forwarded-for")
	if len(ip)>0 && ip[0] != "" {
		info.IP = ip[0]
	}
	return &info
}

type authKey struct{}

// NewContext returns a new Context that carries value u.
func NewContext(r *http.Request) context.Context {
	user, _ := Auth(r)
	return context.WithValue(r.Context(), authKey{}, user)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*model.UserMainInfo, bool) {
	user, ok := ctx.Value(authKey{}).(*model.UserMainInfo)
	return user, ok
}
