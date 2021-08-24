package service

import (
	contexti "github.com/liov/hoper/server/go/lib/tiga/context"
	"github.com/liov/hoper/server/go/lib/tiga/pick"
	"net/http"

	"github.com/liov/hoper/server/go/lib/protobuf/response"
	model "github.com/liov/hoper/server/go/lib/protobuf/user"
)

type TestService struct{}

func (*TestService) Service() (string, string, []http.HandlerFunc) {
	return "测试相关", "/api/test", nil
}

func (*TestService) Test(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test1(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("/").
			Title("用户注册").
			Version(1).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			End()
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test2(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("/a/").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test3(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() {
		pick.Post("/a/:b").
			Title("用户注册").
			Version(3).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}
