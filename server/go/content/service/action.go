package service

import (
	"context"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/scaffold/errcode"

	gormx "github.com/hopeio/gox/database/sql/gorm"
	"github.com/liov/hoper/server/go/content/data"
	dbdao "github.com/liov/hoper/server/go/content/data/db"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/content"
	"github.com/liov/hoper/server/go/protobuf/user"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hopeio/gox/container/set"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/slices"
	"github.com/hopeio/protobuf/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ActionService struct {
	content.UnimplementedActionServiceServer
}

func (*ActionService) Like(ctx context.Context, req *content.LikeReq) (*request.Id, error) {
	if req.Action != content.ActionLike && req.Action != content.ActionUnlike && req.Action != content.ActionBrowse {
		return nil, nil
	}
	ctxi, _ := httpctx.FromContext(ctx)
	span := ctxi.StartSpan("Like")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())

	contentDBDao := dbdao.GetDao(ctxi, db)

	req.UserId = auth.Id

	id, err := contentDBDao.LikeId(req.Type, req.Action, req.RefId, req.UserId)
	if err != nil {
		return &request.Id{Id: id}, err
	}
	if id > 0 {
		return &request.Id{Id: id}, nil
	}

	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		contenttxDBDao := dbdao.GetDao(ctxi, tx)
		err = tx.Table(model.TableNameLike).Create(req).Error
		if err != nil {
			return ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}

		err = contenttxDBDao.ActionCount(req.Type, req.Action, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.HotCount(req.Type, req.RefId, 1)

	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.RedisErr, err, "HotCountRedis")
	}
	return &request.Id{Id: req.Id}, nil
}

func (*ActionService) DelLike(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	span := ctxi.StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	like, err := contentDBDao.GetLike(req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	if like.Id == 0 {
		return nil, errcode.InvalidArgument
	}
	err = contentDBDao.DelByAuth(model.TableNameLike, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	err = contentDBDao.ActionCount(like.Type, like.Action, like.RefId, -1)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.HotCount(like.Type, like.RefId, -1)
	if err != nil {
		return nil, err
	}
	return new(emptypb.Empty), err
}

func (*ActionService) Comment(ctx context.Context, req *content.CommentReq) (*request.Id, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())

	req.UserId = auth.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		contenttxDBDao := dbdao.GetDao(ctxi, tx)
		err = tx.Table(model.TableNameComment).Create(req).Error
		if err != nil {
			return ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}
		err = contenttxDBDao.CreateContextExt(content.ContentComment, req.Id)
		if err != nil {
			return err
		}
		err = contenttxDBDao.ActionCount(req.Type, content.ActionComment, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err != errcode.DBError {
			ctxi.ErrorLog(err, zap.String(log.FieldPosition, "Transaction"))
		}
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.HotCount(req.Type, req.RefId, 1)
	if err != nil {
		return nil, err
	}
	return &request.Id{Id: req.Id}, nil
}

func (*ActionService) DelComment(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	var comment content.Comment
	err = db.Table(model.TableNameComment).First(&comment, "id = ?", req.Id).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "Find")
	}
	if comment.UserId != auth.Id {
		var userId uint64
		err = db.Table(model.ContentTableName(comment.Type)).Select("user_id").
			Where(`id = ?`, comment.RefId).Scan(&userId).Error
		if err != nil {
			return nil, ctxi.RespErrorLog(errcode.DBError, err, "SelectUserId")
		}
		if userId != auth.Id {
			return nil, errcode.PermissionDenied
		}
	}

	err = contentDBDao.Del(model.TableNameComment, req.Id)
	if err != nil {
		return nil, err
	}
	err = contentDBDao.ActionCount(comment.Type, content.ActionComment, comment.RefId, -1)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.HotCount(comment.Type, comment.RefId, -1)
	if err != nil {
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (*ActionService) Collect(ctx context.Context, req *content.CollectReq) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	req.UserId = auth.Id
	collects, err := contentDBDao.GetCollects(req.Type, []uint64{req.RefId}, auth.Id)
	if err != nil {
		return nil, err
	}
	var origin []uint64
	for _, collect := range collects {
		origin = append(origin, collect.FavId)
	}
	diff := slices.DifferenceSet(origin, req.FavIds)
	collect := model.Collect{
		Type:   req.Type,
		RefId:  req.RefId,
		UserId: auth.Id,
		FavId:  0,
	}
	for _, id := range diff {
		collect.FavId = id
		err = db.Table(model.TableNameCollect).Create(&collect).Error
		if err != nil {
			return nil, ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}
	}
	if len(origin) == 0 && len(req.FavIds) > 0 {
		err = contentDBDao.ActionCount(req.Type, content.ActionCollect, req.RefId, 1)
		if err != nil {
			return nil, ctxi.RespErrorLog(errcode.DBError, err, "ActionCount")
		}
	}
	err = db.Table(model.TableNameCollect).Where(`type = ? AND ref_id = ? AND fav_id NOT IN (?)`, req.Type, req.RefId, req.FavIds).
		Update(`deleted_at`, ctxi.RequestTime.String()).Error
	if err != nil {
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "DELETE")
	}
	var hotCount float64
	if len(origin) == 0 && len(req.FavIds) > 0 {
		hotCount = 1
	}
	if len(origin) > 0 && len(req.FavIds) == 0 {
		hotCount = -1
	}
	if hotCount != 0 {
		contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
		err = contentRedisDao.HotCount(req.Type, req.RefId, hotCount)
		if err != nil {
			return nil, err
		}
	}

	return new(emptypb.Empty), nil
}

