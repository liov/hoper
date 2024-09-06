package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/sdk/ffmpeg"
)

func Play(ctx *gin.Context) {
	file := ctx.Param("file")
	ffmpeg.Run("-re -i D:\\Download\\" + file + " -c:v libx265 -preset veryfast -crf 30 -maxrate 500k -bufsize 2000k -c:a aac -b:a 128k -f flv http://localhost:8090/live/stream")
}

func Stream(ctx *gin.Context) {
	fmt.Println(ctx)
}
