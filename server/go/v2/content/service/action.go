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
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/slices"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActionService struct {
	content.UnimplementedActionServiceServer
}

func (*ActionService) Like(ctx context.Context, req *content.LikeReq) (*request.Object, error) {
	if req.Action != content.ActionLike && req.Action != content.ActionUnlike && req.Action != content.ActionBrowse {
		return nil, nil
	}
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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

	id, err := contentDao.LikeIdDB(db, req.Type, req.Action, req.RefId, req.UserId)
	if err != nil {
		return &request.Object{Id: id}, err
	}
	if (id > 0 && !req.Del) || (id == 0 && req.Del) {
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
		return nil, ctxi.ErrorLog(errorcode.RedisErr, err, "HotCountRedis")
	}
	return &request.Object{Id: req.Id}, nil
}

func (*ActionService) DelLike(ctx context.Context, req *request.Object) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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

func (*ActionService) Comment(ctx context.Context, req *content.CommentReq) (*request.Object, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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
		return nil, err
	}
	return &request.Object{Id: req.Id}, nil
}

func (*ActionService) DelComment(ctx context.Context, req *request.Object) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "Find")
	}
	if comment.UserId != auth.Id {
		var userId uint64
		err = db.Table(model.ContentTableName(comment.Type)).Select("user_id").
			Where(`id = ?`, comment.RefId).Row().Scan(&userId)
		if err != nil {
			return nil, ctxi.ErrorLog(errorcode.DBError, err, "SelectUserId")
		}
		if userId != auth.Id {
			return nil, errorcode.PermissionDenied
		}
	}

	err = contentDao.DelDB(db, model.CommentTableName, req.Id)
	if err != nil {
		return nil, err
	}
	err = contentDao.HotCountRedis(dao.Dao.Redis, comment.Type, comment.RefId, -1)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (*ActionService) Collect(ctx context.Context, req *content.CollectReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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
	collects, err := contentDao.GetCollectsDB(db, req.Type, []uint64{req.RefId}, auth.Id)
	if err != nil {
		return nil, err
	}
	var origin []uint64
	for _, collect := range collects {
		origin = append(origin, collect.FavId)
	}
	diff := slices.DiffUint64(origin, req.FavIds)
	collect := model.Collect{
		Type:   req.Type,
		RefId:  req.RefId,
		UserId: auth.Id,
		FavId:  0,
	}
	for _, id := range diff {
		collect.FavId = id
		err = db.Table(model.CollectTableName).Create(&collect).Error
		if err != nil {
			return nil, ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
	}
	if len(origin) == 0 && len(req.FavIds) > 0 {
		err = contentDao.ActionCountDB(db, req.Type, content.ActionCollect, req.RefId, 1)
		if err != nil {
			return nil, ctxi.ErrorLog(errorcode.DBError, err, "ActionCountDB")
		}
	}
	err = db.Table(model.CollectTableName).Where(`type = ? AND ref_id = ? AND fav_id NOT IN (?)`, req.Type, req.RefId, req.FavIds).
		Update(`deleted_at`, ctxi.RequestAt.TimeString).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "DELETE")
	}
	var hotCount float64
	if len(origin) == 0 && len(req.FavIds) > 0 {
		hotCount = 1
	}
	if len(origin) > 0 && len(req.FavIds) == 0 {
		hotCount = -1
	}
	if hotCount != 0 {
		err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, hotCount)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (*ActionService) Report(ctx context.Context, req *content.ReportReq) (*empty.Empty, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
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
			return ctxi.ErrorLog(errorcode.DBError, err, "Create")
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
