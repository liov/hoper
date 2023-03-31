package service

import (
	"context"
	"github.com/hopeio/pandora/context/http_context"
	"github.com/liov/hoper/server/go/mod/protobuf/upload"
	"github.com/liov/hoper/server/go/mod/upload/dao"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	uploadInfos, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	uploadInfos, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
