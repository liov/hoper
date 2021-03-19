package service

import (
	"context"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/utils/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActionService struct {
	content.UnimplementedActionServiceServer
}

func (*ActionService) Like(ctx context.Context, req *content.LikeReq) (*empty.Empty, error) {
	if req.Action != content.ActionLike && req.Action != content.ActionUnlike && req.Action != content.ActionBrowse {
		return nil, nil
	}
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

	exists, err := contentDao.ActionExists(db, req.Type, req.Action, req.RefId, req.UserId)
	if err != nil {
		return nil, err
	}
	if (exists && !req.Del) || (!exists && req.Del) {
		return nil, nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if req.Del {
			err = contentDao.DelActionDB(tx, req.Type, req.Action, req.RefId, auth.Id)
			if err != nil {
				return err
			}
		} else {
			err = db.Table(model.LikeTableName).Create(req).Error
			if err != nil {
				return err
			}
		}
		err = contentDao.ActionCountDB(tx, req.Type, req.Action, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	if req.Del {
		err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, -1)
	} else {
		err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, 1)
	}

	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err,  "HotCountRedis")
	}
	return nil, nil
}

func (*ActionService) DelLike(ctx context.Context, req *request.Object) (*empty.Empty, error) {
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
	err = contentDao.DelByAuthDB(db, model.LikeTableName, req.Id, auth.Id)
	return nil, err
}

func (*ActionService) Comment(ctx context.Context, req *content.CommentReq) (*empty.Empty, error) {
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
			return ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
		err = contentDao.ActionCountDB(tx, req.Type, content.ActionComment, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err != errorcode.DBError {
			ctxi.Error(err.Error(), zap.String(log.Position, "Transaction"))
		}
		return nil, errorcode.DBError
	}
	err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, 1)
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr,  err, "HotCountRedis")
	}
	return nil, nil
}

func (*ActionService) DelComment(ctx context.Context, req *request.Object) (*empty.Empty, error) {
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
	var comment content.Comment
	err = db.Table(model.CommentTableName).First(&comment, "id = ?", req.Id).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError,  err, "Find")
	}
	if comment.UserId != auth.Id {
		var userId uint64
		err = db.Table(model.ContentTableName(comment.Type)).Select("user_id").
			Where(`id = ?`, comment.RefId).Row().Scan(&userId)
		if err != nil {
			return nil, ctxi.ErrorLog(errorcode.DBError,  err, "SelectUserId")
		}
		if userId != auth.Id {
			return nil, errorcode.PermissionDenied
		}
	}

	err = contentDao.DelDB(db, model.CommentTableName, req.Id)
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError,  err, "DelDB")
	}
	err = contentDao.HotCountRedis(dao.Dao.Redis, comment.Type, comment.RefId, -1)
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.RedisErr,  err, "HotCountRedis")
	}
	return nil, nil
}

func (*ActionService) Collect(ctx context.Context, req *content.CollectReq) (*empty.Empty, error) {
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
	exists, err := contentDao.ExistsByAuthDB(db, model.FavoritesTableName, req.FavId, auth.Id)
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err,  "ExistsByAuthDB")
	}
	if !exists {
		return nil, errorcode.PermissionDenied.Message("无效的收藏夹")
	}
	err = db.Table(model.CollectTableName).Create(req).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError,  err, "Create")
	}
	err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, -1)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (*ActionService) DelCollect(ctx context.Context, req *request.Object) (*empty.Empty, error) {
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
	err = contentDao.DelByAuthDB(db, model.CollectTableName, req.Id, auth.Id)
	return nil, err
}

func (*ActionService) Report(ctx context.Context, req *content.ReportReq) (*empty.Empty, error) {
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
			return ctxi.ErrorLog(errorcode.DBError,err, "Create")
		}
		err = contentDao.ActionCountDB(tx, req.Type, content.ActionReport, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err != errorcode.DBError {
			ctxi.Error(err.Error(), zap.String(log.Position, "Transaction"))
		}
		return nil, errorcode.DBError
	}
	return nil, nil
}
