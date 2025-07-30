package service

import (
	"context"
	"github.com/hopeio/gox/sugar"
	"github.com/liov/hoper/server/go/file/api/request"
	"github.com/liov/hoper/server/go/file/api/response"
	"os"
)

func (*FileService) ListDir(ctx context.Context, req *request.ListDir) (*response.ListDir, error) {
	info, err := os.Stat(req.Dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrNotExist
	}
	dirs, err := os.ReadDir(req.Dir)
	if err != nil {
		return nil, err
	}
	var items []*response.Entry
	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			return nil, err
		}
		items = append(items, &response.Entry{
			Name: dir.Name(),
			Type: sugar.TernaryOperator(dir.IsDir(), 0, 1),
			Size: info.Size(),
		})
	}
	return nil, nil
}

func (*FileService) Thumbnail(ctx context.Context, req *request.Path) (*response.ListDir, error) {
	return nil, nil
}
