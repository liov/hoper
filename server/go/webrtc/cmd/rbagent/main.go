package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/liov/hoper/server/go/webrtc/session"
)

func main() {
	room := flag.String("room", "demo", "配对房间码")
	root := flag.String("root", ".", "本地目录或 rfv 根路径")
	signalURL := flag.String("signal", envOr("RB_SIGNAL_WS", "ws://127.0.0.1:8080/rb/signal"), "信令 WebSocket")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := session.RunAgent(ctx, *signalURL, *room, *root); err != nil {
		log.Fatalf("agent: %v", err)
	}
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
