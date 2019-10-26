package router

import (
	"reflect"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/liov/hoper/go/v2/gateway/internal/controller"
)

func route(app *iris.Application) {
	//params:控制器,控制器公共的url路径，公共中间件
	newController :=func(c interface{},handlers ...context.Handler) interface{} {
		value:=reflect.ValueOf(c)
		if value.Kind() != reflect.Ptr{
			panic("必须传入指针")
		}
		value = value.Elem()
		if value.NumField() > 1{
			panic("传入controller不合法")
		}
		value.Field(0).Set(reflect.ValueOf(controller.Controller{App:app,Middle:handlers}))
		return c
	}
	//controller注册
	ctrl := []interface{}{
		newController(&controller.UserController{}),
	}
	register(ctrl)
}




func register(ctrl []interface{}) {
	for _,c:=range ctrl{
		value := reflect.ValueOf(c)
		for i := 0; i < value.NumMethod(); i++ {
			value.Method(i).Call(nil)
		}
	}
}