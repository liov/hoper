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
	root := flag.String("path", "", "列举目录，空则仅测连通")
	signalURL := flag.String("signal", envOr("RB_SIGNAL_WS", "ws://127.0.0.1:8080/rb/signal"), "信令 WebSocket")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	conn, err := session.ConnectViewerRelay(ctx, *signalURL, *room)
	if err != nil {
		log.Fatalf("relay: %v", err)
	}
	defer conn.Close()
	entries, err := session.ListFilesOn(conn, *root)
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	log.Printf("files=%d", len(entries))
	for _, e := range entries {
		log.Printf("%s size=%d", e.GetName(), e.GetSize())
	}
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
