package service

import (
	"context"
	"github.com/liov/hoper/v2/content/client"
	"github.com/liov/hoper/v2/content/conf"
	"github.com/liov/hoper/v2/content/dao"
	"github.com/liov/hoper/v2/content/model"
	"github.com/liov/hoper/v2/protobuf/content"
	"github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/protobuf/utils/empty"
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/v2/protobuf/utils/request"
	contexti "github.com/liov/hoper/v2/tiga/context"
	"github.com/liov/hoper/v2/utils/log"
	"github.com/liov/hoper/v2/utils/slices"
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

	db := ctxi.NewDB(dao.Dao.GORMDB)
	req.UserId = auth.Id

	id, err := contentDao.LikeIdDB(db, req.Type, req.Action, req.RefId, req.UserId)
	if err != nil {
		return &request.Object{Id: id}, err
	}
	if id > 0 {
		return &request.Object{Id: id}, nil
	}

	err = contentDao.Transaction(db, func(tx *gorm.DB) error {

		err = db.Table(model.LikeTableName).Create(req).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}

		err = contentDao.ActionCountDB(tx, req.Type, req.Action, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = contentDao.HotCountRedis(dao.Dao.Redis, req.Type, req.RefId, 1)

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
	db := ctxi.NewDB(dao.Dao.GORMDB)
	like, err := contentDao.GetLikeDB(db, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if like.Id == 0 {
		return nil, errorcode.ParamInvalid
	}
	err = contentDao.DelByAuthDB(db, model.LikeTableName, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	err = contentDao.ActionCountDB(db, like.Type, like.Action, like.RefId, -1)
	if err != nil {
		return nil, err
	}
	err = contentDao.HotCountRedis(dao.Dao.Redis, like.Type, like.RefId, -1)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
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

	db := ctxi.NewDB(dao.Dao.GORMDB)
	req.UserId = auth.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.Table(model.CommentTableName).Create(req).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
		err = contentDao.CreateContextExt(tx, content.ContentComment, req.Id)
		if err != nil {
			return err
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
		return nil, err
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

	db := ctxi.NewDB(dao.Dao.GORMDB)
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
	err = contentDao.ActionCountDB(db, comment.Type, content.ActionComment, comment.RefId, -1)
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

	db := ctxi.NewDB(dao.Dao.GORMDB)
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
	db := ctxi.NewDB(dao.Dao.GORMDB)
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

func (*ActionService) CommentList(ctx context.Context, req *content.CommentListReq) (*content.CommentListRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	db := ctxi.NewDB(dao.Dao.GORMDB)
	total, comments, err := contentDao.GetCommentsDB(db, content.ContentMoment, req.RefId, req.RootId, int(req.PageNo), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	var ids, userIds []uint64
	var m = make(map[uint64]*content.Comment)
	for i := range comments {
		ids = append(ids, comments[i].Id)
		m[comments[i].Id] = comments[i]
		userIds = append(userIds, comments[i].UserId)
		userIds = append(userIds, comments[i].RecvId)
		// 屏蔽字段
		commentMaskField(comments[i])
	}
	// ext
	exts, err := contentDao.GetContentExtDB(db, content.ContentMoment, ids)
	if err != nil {
		return nil, err
	}
	for i := range exts {
		if moment, ok := m[exts[i].RefId]; ok {
			moment.Ext = exts[i]
		}
	}

	//like
	if auth.Id != 0 {
		likes, err := contentDao.GetContentActionsDB(db, content.ActionLike, content.ContentComment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range likes {
			if moment, ok := m[likes[i].RefId]; ok {
				if likes[i].Action == content.ActionLike {
					moment.LikeId = likes[i].Id
				}
				if likes[i].Action == content.ActionUnlike {
					moment.UnlikeId = likes[i].Id
				}
			}
		}
		collects, err := contentDao.GetCollectsDB(db, content.ContentComment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range collects {
			if moment, ok := m[collects[i].RefId]; ok {
				moment.Collects = append(moment.Collects, collects[i].FavId)
			}
		}
	}
	var users []*user.UserBaseInfo
	if len(userIds) > 0 {
		userList, err := client.UserClient.BaseList(ctxi, &user.BaseListReq{Ids: userIds})
		if err != nil {
			return nil, err
		}
		users = userList.List
	}
	return &content.CommentListRep{
		Total: total,
		List:  comments,
		Users: users,
	}, nil
}

// 屏蔽字段
func commentMaskField(comment *content.Comment) {
	comment.DeletedAt = ""
	comment.CreatedAt = comment.CreatedAt[:19]
}
