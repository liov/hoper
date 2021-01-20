package user

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/liov/hoper/go/v2/utils/net/http/request"
)

//Cannot use 'resumes' (type []*model.Resume) as type []CmpKey
//我认为这是一个bug
//[]int可以是interface，却不可以是[]interface
//var test []array.CmpKey
//test = append(test,resumes[0]) 可行
//test = append(test,resumes...) 不可行，可笑
func (x *Resume) CmpKey() uint64 {
	return x.Id
}

var UserserviceServicedesc = &_UserService_serviceDesc

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
func (x *User) TableName() string {
	if x.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(x.Id/1_000_000+49))
}

func (x *UserBaseInfo) TableName() string {
	if x.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(x.Id/1_000_000+49))
}

func (x *UserAuthInfo) TableName() string {
	if x.Id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(x.Id/1_000_000+49))
}

func (x *Resume) TableName() string {
	return "resume"
}

func (x *UserAuthInfo) Valid() error {
	now := time.Now().Unix()
	if x.ExpiredAt != 0 && now <= x.ExpiredAt {
		return UserErr_LoginTimeout
	}
	return nil
}

func Device(r http.Header) *DeviceInfo {
	unknow := true
	var info DeviceInfo
	//Device-Info:device-osInfo-appCode-appVersion
	if deviceInfo := r.Get(request.DeviceInfo); deviceInfo != "" {
		unknow = false
		infos := strings.Split(deviceInfo, "-")
		if len(infos) == 4 {
			info.Device = infos[0]
			info.Os = infos[1]
			info.AppCode = infos[2]
			info.AppVersion = infos[3]
		}
	}
	// area:xxx
	// location:1.23456,2.123456
	if area := r.Get(request.Area); area != "" {
		unknow = false
		info.Area, _ = url.PathUnescape(area)
	}
	if location := r.Get(request.Location); location != "" {
		unknow = false
		infos := strings.Split(location, ",")
		if len(infos) == 2 {
			info.Lng = infos[0]
			info.Lat = infos[1]
		}
	}

	if userAgent := r.Get(request.UserAgent); userAgent != "" {
		unknow = false
		info.UserAgent = userAgent
	}
	if ip := r.Get(request.XForwardedFor); ip != "" {
		unknow = false
		info.IP = ip
	}
	if unknow {
		return nil
	}
	return &info
}

func (x *AuthInfo) UserAuthInfo() *UserAuthInfo {
	return &UserAuthInfo{
		Id:           x.Id,
		Name:         x.Name,
		Role:         x.Role,
		Status:       x.Status,
		LastActiveAt: x.LastActiveAt,
		ExpiredAt:    x.ExpiredAt,
		LoginAt:      x.LoginAt,
	}
}

type AuthInfo struct {
	Id           uint64     `json:"id"`
	Name         string     `json:"name"`
	Role         Role       `json:"role"`
	Status       UserStatus `json:"status"`
	LastActiveAt int64      `json:"lastActiveAt,omitempty"`
	ExpiredAt    int64      `json:"expiredAt,omitempty"`
	LoginAt      int64      `json:"loginAt,omitempty"`
}

type DeviceInfo struct {
	//设备
	Device     string `json:"device" gorm:"size:255"`
	Os         string `json:"os" gorm:"size:255"`
	AppCode    string `json:"appCode" gorm:"size:255"`
	AppVersion string `json:"appVersion" gorm:"size:255"`
	IP         string `json:"IP" gorm:"size:255"`
	Lng        string `json:"lng" gorm:"type:numeric(10,6)"`
	Lat        string `json:"lat" gorm:"type:numeric(10,6)"`
	Area       string `json:"area" gorm:"size:255"`
	UserAgent  string `json:"userAgent" gorm:"size:255"`
}

func (x *DeviceInfo) UserDeviceInfo() *UserDeviceInfo {
	return &UserDeviceInfo{
		Device:     x.Device,
		Os:         x.Os,
		AppCode:    x.AppCode,
		AppVersion: x.AppVersion,
		IP:         x.IP,
		Lng:        x.Lng,
		Lat:        x.Lat,
		Area:       x.Area,
		UserAgent:  x.UserAgent,
	}
}
