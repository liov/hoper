package model

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/reconn"
	"google.golang.org/grpc"
)

//Cannot use 'resumes' (type []*model.Resume) as type []CmpKey
//我认为这是一个bug
//[]int可以是interface，却不可以是[]interface
//var test []array.CmpKey
//test = append(test,resumes[0]) 可行
//test = append(test,resumes...) 不可行，可笑
func (m *Resume) CmpKey() uint64 {
	return m.Id
}

var UserService_serviceDesc = &_UserService_serviceDesc

func (m *LoginRep) GetCookie() string {
	return m.Cookie
}

func (m *LogoutRep) GetCookie() string {
	return m.Cookie
}

func RegisterUserServiceHandlerFromModuleWithReConnect(ctx context.Context, mux *runtime.ServeMux, module string, opts []grpc.DialOption, reConnect map[string]func() error) (err error) {
	conn, err := grpc.Dial(initialize.BasicConfig.NacosConfig.GetServiceEndPort(module), opts...)
	if err != nil {
		return err
	}
	client := NewUserServiceClient(conn)
	reConnect[module] = reconn.ReConnect(client, module, opts)
	return RegisterUserServiceHandlerClient(ctx, mux, client)
}
