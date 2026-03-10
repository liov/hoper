package service

import (
	"context"

	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/data"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/content"
	"google.golang.org/protobuf/types/known/emptypb"

	"net/http"

	"github.com/hopeio/protobuf/request"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiaryService struct {
	content.UnimplementedDiaryServiceServer
}

func (m *DiaryService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "日记相关", "/api/diary", nil
}

func (*DiaryService) DiaryBook(ctx context.Context, req *content.DiaryBookReq) (*content.DiaryBookResp, error) {

	return nil, status.Errorf(codes.Unimplemented, "method DiaryBook not implemented")
}
func (*DiaryService) DiaryBookList(ctx context.Context, req *content.DiaryBookListReq) (*content.DiaryBookListResp, error) {

	_, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}

	return nil, status.Errorf(codes.Unimplemented, "method DiaryBookList not implemented")
}
func (*DiaryService) AddDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*request.Id, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiaryBook).Create(req).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return &request.Id{Id: req.Id}, nil
}
func (*DiaryService) EditDiaryBook(ctx context.Context, req *content.AddDiaryBookReq) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiaryBook).Where(`id =? AND user_id =?`, req.Id, auth.Id).
		Updates(req).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}
func (*DiaryService) Info(ctx context.Context, req *request.Id) (*content.Diary, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	rclient := global.Dao.Redis.WithContext(ctx)
	contentRedisDao := data.GetRedisDao(rclient)
	err = contentRedisDao.Limit(ctx, &global.Conf.Moment.Limit, auth.Id)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	var diary content.Diary
	err = db.Where(`id = ? AND user_id = ?`, req.Id, auth.Id).First(&diary).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}
func (*DiaryService) Add(ctx context.Context, req *content.AddDiaryReq) (*request.Id, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	req.UserId = auth.Id
	err = db.Table(model.TableNameDiary).Create(req).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}

func (*DiaryService) Edit(context.Context, *content.AddDiaryReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*DiaryService) List(context.Context, *content.DiaryListReq) (*content.DiaryListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (*DiaryService) Delete(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	contentDBDao := data.GetDBDao(db)
	err = contentDBDao.DelByAuth(model.TableNameDiary, req.Id, auth.Id)
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}
