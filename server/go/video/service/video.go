package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/sdk/ffmpeg"
	"github.com/notedit/rtmp-lib/flv"
	"github.com/notedit/rtmp-lib/pubsub"
	"io"
	"net/http"
	"sync"
)

func Video(ctx *gin.Context) {
	file := ctx.Param("file")
	ffmpeg.Run("-i D:\\Download\\" + file + "-an -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp rtp://127.0.0.1:5004?pkt_size=1200 -vn -c:a libopus -f rtp rtp://127.0.0.1:5006?pkt_size=1200")
}

type Channel struct {
	que *pubsub.Queue
}

var channels = map[string]*Channel{}

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (self writeFlusher) Flush() error {
	self.httpflusher.Flush()
	return nil
}

var l sync.RWMutex

func Play(ctx *gin.Context) {
	l.RLock()
	ch := channels[ctx.Request.URL.Path]
	l.RUnlock()

	if ch != nil {
		ctx.Header("Content-Type", "video/x-flv")
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Status(200)
		flusher := ctx.Writer.(http.Flusher)
		flusher.Flush()

		muxer := flv.NewMuxerWriteFlusher(writeFlusher{httpflusher: flusher, Writer: ctx.Writer})
		cursor := ch.que.Latest()

		streams, err := cursor.Streams()

		if err != nil {
			panic(err)
		}

		muxer.WriteHeader(streams)

		for {
			packet, err := cursor.ReadPacket()
			if err != nil {
				break
			}
			muxer.WritePacket(packet)
		}
	} else {
		http.NotFound(ctx.Writer, ctx.Request)
	}
}
