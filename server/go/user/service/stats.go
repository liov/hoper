package service

import (
	"context"

	"github.com/hopeio/gox/context/httpctx"
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
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()

	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	userDao := data.GetDBDao(ctxi, global.Dao.GORMDB.DB)
	exists, err := userDao.FollowExistsDB(req.Id, auth.Id)
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
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
	}
	return new(emptypb.Empty), nil
}

// 取消关注
func (u *UserService) DelFollow(ctx context.Context, req *user.FollowReq) (*user.BaseListResp, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()

	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	userDao := data.GetDBDao(ctxi, global.Dao.GORMDB.DB)
	exists, err := userDao.FollowExistsDB(req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	err = userDao.Table(model.TableNameFollow).Where("user_id = ? AND follow_id = ?"+sqlx.WithNotDeleted, req.Id, auth.Id).
		UpdateColumn("deleted_at", ctxi.RequestAt.String()).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
	}
	return nil, nil
}
