package user

import (
	"time"
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
	now:=time.Now().Unix()
	if x.ExpiredAt!=0 && now <= x.ExpiredAt {
		return UserErr_LoginTimeout
	}
	return nil
}