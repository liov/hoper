package api

import (
	"github.com/gin-gonic/gin"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/file/remotebrowse"
)

func RegisterRemoteBrowse(app *gin.Engine) {
	pb.RegisterRemoteBrowseServiceHandlerServer(app, remotebrowse.GetService())
	app.GET("/rb/signal", func(c *gin.Context) {
		proxyRBSignal(c, remotebrowse.SignalUpstream())
	})
}
