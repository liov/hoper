package gin_build

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/log"
	http2 "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/handlerconv"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/middleware"
)

func Http(ginHandle func(engine *gin.Engine)) http.Handler {

	//openapi
	r := gin.New()
	//r.Use(gin.Logger())
	logger := (&log.Config{Development: initialize.InitConfig.Env == initialize.PRODUCT}).NewLogger()
	middleware.SetLog(r, logger, false)
	r.Use(gin.Recovery())
	r.Any("/debug/*path", handlerconv.FromStd(http2.Debug()))
	if ginHandle != nil {
		ginHandle(r)
	}
	return r
}
