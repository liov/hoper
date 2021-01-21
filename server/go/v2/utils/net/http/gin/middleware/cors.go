package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	//method := c.Request.Method

	origin := ctx.Request.Header.Get("Origin")
	var headerKeys []string
	for k := range ctx.Request.Header {
		headerKeys = append(headerKeys, k)
	}
	headerStr := strings.Join(headerKeys, ", ")
	if headerStr != "" {
		headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
	} else {
		headerStr = "access-control-allow-origin, access-control-allow-headers"
	}
	if origin != "" {
		//下面的都是乱添加的-_-~

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", headerStr)
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-Control-Allow-Headers", "HeaderAuthorization, Content-Length, X-CSRF-HeaderCookieToken, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// c.Header("Access-Control-Max-Age", "172800")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Set("content-type", "application/json")
	}

	//放行所有OPTIONS方法
	if ctx.Request.Method == "OPTIONS" {
		ctx.JSON(http.StatusOK, "Options Request!")
	}
	ctx.Next()
}
