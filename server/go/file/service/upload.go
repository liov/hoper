package service

import (
	"context"
	"github.com/hopeio/context/httpctx"
	gormi "github.com/hopeio/utils/dao/database/gorm"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/protobuf/upload"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi, _ := httpctx.FromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	files, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{Files: files}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi, _ := httpctx.FromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	files, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{Files: files}, nil
}
