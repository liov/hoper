package service

import (
	pick2 "github.com/liov/hoper/server/go/lib/pick"
	"github.com/liov/hoper/server/go/lib/pick/_example/middle"
	"net/http"
)

type UserService struct{}

func (*UserService) Service() (string, string, []http.HandlerFunc) {
	return "用户相关", "/api/user", []http.HandlerFunc{middle.Log}
}

func (*UserService) Add(ctx *contexti.Ctx, req *SignupReq) (*TinyRep, error) {
	//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
	pick2.Api(func() {
		pick2.Post("").
			Title("用户注册").
			Version(2).
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
	})

	return &TinyRep{Message: "测试"}, nil
}

type EditReq struct {
}
type EditReq_EditDetails struct {
}

func (*UserService) Edit(ctx *contexti.Ctx, req *EditReq) (*EditReq_EditDetails, error) {
	pick2.Api(func() {
		pick2.Put("/:id").
			Title("用户编辑").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
			Deprecated("1.0.0", "jyb", "2019/12/16", "删除").End()
	})

	return nil, nil
}

type Object struct {
	Id uint64 `json:"id"`
}

func (*UserService) Get(ctx *contexti.Ctx, req *Object) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Get("/:id").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Code: uint32(req.Id), Message: "测试"}, nil
}

type StaticService struct{}

func (*StaticService) Service() (string, string, []http.HandlerFunc) {
	return "静态资源", "/api/static", nil
}

type TinyRep struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

func (*StaticService) Get2(ctx *contexti.Ctx, req *SignupReq) (*TinyRep, error) {
	pick2.Api(func() {
		pick2.Get("/*mail").
			Title("用户注册").
			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").End()
	})

	return &TinyRep{Message: req.Mail}, nil
}
