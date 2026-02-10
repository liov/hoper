package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/scaffold/errcode"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"unicode/utf8"

	"github.com/hopeio/gox/container/set"
	"github.com/hopeio/protobuf/request"
	"github.com/liov/hoper/server/go/protobuf/common"
	"google.golang.org/protobuf/types/known/emptypb"

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
	return "瞬间相关", "/api/moment", nil
}

func (*MomentService) Info(ctx context.Context, req *request.Id) (*content.Moment, error) {
	ctx, span := Tracer.Start(ctx, "Content.Info")
	defer span.End()
	auth, _ := auth(ctx, true)

	contentDBDao := data.GetDBDao(ctx, global.Dao.GORMDB.DB)

	var moment content.Moment
	err := contentDBDao.Table(model.TableNameMoment).First(&moment, req.Id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errcode.NotFound
		}
		log.Errorw("Find", zap.Error(err))
		return nil, errcode.DBError.Wrap(err)
	}
	// tags
	contentTags, err := contentDBDao.GetContentTag(content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
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
			return nil, errcode.DBError.Wrap(err)
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
			return nil, errcode.DBError.Wrap(err)
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
		return nil, errcode.DBError.Wrap(err)
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
		}
	}
	momentMaskField(&moment)
	return &moment, nil
}

// 屏蔽字段
func momentMaskField(moment *content.Moment) {
	moment.ModelTime.DeletedAt = nil
	moment.Anonymous = 0
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Id, error) {

	if utf8.RuneCountInString(req.Content) < global.Conf.Customize.Moment.MaxContentLen {
		return nil, errcode.InvalidArgument.Msg(fmt.Sprintf("文章内容不能小于%d个字", global.Conf.Customize.Moment.MaxContentLen))
	}

	ctx, span := Tracer.Start(ctx, "Content.Add")
	defer span.End()
	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}
	contentDBDao := data.GetDBDao(ctx, comconfdao.Dao.GORMDB.DB)
	commonDBDao := comdata.GetDBDao(ctx, comconfdao.Dao.GORMDB.DB)

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
			return nil, errcode.DBError.Wrap(err)
		}
	}

	req.UserId = auth.Id
	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		if req.Permission == 0 {
			req.Permission = content.ViewPermissionAll
		}
		contenttxDBDao := data.GetDBDao(ctx, tx)
		err = tx.Table(model.TableNameMoment).Create(req).Error
		if err != nil {
			return errcode.DBError.Wrap(err)
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
				return errcode.DBError.Wrap(err)
			}
		}
		if len(noExist) > 1 {
			if err = tx.Create(&noExist).Error; err != nil {
				return errcode.DBError.Wrap(err)
			}
		}
		for i := range noExist {
			contentTags = append(contentTags, model.ContentTag{
				Type:  content.ContentMoment,
				RefId: req.Id,
				TagId: noExist[i].Basic.Id,
			})
		}
		if len(contentTags) > 0 {
			if err = tx.Create(&contentTags).Error; err != nil {
				return errcode.DBError.Wrap(err)
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

func (*MomentService) List(ctx context.Context, req *content.MomentListReq) (*content.MomentListResp, error) {
	ctx, span := Tracer.Start(ctx, "Content.List")
	defer span.End()
	auth, _ := auth(ctx, true)

	contentDBDao := data.GetDBDao(ctx, global.Dao.GORMDB.DB)

	total, moments, err := contentDBDao.GetMomentList(req)
	if err != nil {
		return nil, err
	}

	if len(moments) == 0 {
		return &content.MomentListResp{
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
		userList, err := global.UserClient().BaseList(metadata.AppendToOutgoingContext(ctx, httpx.HeaderGrpcInternal, httpx.HeaderGrpcInternal), &user.BaseListReq{Ids: userIds.ToSlice()})
		if err != nil {
			return nil, err
		}
		users = userList.List
	}
	return &content.MomentListResp{
		Total: total,
		List:  moments,
		Users: users,
	}, nil
}

func (*MomentService) Delete(ctx context.Context, req *request.Id) (*emptypb.Empty, error) {
	ctx, span := Tracer.Start(ctx, "Content.Delete")
	defer span.End()
	auth, err := auth(ctx, true)
	if err != nil {
		return nil, err
	}

	contentDBDao := data.GetDBDao(ctx, global.Dao.GORMDB.DB)

	err = contentDBDao.DelByAuth(model.TableNameMoment, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
