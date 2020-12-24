package handlerconv

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Convert(handlers []http.HandlerFunc) []gin.HandlerFunc{
	var rets []gin.HandlerFunc
	for _,handler:=range handlers{
		rets = append(rets,FromStd(handler))
	}
	return rets
}
