package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/httptpl/httptpl1/internal/grpcclient"
	"github.com/liov/hoper/go/v2/httptpl/httptpl1/internal/service"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
)

type UserController struct {
	Service *service.UserService
}

func (u *UserController) Middle() []iris.Handler {
	return nil
}

func (u *UserController) Describe() string {
	return "用户相关接口"
}

func (u *UserController) VerificationCode() {
	api.Path("/user/verification").
		Method(http.MethodPost).
		//这些信息不应该放在源代码里,都会打包进二进制文件
		Service(u.Service.VerificationCode).
		Title("获取验证码").
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
	api.Path("/user").
		Method(http.MethodPost).
		Service(u.Service.Add).
		Title("新增用户").
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
	api.Path("/user/{id:uint64}").
		Method(http.MethodGet).
		Service(u.Service.Add).
		Title("获取用户信息").
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
	api.Path("/user/{id:uint64}").
		Method(http.MethodPut).
		Service(u.Service.Add).
		Title("用户编辑").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
}

func (u *UserController) Test() {
	api.Path("/test").
		Method(http.MethodGet).
		Service(u.Service.Add).
		Title("测试").
		Version(1).
		CreateLog("1.0.0", "jyb", "2019/12/16", "创建")
}

func (u *UserController) Shutdown() {
	api.Path("/restart").
		Method(http.MethodGet).
		Title("系统重启").
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
