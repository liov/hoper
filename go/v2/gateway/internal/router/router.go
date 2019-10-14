package router

import (
	"reflect"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/liov/hoper/go/v2/gateway/internal/controller"
)

func route(app *iris.Application) {
	getController :=func(partyPath string, handlers ...context.Handler) controller.Controller {
		return controller.Controller{App:app.Party("/api"+partyPath,handlers...)}
	}

	ctrl := []interface{}{
		&controller.UserController{Controller: getController("/user")},
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