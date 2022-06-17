package handler

import (
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/gin-gonic/gin"
)

func Response(ctx *gin.Context, res ...interface{}) {
	httpi.Resp(ctx.Writer, res...)
}

func Res(ctx *gin.Context, code errorcode.ErrCode, msg string, data interface{}) {
	httpi.Response(ctx.Writer, code, msg, data)
}
