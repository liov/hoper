package service

import (
	"context"
	"github.com/liov/hoper/v2/protobuf/upload"
	contexti "github.com/liov/hoper/v2/tiga/context"
	"github.com/liov/hoper/v2/upload/dao"
)

type UploadService struct {
	upload.UnimplementedUploadServiceServer
}

func (*UploadService) GetUrls(ctx context.Context, req *upload.GetUrlsReq) (*upload.GetUrlsRep, error) {
	ctxi := contexti.CtxFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(dao.Dao.GORMDB)
	uploadInfos, err := uploadDao.GetUrls(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
func (*UploadService) GetUrlsByStrId(ctx context.Context, req *upload.GetUrlsByStrIdReq) (*upload.GetUrlsRep, error) {
	ctxi := contexti.CtxFromContext(ctx)

	uploadDao := dao.GetDao(ctxi)

	db := ctxi.NewDB(dao.Dao.GORMDB)
	uploadInfos, err := uploadDao.GetUrlsByStrId(db, req.Ids)
	if err != nil {
		return nil, err
	}
	return &upload.GetUrlsRep{UploadInfos: uploadInfos}, nil
}
