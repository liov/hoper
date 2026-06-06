package model

import (
	"strconv"

"github.com/hopeio/scaffold/context"
	"github.com/liov/hoper/server/go/protobuf/user"
)

type ClientInfo struct {
	Auth   *AuthInfo
	Device *context.DeviceInfo
}

type AuthInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Role user.Role   `json:"role"`
}

func (x *AuthInfo) GetId() string {
	return strconv.FormatUint(x.Id, 10)
}

func (x *AuthInfo) Proto() *user.Auth {
	return &user.Auth{
		Id:   x.Id,
		Name: x.Name,
		Role: x.Role,
	}
}

func ConvDeviceInfo(x *context.DeviceInfo) *user.AccessDevice {
	return &user.AccessDevice{
		Device:    x.Device,
		OS:        x.OS,
		AppCode:   x.AppCode,
		AppVer:    x.AppVer,
		IP:        x.IP.String(),
		Lng:       x.Lng,
		Lat:       x.Lat,
		Area:      x.Area,
		UserAgent: x.UserAgent,
	}
}
