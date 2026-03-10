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

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	uploadDao := data.GetDao(db)

	files, err := uploadDao.GetUrls(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
func (*FileService) GetUrlsByStrId(ctx context.Context, req *file.GetUrlsByStrIdReq) (*file.GetUrlsResp, error) {

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	uploadDao := data.GetDao(db)
	files, err := uploadDao.GetUrlsByStrId(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
