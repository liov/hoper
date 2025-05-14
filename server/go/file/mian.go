package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	uploadapi "github.com/liov/hoper/server/go/file/api"
	"github.com/liov/hoper/server/go/global"
	uconf "github.com/liov/hoper/server/go/user/global"
)

func main() {
	//配置初始化应该在第一位
	defer global.Global.Cleanup()
	uconf.Conf.Server.WithOptions(func(s *cherry.Server) {
		s.GinHandler = func(app *gin.Engine) {
			uploadapi.GinRegister(app)
		}
		//s.GraphqlHandler= graphql.NewExecutableSchema(),
	}).Run()

}
