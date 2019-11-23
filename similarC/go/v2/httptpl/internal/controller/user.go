package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/httptpl/internal/grpcclient"
	model "github.com/liov/hoper/go/v2/protobuf/user"

	"github.com/liov/hoper/go/v2/utils/log"
)

type UserController struct{
	Controller
}

func (u *UserController) VerificationCode() {
	u.api(
		path("/user/verificationCode"),
		method(http.MethodPost),
		describe("获取验证码"),
		auth("jyb"),
		version(1),
		handle(
			func(c iris.Context) {
				var req model.VerifyReq
				c.ReadJSON(&req)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				rep,err:= grpcclient.Client.UserClient.Verify(ctx,&req)
				if err != nil {
					log.Errorf("could not greet: %v", err)
				}
				c.JSON(rep)
			}),
	)

}


func (u *UserController) Add() {
	u.api(
		path("/user"),
		method(http.MethodPost),
		describe("新增用户"),
		auth("jyb"),
		version(1),
		handle(
			func(c iris.Context) {
				var req model.SignupReq
				c.ReadJSON(&req)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				rep,err:= grpcclient.Client.UserClient.Signup(ctx,&req)
				if err != nil {
					log.Errorf("could not greet: %v", err)
				}
				c.JSON(rep)
			}),
	)

}

func (u *UserController) Get() {
	u.api(
		path("/user/:id"),
		method(http.MethodGet),
		describe("获取用户信息"),
		auth("jyb"),
		version(1),
		handle(
			func(c iris.Context) {
				id:= c.Params().GetUint64Default("id",0)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				rep,err:= grpcclient.Client.UserClient.GetUser(ctx,&model.GetReq{ID: id})
				if err != nil {
					log.Errorf("could not greet: %v", err)
				}
				c.JSON(rep)
			}),
	)
}

func (u *UserController) Shutdown() {
	u.api(
		path("/shutdown"),
		method(http.MethodGet),
		describe("get"),
		auth("jyb"),
		version(1),
		handle(
			func(c iris.Context) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				u.App.Shutdown(ctx)
				c.WriteString("重启了")
			}),
	)
}