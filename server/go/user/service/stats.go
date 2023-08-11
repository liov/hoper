package service

import (
	"context"
	"github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/actliboy/hoper/server/go/user/confdao"
	"github.com/actliboy/hoper/server/go/user/dao"
	"github.com/actliboy/hoper/server/go/user/model"
	"github.com/hopeio/zeta/context/http_context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hopeio/zeta/protobuf/errorcode"
	dbi "github.com/hopeio/zeta/utils/dao/db/const"
)

// 关注
func (u *UserService) Follow(ctx context.Context, req *user.FollowReq) (*emptypb.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	exists, err := userDao.FollowExistsDB(db, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}
	err = db.Table(model.FollowTableName).Create(&user.UserFollow{
		UserId:   req.Id,
		FollowId: auth.Id,
	}).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return new(emptypb.Empty), nil
}

// 取消关注
func (u *UserService) DelFollow(ctx context.Context, req *user.FollowReq) (*user.BaseListRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	userDao := dao.GetDao(ctxi)
	exists, err := userDao.FollowExistsDB(db, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	err = db.Table(model.FollowTableName).Where("user_id = ? AND follow_id = ?"+dbi.WithNotDeleted, req.Id, auth.Id).
		UpdateColumn("deleted_at", ctxi.RequestAt.TimeString).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return nil, nil
}
