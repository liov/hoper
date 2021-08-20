package service

import (
	"context"
	"github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/protobuf/utils/empty"
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
	contexti "github.com/liov/hoper/v2/tiga/context"
	"github.com/liov/hoper/v2/user/dao"
	"github.com/liov/hoper/v2/user/model"
	dbi "github.com/liov/hoper/v2/utils/dao/db"
)

// 关注
func (u *UserService) Follow(ctx context.Context, req *user.FollowReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB)
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
	return new(empty.Empty), nil
}

// 取消关注
func (u *UserService) DelFollow(ctx context.Context, req *user.FollowReq) (*user.BaseListRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(dao.Dao.GORMDB)
	userDao := dao.GetDao(ctxi)
	exists, err := userDao.FollowExistsDB(db, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	err = db.Table(model.FollowTableName).Where("user_id = ? AND follow_id = ? AND "+dbi.PostgreNotDeleted, req.Id, auth.Id).
		UpdateColumn("deleted_at", ctxi.RequestAt.TimeString).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
	}
	return nil, nil
}
