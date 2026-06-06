package service

import (
	"context"
	"time"

	sqlx "github.com/hopeio/gox/database/sql"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/data"
	"github.com/liov/hoper/server/go/user/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 关注
func (u *UserService) Follow(ctx context.Context, req *user.FollowReq) (*emptypb.Empty, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDao := data.GetDBDao(db)
	exists, err := userDao.FollowExistsDB(ctx, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}
	err = userDao.Table(model.TableNameFollow).Create(&user.Follow{
		UserId:   req.Id,
		FollowId: auth.Id,
	}).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return new(emptypb.Empty), nil
}

// 取消关注
func (u *UserService) DelFollow(ctx context.Context, req *user.FollowReq) (*user.BaseListResp, error) {

	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	userDao := data.GetDBDao(db)
	exists, err := userDao.FollowExistsDB(ctx, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	err = userDao.Table(model.TableNameFollow).Where("user_id = ? AND follow_id = ?"+sqlx.WithNotDeleted, req.Id, auth.Id).
		UpdateColumn("deleted_at", time.Now()).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return nil, nil
}
