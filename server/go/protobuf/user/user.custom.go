package user

import (
	"strconv"

	"github.com/hopeio/scaffold/context"
)

var UserserviceServicedesc = &UserService_ServiceDesc

/*----------------------------ORM-------------------------------*/

/*----------------------------AuthInfo-------------------------------*/

type ClientInfo struct {
	Auth   *AuthInfo
	Device *context.DeviceInfo
}

type AuthInfo struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Role Role   `json:"role"`
}

func (x *AuthInfo) GetId() string {
	return strconv.FormatUint(x.Id, 10)
}

func (x *AuthInfo) Proto() *Auth {
	return &Auth{
		Id:   x.Id,
		Name: x.Name,
		Role: x.Role,
	}
}

func ConvDeviceInfo(x *context.DeviceInfo) *AccessDevice {
	return &AccessDevice{
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
