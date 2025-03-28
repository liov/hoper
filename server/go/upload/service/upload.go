package service

import (
	"context"
	"github.com/hopeio/context/httpctx"
	gormi "github.com/hopeio/utils/dao/database/gorm"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"github.com/liov/hoper/server/go/upload/data"
	"github.com/liov/hoper/server/go/upload/global"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi, _ := httpctx.FromContextValue(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	uploadInfos, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi, _ := httpctx.FromContextValue(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	uploadInfos, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
