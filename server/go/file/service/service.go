package service

import (
	"context"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/database/sql/gorm"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/protobuf/file"
)

type FileService struct {
	file.UnimplementedFileServiceServer
}

func (*FileService) GetUrls(ctx context.Context, req *file.GetUrlsReq) (*file.GetUrlsResp, error) {
	ctxi, _ := httpctx.FromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gorm.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	files, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
func (*FileService) GetUrlsByStrId(ctx context.Context, req *file.GetUrlsByStrIdReq) (*file.GetUrlsResp, error) {
	ctxi, _ := httpctx.FromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gorm.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	files, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}
