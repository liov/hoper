package gini

import (
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/net/http/gin/handler"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug(r *gin.Engine) {
	r.Any("/debug/*path", handler.Wrap(httpi.Debug()))
	// Register Prometheus metrics handler.
	r.Any("/metrics", handler.Wrap(promhttp.Handler()))
}
