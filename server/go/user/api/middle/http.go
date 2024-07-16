package middle

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/hopeio/utils/log"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RequestURI)
}

func GinLog(ctx *gin.Context) {
	log.Debug(ctx.Request.RequestURI)
}
