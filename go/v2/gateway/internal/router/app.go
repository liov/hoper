package router

import (
	"context"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/api"
	"github.com/liov/hoper/go/v2/utils/log/mid"
)

func App() *iris.Application {
	app := iris.New()

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		//关闭所有主机
		app.Shutdown(ctx)
	})
	app.StaticWeb("/api/static", "../../static")
	app.Use(recover.New())
	if config.Conf.Env != initialize.PRODUCT {
		app.Use(api.ApiMiddle)
	}
	//other.Wrap(app)
	//api文档
	//other.Api(app)
	//https://rpm.newrelic.com/accounts/2269290/applications
	/*	config := newrelic.config("hoper", "199e00247f278548fe92d6c81aeaadac0fc52b4b")
		m, err := newrelic.New(config)
		if err != nil {
			app.Logger().Fatal(err)
		}
		app.Use(m.ServeHTTP)*/

	/*	prometheus := prometheus.New("hoper")
		app.Use(prometheus.ServeHTTP)
		app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
			//错误代码处理程序不与其他路由共享相同的中间件，所以单独执行错误
			prometheus.ServeHTTP(ctx)
			ctx.Writef("Not Found")
		})
	*/
	/*middleware必须要写ctx.next(),且写在路由前，路由后的midddleware在请求之前的路由时不生效
	  iris.FromStd()将其他Handler转为iris的Handler
	*/
	//i18n
	/*		globalLocale := i18n.New(i18n.config{
			Default:      "en-US",
			URLParameter: "lang",
			Languages: map[string]string{
				"en-US": "../../data/i18n/locale_en-US.ini",
				"zh-CN": "../../data/i18n/locale_zh-CN.ini"}})
		app.Use(globalLocale)*/
	//请求日志
	app.Use(mid.LogMid(false))

	app.OnAnyErrorCode(mid.LogMid(true), func(ctx iris.Context) {
		//这应该被添加到日志中，因为`logger.config＃MessageContextKey`
		ctx.Values().Set("logger_message",
			ctx.Request())
		ctx.Writef("My Custom error page")
	})


	route(app)

	return app
}


