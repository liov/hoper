package upload

import (
	"github.com/gin-gonic/gin"
)

func ExistsGin(ctx *gin.Context) {
	md5 := ctx.Param("md5")
	size := ctx.Param("size")
	exists(ctx.Request.Context(), ctx.Writer, md5, size)
}
