package service

import (
	"net/http"

	"github.com/liov/hoper/go/v2/httptpl/httptpl/internal/claims"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

type UserService struct{}

func (*UserService) Service() (string, string, http.HandlerFunc) {
	return "用户相关", "/api/${version}/user", nil
}

func (*UserService) Add(ctx *claims.Claims, req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("").
			Method(http.MethodPost).
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
	})

	return &model.SignupRep{Message: "测试"}, nil
}

func (*UserService) Edit(ctx *claims.Claims, req *model.EditReq) (*model.EditReq_EditDetails, error) {
	pick.Api(func() interface{} {
		return pick.Path("/id").
			Method(http.MethodGet).
			Title("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除")
	})

	return nil, nil
}

func (*UserService) Get(ctx *claims.Claims, req *model.GetReq) (*model.GetRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick.Api(func() interface{} {
		return pick.Path("/id/*name").
			Method(http.MethodGet).
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return &model.GetRep{Code: uint32(req.Id), Message: "测试"}, nil
}
