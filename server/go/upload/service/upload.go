package service

import (
	"context"
	"github.com/hopeio/tiga/context/http_context"
	gormi "github.com/hopeio/tiga/utils/dao/db/gorm"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"github.com/liov/hoper/server/go/upload/confdao"
	"github.com/liov/hoper/server/go/upload/data"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	uploadInfos, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := data.GetDao(ctxi)

	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctxi.TraceID)
	uploadInfos, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
