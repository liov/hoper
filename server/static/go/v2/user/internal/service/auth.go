package service

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/utils/net/http/auth/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(r *http.Request) (*model.UserMainInfo, error) {
	var auth string
	cookie, _ := r.Cookie("token")
	value, _ := url.QueryUnescape(cookie.Value)
	if value == "" {
		auth = r.Header.Get("authorization")
	} else {
		auth = value
	}

	if auth == "" {
		return nil, model.UserErr_NoLogin
	}
	claims, err := jwt.ParseToken(auth, config.Conf.Customize.TokenSecret)
	if err != nil {
		return nil, model.UserErr_InvalidToken
	}
	conn := dao.NewUserRedis()
	defer conn.Close()
	user, err := conn.EfficientUserHashFromRedis(claims.UserID)
	if err != nil {
		log.Error(err)
		return nil, model.UserErr_InvalidToken
	}
	return user, nil
}

func (*UserService) Auth(ctx context.Context) (*model.UserMainInfo, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("auth")
	if len(tokens) == 0 || tokens[0] == "" {
		return nil, model.UserErr_NoLogin
	}
	claims, err := jwt.ParseToken(tokens[0], config.Conf.Customize.TokenSecret)
	if err != nil {
		return nil, err
	}
	//redis
	conn := dao.NewUserRedis()
	defer conn.Close()
	user, err := conn.EfficientUserHashFromRedis(claims.UserID)
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
	if deviceInfo[0] != "" {
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
	if location[0] != "" {
		info.Area, _ = url.PathUnescape(location[0])
	}

	if location[1] != "" {
		infos := strings.Split(location[1], ",")
		if len(infos) == 2 {
			info.Lng = infos[0]
			info.Lat = infos[1]
		}
	}

	userAgent := md.Get("user-agent")
	if userAgent[0] != "" {
		info.UserAgent = userAgent[0]
	}
	ip := md.Get("x-forwarded-for")
	if ip[0] != "" {
		info.IP = ip[0]
	}
	return &info
}
