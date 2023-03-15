package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
)

func RespErrcode(ctx *gin.Context, code errorcode.ErrCode) {
	httpi.RespErrcode(ctx.Writer, code)
}

func RespErr(ctx *gin.Context, err error) {
	httpi.RespErr(ctx.Writer, err)
}

func RespErrMsg(ctx *gin.Context, msg string) {
	httpi.RespErrMsg(ctx.Writer, msg)
}

func RespErrRep(ctx *gin.Context, rep *errorcode.ErrRep) {
	httpi.RespErrRep(ctx.Writer, rep)
}

func Response(ctx *gin.Context, code errorcode.ErrCode, msg string, data interface{}) {
	httpi.Response(ctx.Writer, code, msg, data)
}
