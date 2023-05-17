package service

import (
	"context"
	"github.com/actliboy/hoper/server/go/protobuf/upload"
	"github.com/actliboy/hoper/server/go/upload/confdao"
	"github.com/actliboy/hoper/server/go/upload/dao"
	"github.com/hopeio/pandora/context/http_context"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	uploadInfos, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi := http_context.ContextFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	uploadInfos, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
