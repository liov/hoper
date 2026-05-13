package api

import (
	"github.com/gin-gonic/gin"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/service"
)

func GinRegister(app *gin.Engine) {
	app.GET("/video/*file", service.Video)
	app.GET("/live/stream", service.Play)
	pb.RegisterRemoteBrowseServiceHandlerServer(app, service.GetRemoteBrowseService())
	app.GET("/rb/signal", remoteBrowseSignal)
}
