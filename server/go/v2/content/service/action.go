package service

import (
	"context"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"gorm.io/gorm"
)

type ActionService struct {
	content.UnimplementedActionServiceServer
}

func (*ActionService) Like(ctx context.Context, req *content.LikeReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
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
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(model.LikeTableName).Create(req).Error
		if err != nil {
			return ctxi.Log(err, "Create", err.Error())
		}
		err = contentDao.ActionCountDB(tx, req.Type, req.Action, req.RefId, 1)
		if err != nil {
			return ctxi.Log(err, "ActionCountDB", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	err = contentDao.LikeCountRedis(dao.Dao.Redis, req.Type, req.RefId, 1)
	if err != nil {
		return nil, ctxi.Log(errorcode.RedisErr, "LikeCountRedis", err.Error())
	}
	return nil, nil
}
func (*ActionService) Comment(ctx context.Context, req *content.CommentReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
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
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(model.CommentTableName).Create(req).Error
		if err != nil {
			return ctxi.Log(err, "Create", err.Error())
		}
		err = contentDao.ActionCountDB(tx, req.Type, content.ActionComment, req.RefId, 1)
		if err != nil {
			return ctxi.Log(err, "ActionCountDB", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
func (*ActionService) Collect(ctx context.Context, req *content.CollectReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
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
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(model.CollectTableName).Create(req).Error
		if err != nil {
			return ctxi.Log(err, "Create", err.Error())
		}
		err = contentDao.ActionCountDB(tx, req.Type, content.ActionCollect, req.RefId, 1)
		if err != nil {
			return ctxi.Log(err, "ActionCountDB", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
func (*ActionService) Report(ctx context.Context, req *content.ReportReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
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
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(model.ReportTableName).Create(req).Error
		if err != nil {
			return ctxi.Log(err, "Create", err.Error())
		}
		err = contentDao.ActionCountDB(tx, req.Type, content.ActionReport, req.RefId, 1)
		if err != nil {
			return ctxi.Log(err, "ActionCountDB", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
