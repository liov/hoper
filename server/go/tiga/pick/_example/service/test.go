package service

import (
	"context"
	"github.com/liov/hoper/v2/tiga/pick"
	"net/http"

	model "github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/protobuf/utils/response"
)

type TestService struct{}

func (*TestService) Service() (string, string, []http.HandlerFunc) {
	return "测试相关", "/api/test", nil
}

func (*TestService) Test(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Post("").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test1(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test2(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/a/").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*TestService) Test3(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
	pick.Api(func() interface{} {
		return pick.Post("/a/:b").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &response.TinyRep{Message: "测试"}, nil
}
