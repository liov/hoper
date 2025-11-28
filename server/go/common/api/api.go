package api

import (
	"github.com/gin-gonic/gin"
	commonService "github.com/liov/hoper/server/go/common/service"
	"github.com/liov/hoper/server/go/protobuf/common"
)

func GinRegister(app *gin.Engine) {
	common.RegisterCommonServiceHandlerServer(app, commonService.GetCommonService())
}