func (*ActionService) Report(ctx context.Context, req *content.ReportReq) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentRedisDao := data.GetRedisDao(ctxi, global.Dao.Redis)
	err = contentRedisDao.Limit(&global.Conf.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	req.UserId = auth.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		contenttxDBDao := data.GetDBDao(ctxi, tx)
		err = tx.Table(model.TableNameReport).Create(req).Error
		if err != nil {
			return ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}
		err = contenttxDBDao.ActionCount(req.Type, content.ActionReport, req.RefId, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err != errcode.DBError {
			ctxi.ErrorLog(err, zap.String(log.FieldPosition, "Transaction"))
		}
		return nil, errcode.DBError
	}
	return new(emptypb.Empty), nil
}

func (*ActionService) CommentList(ctx context.Context, req *content.CommentListReq) (*content.CommentListResp, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	total, comments, err := contentDBDao.GetComments(content.ContentMoment, req.RefId, req.RootId, req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}
	var ids []uint64
	userIds := set.New[uint64]()
	var m = make(map[uint64]*content.Comment)
	for i := range comments {
		ids = append(ids, comments[i].Id)
		m[comments[i].Id] = comments[i]
		userIds.Add(comments[i].UserId)
		userIds.Add(comments[i].RecvId)
		// 屏蔽字段
		commentMaskField(comments[i])
	}
	// ext
	statistics, err := contentDBDao.GetStatistics(content.ContentComment, ids)
	if err != nil {
		return nil, err
	}
	for i := range statistics {
		if comment, ok := m[statistics[i].Id]; ok {
			comment.Statistics = statistics[i]
		}
	}

	//like
	if auth.Id != 0 {
		likes, err := contentDBDao.GetContentActions(content.ActionLike, content.ContentComment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range likes {
			if comment, ok := m[likes[i].RefId]; ok {
				if comment.Action == nil {
					comment.Action = &content.UserAction{}
				}
				if likes[i].Action == content.ActionLike {
					comment.Action.LikeId = likes[i].Id
				}
				if likes[i].Action == content.ActionUnlike {
					comment.Action.UnlikeId = likes[i].Id
				}
			}
		}
		collects, err := contentDBDao.GetCollects(content.ContentComment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range collects {
			if comment, ok := m[collects[i].RefId]; ok {
				if comment.Action == nil {
					comment.Action = &content.UserAction{}
				}
				comment.Action.CollectIds = append(comment.Action.CollectIds, collects[i].RefId)
			}
		}
	}
	var users []*user.UserBase
	if len(userIds) > 0 {
		userList, err := data.UserClient().BaseList(ctxi.Base(), &user.BaseListReq{Ids: userIds.ToSlice()})
		if err != nil {
			return nil, err
		}
		users = userList.List
	}
	return &content.CommentListResp{
		Total: total,
		List:  comments,
		Users: users,
	}, nil
}

// 屏蔽字段
func commentMaskField(comment *content.Comment) {
	comment.ModelTime.DeletedAt = nil
}

func (*ActionService) GetUserAction(ctx context.Context, req *content.ContentReq) (*content.UserAction, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}

	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := dbdao.GetDao(ctxi, db)

	action := &content.UserAction{}
	likes, err := contentDBDao.GetContentActions(content.ActionLike, content.ContentMoment, []uint64{req.RefId}, auth.Id)
	if err != nil {
		return nil, err
	}

	for i := range likes {
		if likes[i].Action == content.ActionLike {
			action.LikeId = likes[i].Id
		}
		if likes[i].Action == content.ActionUnlike {
			action.UnlikeId = likes[i].Id
		}
	}
	collects, err := contentDBDao.GetCollects(content.ContentMoment, []uint64{req.RefId}, auth.Id)
	if err != nil {
		return nil, err
	}
	for i := range collects {
		action.CollectIds = append(action.CollectIds, collects[i].RefId)
	}
	return action, nil
}
