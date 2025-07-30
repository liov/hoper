package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/scaffold/errcode"

	"github.com/hopeio/protobuf/request"
	"github.com/hopeio/gox/datastructure/set"
	gormi "github.com/hopeio/gox/datax/database/gorm"
	"github.com/liov/hoper/server/go/protobuf/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"unicode/utf8"

	comdata "github.com/liov/hoper/server/go/common/data"
	comconfdao "github.com/liov/hoper/server/go/common/global"
	"github.com/liov/hoper/server/go/content/data"
	"github.com/liov/hoper/server/go/content/global"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"github.com/liov/hoper/server/go/protobuf/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MomentService struct {
	content.UnimplementedMomentServiceServer
}

func (*MomentService) Service() (describe, prefix string, middleware []gin.HandlerFunc) {
	return "瞬间相关", "/api/v1/moment", nil
}

func (*MomentService) Info(ctx context.Context, req *request.Id) (*content.Moment, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, _ := auth(ctxi, true)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	var moment content.Moment
	err := db.Table(model.TableNameMoment).First(&moment, req.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.NotFound
		}
		return nil, ctxi.RespErrorLog(errcode.DBError, err, "First")
	}
	// tags
	contentTags, err := contentDBDao.GetContentTag(content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, err
	}
	var tags = make([]*common.TinyTag, len(contentTags))
	for i := range contentTags {
		tags[i] = &contentTags[i].TinyTag
	}
	moment.Tags = tags

	//like
	if auth != nil {
		action := &content.UserAction{}
		likes, err := contentDBDao.GetContentActions(content.ActionLike, content.ContentMoment, []uint64{req.Id}, auth.Id)
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
		collects, err := contentDBDao.GetCollects(content.ContentMoment, []uint64{req.Id}, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range collects {
			action.CollectIds = append(action.CollectIds, collects[i].RefId)
		}
		if len(likes) != 0 && len(collects) != 0 {
			moment.Action = action
		}
	}
	// ext
	statistics, err := contentDBDao.GetStatistics(content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, err
	}
	moment.Statistics = statistics[0]

	// 匿名
	if moment.Anonymous == 1 {
		moment.UserId = 0
	} else {
		moment.User = &user.UserBase{
			Id:     auth.Id,
			Name:   auth.Name,
			Score:  0,
			Gender: 0,
			Avatar: auth.Avatar,
		}
	}
	momentMaskField(&moment)
	return &moment, nil
}

// 屏蔽字段
func momentMaskField(moment *content.Moment) {
	moment.DeletedAt = nil
	moment.Anonymous = 0
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Id, error) {

	if utf8.RuneCountInString(req.Content) < global.Conf.Customize.Moment.MaxContentLen {
		return nil, errcode.InvalidArgument.Msg(fmt.Sprintf("文章内容不能小于%d个字", global.Conf.Customize.Moment.MaxContentLen))
	}

	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	codb := gormi.NewTraceDB(comconfdao.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)
	commonDBDao := comdata.GetDBDao(ctxi, codb)

	req.UserId = auth.Id

	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errcode.ParamInvalid.Msg("心情不存在")
		}*/
	var tags []model.TinyTag
	if len(req.Tags) > 0 {
		tags, err = commonDBDao.GetTagsByName(req.Tags)
		if err != nil {
			return nil, err
		}
	}

	req.UserId = auth.Id
	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		if req.Permission == 0 {
			req.Permission = content.ViewPermissionAll
		}
		contenttxDBDao := data.GetDBDao(ctxi, tx)
		err = tx.Table(model.TableNameMoment).Create(req).Error
		if err != nil {
			return ctxi.RespErrorLog(errcode.DBError, err, "tx.CreateReq")
		}
		err = contenttxDBDao.CreateContextExt(content.ContentMoment, req.Id)
		if err != nil {
			return err
		}
		var contentTags []model.ContentTag
		var noExist []common.Tag
	Loop:
		for i := range req.Tags {
			// 性能可以优化
			for j := range tags {
				if req.Tags[i] == tags[j].Name {
					contentTags = append(contentTags, model.ContentTag{
						Type:  content.ContentMoment,
						RefId: req.Id,
						TagId: tags[j].Id,
					})
					continue Loop
				}
			}
			noExist = append(noExist, common.Tag{Name: req.Tags[i], UserId: auth.Id})
		}
		if len(noExist) == 1 {
			if err = tx.Create(&noExist[1]).Error; err != nil {
				return ctxi.RespErrorLog(errcode.DBError, err, "db.CreateNoExist")
			}
		}
		if len(noExist) > 1 {
			if err = tx.Create(&noExist).Error; err != nil {
				return ctxi.RespErrorLog(errcode.DBError, err, "db.CreateNoExist")
			}
		}
		for i := range noExist {
			contentTags = append(contentTags, model.ContentTag{
				Type:  content.ContentMoment,
				RefId: req.Id,
				TagId: noExist[i].Id,
			})
		}
		if len(contentTags) > 0 {
			if err = tx.Create(&contentTags).Error; err != nil {
				return ctxi.RespErrorLog(errcode.DBError, err, "db.CreateContentTags")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &request.Id{Id: req.Id}, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(ctx context.Context, req *content.MomentListReq) (*content.MomentListRep, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, _ := auth(ctxi, true)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	total, moments, err := contentDBDao.GetMomentList(req)
	if err != nil {
		return nil, err
	}

	if len(moments) == 0 {
		return &content.MomentListRep{
			Total: total,
			List:  nil,
			Users: nil,
		}, nil
	}

	var m = make(map[uint64]*content.Moment)
	var ids []uint64
	userIds := set.New[uint64]()
	for i := range moments {
		ids = append(ids, moments[i].Id)
		m[moments[i].Id] = moments[i]
		userIds.Add(moments[i].UserId)
		// 屏蔽字段
		momentMaskField(moments[i])
	}

	// tag
	tags, err := contentDBDao.GetContentTag(content.ContentMoment, ids)
	if err != nil {
		return nil, err
	}

	for i := range tags {
		if moment, ok := m[tags[i].RefId]; ok {
			moment.Tags = append(moment.Tags, &tags[i].TinyTag)
		}
	}
	// ext
	statistics, err := contentDBDao.GetStatistics(content.ContentMoment, ids)
	if err != nil {
		return nil, err
	}
	for i := range statistics {
		if moment, ok := m[statistics[i].Id]; ok {
			moment.Statistics = statistics[i]
		}
	}
	//like
	if auth != nil {
		likes, err := contentDBDao.GetContentActions(content.ActionLike, content.ContentMoment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range likes {
			if moment, ok := m[likes[i].RefId]; ok {
				if moment.Action == nil {
					moment.Action = &content.UserAction{}
				}
				if likes[i].Action == content.ActionLike {
					moment.Action.LikeId = likes[i].Id
				}
				if likes[i].Action == content.ActionUnlike {
					moment.Action.UnlikeId = likes[i].Id
				}
			}
		}
		collects, err := contentDBDao.GetCollects(content.ContentMoment, ids, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range collects {
			if moment, ok := m[collects[i].RefId]; ok {
				if moment.Action == nil {
					moment.Action = &content.UserAction{}
				}
				moment.Action.CollectIds = append(moment.Action.CollectIds, collects[i].RefId)
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
	return &content.MomentListRep{
		Total: total,
		List:  moments,
		Users: users,
	}, nil
}

func (*MomentService) Delete(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	contentDBDao := data.GetDBDao(ctxi, db)

	err = contentDBDao.DelByAuth(model.TableNameMoment, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
