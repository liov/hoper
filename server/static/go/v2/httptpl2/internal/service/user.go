package service

import (
	"net/http"

	"github.com/kataras/iris/v12"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
)

type UserService struct{}

func (*UserService) Name() string {
	return "user"
}

func (*UserService) DocTag() string {
	return "用户相关"
}

func (*UserService) Middle() []iris.Handler {
	return nil
}

func (*UserService) Add(req *model.SignupReq) (*model.SignupRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	api.Api(func() interface{} {
		return api.Method(http.MethodPut).
			Describe("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return nil, nil
}

func (*UserService) Edit(req *model.EditReq) (*model.EditReq_EditDetails, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	api.Api(func() interface{} {
		return api.Method(http.MethodPut).
			Describe("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
	})

	return nil, nil
}
