package user

import (
	contexti "github.com/hopeio/context/reqctx"
	"strconv"
)

// Cannot use 'resumes' (type []*model.Resume) as type []CmpKey
// 我认为这是一个bug
// []int可以是interface，却不可以是[]interface
// var test []array.CmpKey
// test = append(test,resumes[0]) 可行
// test = append(test,resumes...) 不可行，可笑
func (x *Resume) CompareKey() uint64 {
	return x.Id
}

var UserserviceServicedesc = &UserService_ServiceDesc

/*
func RegisterUserServiceHandlerFromModuleWithReConnect(ctx context.Context, mux *runtime.ServeMux, getEndPort func() string, opts []grpc.DialOption) (err error) {
	endPort:=getEndPort()
	conn, err := grpc.Dial(endPort, opts...)
	if err != nil {
		return err
	}
	client := NewUserServiceClient(conn)
	reconn.ReConnectMap[endPort] = reconn.ReConnect(client, getEndPort, opts)
	return RegisterUserServiceHandlerClient(ctx, mux, client)
}
*/

/*----------------------------ORM-------------------------------*/

/*----------------------------AuthInfo-------------------------------*/

type AuthInfo struct {
	Id     uint64     `json:"id"`
	Name   string     `json:"name"`
	Role   Role       `json:"role"`
	Status UserStatus `json:"status"`
}

func (x *AuthInfo) IdStr() string {
	return strconv.FormatUint(x.Id, 10)
}

func (x *AuthInfo) UserAuthInfo() *AuthInfoRep {
	return &AuthInfoRep{
		Id:     x.Id,
		Name:   x.Name,
		Role:   x.Role,
		Status: x.Status,
	}
}

func ConvDeviceInfo(x *contexti.DeviceInfo) *AccessDevice {
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
