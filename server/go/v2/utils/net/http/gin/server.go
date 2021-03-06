package gin_build

import (
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/verification/validator"
)

func Http(confPath string,ginHandle func(engine *gin.Engine)) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	validator.DefaultValidator = nil // 自己做校验

	//openapi
	r := gin.New()
	WithConfiguration(r, confPath)
	//r.Use(gin.Logger())
	/*logger := (&log.Config{Development: initialize.InitConfig.Env == initialize.PRODUCT}).NewLogger()
	middleware.SetLog(r, logger, false)*/
	//r.Use(gin.Recovery())
	Debug(r)
	// r.Any("/*any", handler.FromStd(http.DefaultServeMux))
	if ginHandle != nil {
		ginHandle(r)
	}
	return r
}
