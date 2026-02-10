package service

import (
	"context"

	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/protobuf/file"
)

type FileService struct {
	file.UnimplementedFileServiceServer
}

func (*FileService) GetUrls(ctx context.Context, req *file.GetUrlsReq) (*file.GetUrlsResp, error) {
	ctx, span := Tracer.Start(ctx, "FileService.GetUrls")
	defer span.End()

	uploadDao := data.GetDao(ctx, global.Dao.GORMDB.DB)

	files, err := uploadDao.GetUrls(req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
func (*FileService) GetUrlsByStrId(ctx context.Context, req *file.GetUrlsByStrIdReq) (*file.GetUrlsResp, error) {
	ctx, span := Tracer.Start(ctx, "FileService.GetUrlsByStrId")
	defer span.End()
	uploadDao := data.GetDao(ctx, global.Dao.GORMDB.DB)
	files, err := uploadDao.GetUrlsByStrId(req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
