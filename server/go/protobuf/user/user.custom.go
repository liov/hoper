package user

import (
	"strconv"

	"github.com/hopeio/gox/context/reqctx"
)

var UserserviceServicedesc = &UserService_ServiceDesc

/*----------------------------ORM-------------------------------*/

/*----------------------------AuthInfo-------------------------------*/

type AuthBase struct {
	Id     uint64     `json:"id"`
	Name   string     `json:"name"`
	Role   Role       `json:"role"`
	Status UserStatus `json:"status"`
	Avatar string     `json:"avatar"`
}

func (x *AuthBase) GetId() string {
	return strconv.FormatUint(x.Id, 10)
}

func (x *AuthBase) Proto() *Auth {
	return &Auth{
		Id:     x.Id,
		Name:   x.Name,
		Role:   x.Role,
		Status: x.Status,
	}
}

func ConvDeviceInfo(x *reqctx.DeviceInfo) *AccessDevice {
	return &AccessDevice{
		Device:    x.Device,
		OS:        x.OS,
		AppCode:   x.AppCode,
		AppVer:    x.AppVer,
		IP:        x.IP,
		Lng:       x.Lng,
		Lat:       x.Lat,
		Area:      x.Area,
		UserAgent: x.UserAgent,
	}
}
