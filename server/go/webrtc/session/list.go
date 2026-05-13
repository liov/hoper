package session

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/webrtc/rfvclient"
)

func listEntries(ctx context.Context, root string) ([]*pb.FileEntry, error) {
	if useRFV() {
		return rfvclient.ListRemoteFiles(ctx, root)
	}
	return listLocalDir(root)
}

func useRFV() bool {
	return strings.TrimSpace(os.Getenv("RB_AGENT_USE_RFV")) != "0"
}

func listLocalDir(root string) ([]*pb.FileEntry, error) {
	ents, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]*pb.FileEntry, 0, len(ents))
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		name := e.Name()
		abs := filepath.Join(root, name)
		out = append(out, &pb.FileEntry{
			Id: abs, Name: name, Size: info.Size(), MtimeUnixMs: info.ModTime().UnixMilli(),
		})
	}
	return out, nil
}
