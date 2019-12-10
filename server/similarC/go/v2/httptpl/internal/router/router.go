package router

import (
	"reflect"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/liov/hoper/go/v2/httptpl/internal/controller"
)

func route(app *iris.Application) {
	//params:控制器,控制器公共的url路径，公共中间件
	newController := func(c interface{}, handlers ...context.Handler) interface{} {
		value := reflect.ValueOf(c)
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}
		value = value.Elem()
		if value.NumField() > 1 {
			panic("传入controller不合法")
		}
		handler := controller.Handler{ApiInfo: &controller.ApiInfo{}, App: app}
		value.Field(0).Set(reflect.ValueOf(controller.Controller{Handler: &handler, Middle: handlers}))
		return c
	}
	//controller注册
	ctrl := []interface{}{
		newController(&controller.UserController{}),
	}
	register(ctrl)
}

func register(ctrl []interface{}) {
	for i := range ctrl {
		value := reflect.ValueOf(ctrl[i])
		for j := 0; j < value.NumMethod(); j++ {
			value.Method(j).Call(nil)
		}
		c := value.Elem().Field(0).Interface().(controller.Controller)
		c.Middle = nil
		if i == len(ctrl) - 1 {
			c.ApiInfo = nil
		}
	}
}
