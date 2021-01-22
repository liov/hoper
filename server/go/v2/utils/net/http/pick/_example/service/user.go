package service

import (
	"net/http"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/http/pick/_example/middle"
)

type UserService struct{}

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middle.Log}
}

func (*UserService) Add(ctx *Claims, req *model.SignupReq) (*response.TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Method(http.MethodPost).
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})

	return &response.TinyRep{Message: "测试"}, nil
}

func (*UserService) Edit(ctx *Claims, req *model.EditReq) (*model.EditReq_EditDetails, error) {
	pick.Api(func() interface{} {
		return pick.Path("/:id").
			Method(http.MethodPut).
			Title("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除")
	})

	return nil, nil
}

func (*UserService) Get(ctx *Claims, req *model.GetReq) (*model.GetRep, error) {
	pick.Api(func() interface{} {
		return pick.Path("/:id").
			Method(http.MethodGet).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.GetRep{Code: uint32(req.Id), Message: "测试"}, nil
}

type StaticService struct{}

func (*StaticService) Service() (string, string, []http.HandlerFunc) {
	return "静态资源", "/api/static", nil
}

func (*StaticService) Get2(ctx *Claims, req *model.SignupReq) (*model.GetRep, error) {
	pick.Api(func() interface{} {
		return pick.Path("/*mail").
			Method(http.MethodGet).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.GetRep{Message: req.Mail}, nil
}
