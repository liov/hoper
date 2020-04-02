package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/httptpl/internal/grpcclient"
	"github.com/liov/hoper/go/v2/httptpl/internal/service"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"

	"github.com/liov/hoper/go/v2/utils/log"
)

type UserController struct {
	Handler *api.Handler
	Service *service.UserService
}

func (u *UserController) Middle() []iris.Handler {
	return nil
}

func (u *UserController) Name() string {
	return "user"
}

func (u *UserController) VerificationCode() {
	u.Handler.
		Path("/user/verification").
		Method(http.MethodPost).
		//这些信息不应该放在源代码里,都会打包进二进制文件
		Service(u.Service.VerificationCode).
		Describe("获取验证码").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
		ChangeLog("1.0.1", "jyb", "2019/12/16", "修改").
		Handle(
			//不用反射做
			func(c iris.Context) {
				var req empty.Empty
				c.ReadJSON(&req)
				c.JSON(u.Service.VerificationCode(&req))
			})
}

func (u *UserController) Add() {
	u.Handler.
		Path("/user").
		Method(http.MethodPost).
		Service(u.Service.Add).
		Describe("新增用户").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
		Handle(
			func(c iris.Context) {
				var req model.SignupReq
				c.ReadJSON(&req)
				c.JSON(u.Service.Add(&req))
			})

}

func (u *UserController) Get() {
	u.Handler.
		Path("/user/{id:uint64}").
		Method(http.MethodGet).
		Request(new(model.GetReq)).
		Response(new(model.GetRep)).
		Describe("获取用户信息").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
		Handle(
			func(c iris.Context) {
				id := c.Params().GetUint64Default("id", 0)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				rep, err := grpcclient.UserClient.GetUser(ctx, &model.GetReq{Id: id})
				if err != nil {
					log.Errorf("could not greet: %v", err)
				}
				c.JSON(rep)
			})
}

func (u *UserController) Edit() {
	u.Handler.
		Path("/user/{id:uint64}").
		Method(http.MethodPut).
		Request(new(model.EditReq)).
		Response(new(model.LoginRep)).
		Describe("用户编辑").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
		Handle(func(c iris.Context) {})
}

func (u *UserController) Shutdown() {
	u.Handler.
		Path("/restart").
		Method(http.MethodGet).
		Describe("系统重启").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
		Handle(
			func(c iris.Context) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				c.Application().(*iris.Application).Shutdown(ctx)
				c.WriteString("重启了")
			})
}
