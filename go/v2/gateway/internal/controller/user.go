package controller

import (
	"net/http"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/liov/hoper/go/v2/gateway/internal/api/request"
	"github.com/liov/hoper/go/v2/gateway/internal/service"
)

type UserController struct{}
var userService = &service.UserService{}
func (u *UserController) Add(app *iris.Application) {
	apiInfo(
		path("/user"),
		method(http.MethodPost),
		describe("新增用户"),
		auth("jyb"),
		version(1),
		handle(app,
			func(ctx context.Context) {
				var req request.AddUserReq
				ctx.ReadJSON(&req)
				userService.Add(&req)
			}),
	)


}

func (u *UserController) Get(app *iris.Application) {
	apiInfo(
		path("/add/:id"),
		method(http.MethodGet),
		describe("get"),
		auth("jyb"),
		version(1),
		handle(app,
			func(ctx context.Context) {
				ctx.Writef("返回")
			}),
	)
}
