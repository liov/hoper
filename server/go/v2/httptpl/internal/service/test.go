package service

import (
	"net/http"

	"github.com/liov/hoper/go/v2/httptpl/internal/claims"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

type TestService struct{}

func (*TestService) Service() (string, string, []http.HandlerFunc) {
	return "测试相关", "/api/${version}/test", nil
}

func (*TestService) Test(ctx *claims.Claims, req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("").
			Method(http.MethodPost).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.SignupRep{Message: "测试"}, nil
}

func (*TestService) Test1(ctx *claims.Claims, req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("/").
			Method(http.MethodPost).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.SignupRep{Message: "测试"}, nil
}

func (*TestService) Test2(ctx *claims.Claims, req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("/a/").
			Method(http.MethodPost).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.SignupRep{Message: "测试"}, nil
}

func (*TestService) Test3(ctx *claims.Claims, req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("/a/:b").
			Method(http.MethodPost).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.SignupRep{Message: "测试"}, nil
}