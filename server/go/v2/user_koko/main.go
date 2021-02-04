package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/oauth"
	"github.com/liov/hoper/go/v2/utils/net/http/koko"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(conf.Conf, dao.Dao)()
	pick.RegisterFiberService(service.GetUserService())
	app := fiber.New()
	pick.FiberWithCtx(app, service.FasthttpCtx, true, initialize.InitConfig.Module)
	go app.Listen(":3000")
	(&koko.Server{
		GRPCHandle: func(gs *grpc.Server)  {
			//grpc.OpenTracing = true
			model.RegisterUserServiceServer(gs, service.GetUserService())
			model.RegisterOauthServiceServer(gs, service.GetOauthService())
		},

		GinHandle: func(app *gin.Engine) {
			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.RegisterService(service.GetUserService())
			pick.Gin(app, model.ConvertContext, true, initialize.InitConfig.Module)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/

	}).Start(&koko.AuthInfoDao{
		Secret:    conf.Conf.Customize.TokenSecret,
		AuthCache: dao.Dao.Cache,
		AuthPool:  &sync.Pool{New: func() interface {}{return new(model.AuthInfo)}},
		Update:    service.Auth,
	})
}
