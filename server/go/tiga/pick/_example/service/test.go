package service

import (
	contexti "github.com/liov/hoper/v2/tiga/context"
	"github.com/liov/hoper/v2/tiga/pick"
	"net/http"

	model "github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/protobuf/utils/response"
)

type TestService struct{}

func (*TestService) Service() (string, string, []http.HandlerFunc) {
	return "测试相关", "/api/test", nil
}

func (*TestService) Test(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test1(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test2(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/a/").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test3(ctx *contexti.Ctx, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/a/:b").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}
