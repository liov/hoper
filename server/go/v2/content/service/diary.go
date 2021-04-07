package service

import (
	"context"
	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type DiaryService struct {
	content.UnimplementedDiaryServiceServer
}

func (m *DiaryService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "日记相关", "/api/diary", nil
}

func (*DiaryService) DiaryBook(ctx context.Context, req *content.DiaryBookReq) (*content.DiaryBookRep, error) {

	return nil, status.Errorf(codes.Unimplemented, "method DiaryBook not implemented")
}
func (*DiaryService) DiaryBookList(context.Context, *content.DiaryBookListReq) (*content.DiaryBookListRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiaryBookList not implemented")
}
func (*DiaryService) AddDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*request.Object, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	req.UserId = auth.Id
	err = db.Table(model.DiaryBookTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return &request.Object{Id:req.Id}, nil
}
func (*DiaryService) EditDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	req.UserId = auth.Id
	err = db.Table(model.DiaryBookTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return nil, nil
}
func (*DiaryService) Info(ctx context.Context,req *request.Object) (*content.Diary, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	var diary content.Diary
	err = db.Where(`id = ? AND user_id = ?`,req.Id,auth.Id).First(&diary).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "First")
	}
	return nil, nil
}
func (*DiaryService) Add(ctx context.Context,req *content.AddDiaryReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	req.UserId = auth.Id
	err = db.Table(model.DiaryTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return nil, nil
}
func (*DiaryService) Edit(context.Context, *content.AddDiaryReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (*DiaryService) List(context.Context, *content.DiaryListReq) (*content.DiaryListRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (*DiaryService) Delete(ctx context.Context,req *request.Object) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi,true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	err = contentDao.DelByAuthDB(db,model.DiaryTableName,req.Id,auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
