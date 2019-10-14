package router

import (
	"reflect"

	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/gateway/internal/controller"
)

func route(app *iris.Application) {
	ctrl := []interface{}{
		&controller.UserController{},
	}
	reg(app,ctrl)
}


func reg(app *iris.Application, ctrl []interface{}) {
	appV := reflect.ValueOf(app)
	apiInfo:=&controller.ApiInfo{}
	infoV:= reflect.ValueOf(apiInfo)
	for _,c:=range ctrl{
		value := reflect.ValueOf(c)
		value.Elem().FieldByName("App").Set(appV)
		value.Elem().FieldByName("ApiInfo").Set(infoV)
		for i := 0; i < value.NumMethod(); i++ {
			value.Method(i).Call(nil)
		}
	}
}