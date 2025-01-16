package service

import (
	"context"
	"github.com/hopeio/context/httpctx"
	gormi "github.com/hopeio/utils/dao/database/gorm"
	"github.com/liov/hoper/server/go/content/data"
	"github.com/liov/hoper/server/go/content/global"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hopeio/protobuf/errcode"
	"github.com/hopeio/protobuf/request"
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
	ctxi := httpctx.FromContextValue(ctx)
	_, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	return nil, status.Errorf(codes.Unimplemented, "method DiaryBookList not implemented")
}
func (*DiaryService) AddDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*request.Id, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiaryBook).Create(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
	}
	return &request.Id{Id: req.Id}, nil
}
func (*DiaryService) EditDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*emptypb.Empty, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiaryBook).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
	}
	return nil, nil
}
func (*DiaryService) Info(ctx context.Context, req *request.Id) (*content.Diary, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.Limit(&global.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	var diary content.Diary
	err = db.Where(`id = ? AND user_id = ?`, req.Id, auth.Id).First(&diary).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "First")
	}
	return nil, nil
}
func (*DiaryService) Add(ctx context.Context, req *content.AddDiaryReq) (*request.Id, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiary).Create(req).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
	}
	return nil, nil
}

func (*DiaryService) Edit(context.Context, *content.AddDiaryReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*DiaryService) List(context.Context, *content.DiaryListReq) (*content.DiaryListRep, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (*DiaryService) Delete(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)
	err = contentDBDao.DelByAuth(model.TableNameDiary, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
