package gini

import (
	"github.com/liov/hoper/server/go/lib/utils/verification/validator"
	"github.com/gin-gonic/gin"
)

func Http(conf *Config, ginHandle func(engine *gin.Engine)) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	validator.DefaultValidator = nil // 自己做校验

	//openapi
	r := gin.New()
	conf.SetConfig(r)
	//r.Use(gin.Logger())
	/*logger := (&log.Config{Development: initialize.SpecialConfig.Env == initialize.PRODUCT}).NewLogger()
	middleware.SetLog(r, logger, false)*/
	//r.Use(gin.Recovery())
	Debug(r)
	// r.Any("/*any", handler.FromStd(http.DefaultServeMux))
	if ginHandle != nil {
		ginHandle(r)
	}
	return r
}
