package service

import (
	"context"
	contexti "github.com/actliboy/hoper/server/go/lib/context"
	"github.com/actliboy/hoper/server/go/lib/protobuf/empty"
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	"github.com/actliboy/hoper/server/go/lib/protobuf/request"
	"github.com/actliboy/hoper/server/go/mod/content/conf"
	"github.com/actliboy/hoper/server/go/mod/content/dao"
	"github.com/actliboy/hoper/server/go/mod/content/model"
	"github.com/actliboy/hoper/server/go/mod/protobuf/content"
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
func (*DiaryService) DiaryBookList(ctx context.Context, req *content.DiaryBookListReq) (*content.DiaryBookListRep, error) {
	ctxi := contexti.CtxFromContext(ctx)
	_, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	return nil, status.Errorf(codes.Unimplemented, "method DiaryBookList not implemented")
}
func (*DiaryService) AddDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*request.Object, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	req.UserId = auth.Id
	err = db.Table(model.DiaryBookTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return &request.Object{Id: req.Id}, nil
}
func (*DiaryService) EditDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	req.UserId = auth.Id
	err = db.Table(model.DiaryBookTableName).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return nil, nil
}
func (*DiaryService) Info(ctx context.Context, req *request.Object) (*content.Diary, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := dao.GetRedisDao(ctxi, dao.Dao.Redis)
	err = contentRedisDao.Limit(&conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	var diary content.Diary
	err = db.Where(`id = ? AND user_id = ?`, req.Id, auth.Id).First(&diary).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "First")
	}
	return nil, nil
}
func (*DiaryService) Add(ctx context.Context, req *content.AddDiaryReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
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
func (*DiaryService) Delete(ctx context.Context, req *request.Object) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB.DB)
	contentDBDao := dao.GetDBDao(ctxi, db)
	err = contentDBDao.DelByAuth(model.DiaryTableName, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
