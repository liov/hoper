package service

import (
	"context"
	"fmt"
	"github.com/hopeio/pandora/context/http_context"
	"github.com/hopeio/pandora/protobuf/request"
	"github.com/liov/hoper/server/go/mod/content/client"
	"net/http"
	"unicode/utf8"

	"github.com/hopeio/pandora/protobuf/empty"
	"github.com/hopeio/pandora/protobuf/errorcode"
	"github.com/liov/hoper/server/go/mod/content/confdao"
	"github.com/liov/hoper/server/go/mod/content/dao"
	"github.com/liov/hoper/server/go/mod/content/model"
	"github.com/liov/hoper/server/go/mod/protobuf/content"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MomentService struct {
	content.UnimplementedMomentServiceServer
}

func (*MomentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "瞬间相关", "/api/moment", nil
}

func (*MomentService) Info(ctx context.Context, req *request.Id) (*content.Moment, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, _ := auth(ctxi, true)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	contentDBDao := dao.GetDBDao(ctxi, db)

	var moment content.Moment
	err := db.Table(model.MomentTableName).
		Where(`id = ?`, req.Id).First(&moment).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "First")
	}
	// tags
	contentTags, err := contentDBDao.GetContentTag(content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, err
	}
	var tags = make([]*content.TinyTag, len(contentTags))
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
			action.Collects = append(action.Collects, collects[i].FavId)
		}
		if len(likes) != 0 && len(collects) != 0 {
			moment.Action = action
		}
	}
	// ext
	exts, err := contentDBDao.GetContentExt(content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, err
	}
	moment.Ext = exts[0]

	var userIds []uint64

	// 匿名
	if moment.Anonymous == 1 {
		moment.UserId = 0
	} else {
		userIds = append(userIds, moment.UserId)
	}
	if len(userIds) > 0 {
		userList, err := client.UserClient.BaseList(ctxi, &user.BaseListReq{Ids: userIds})
		if err != nil {
			return nil, err
		}
		/*	var m = make(map[uint64]*user.UserBaseInfo)
			for _,u:=range userList.List{
				m[u.Id] = u
			}
			// 这个可以放到前端做，减少数据返回
			for i := range comments{
				comments[i].RecvUser = m[comments[i].RecvId]
				comments[i].User = m[comments[i].UserId]
			}
			moment.User = m[moment.UserId]*/
		// 客户端组装
		moment.Users = userList.List
	}

	momentMaskField(&moment)
	return &moment, nil
}

// 屏蔽字段
func momentMaskField(moment *content.Moment) {
	moment.AreaVisibility = 0
	moment.DeletedAt = ""
	moment.CreatedAt = moment.CreatedAt[:19]
	moment.Anonymous = 0
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Id, error) {

	if utf8.RuneCountInString(req.Content) < confdao.Conf.Customize.Moment.MaxContentLen {
		return nil, errorcode.InvalidArgument.Message(fmt.Sprintf("文章内容不能小于%d个字", confdao.Conf.Customize.Moment.MaxContentLen))
	}

	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	contentDBDao := dao.GetDBDao(ctxi, db)

	req.UserId = auth.Id

	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errorcode.ParamInvalid.Message("心情不存在")
		}*/
	var tags []model.TinyTag
	if len(req.Tags) > 0 {
		tags, err = contentDBDao.GetTags(req.Tags)
		if err != nil {
			return nil, err
		}
	}

	req.UserId = auth.Id
	err = contentDBDao.Transaction(func(tx *gorm.DB) error {
		if req.Permission == 0 {
			req.Permission = content.ViewPermissionAll
		}
		contenttxDBDao := dao.GetDBDao(ctxi, tx)
		err = tx.Table(model.MomentTableName).Create(req).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "tx.CreateReq")
		}
		err = contenttxDBDao.CreateContextExt(content.ContentMoment, req.Id)
		if err != nil {
			return err
		}
		var contentTags []model.ContentTag
		var noExist []content.Tag
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
			noExist = append(noExist, content.Tag{Name: req.Tags[i], UserId: auth.Id})
		}
		if len(noExist) == 1 {
			if err = tx.Create(&noExist[1]).Error; err != nil {
				return ctxi.ErrorLog(errorcode.DBError, err, "db.CreateNoExist")
			}
		}
		if len(noExist) > 1 {
			if err = tx.Create(&noExist).Error; err != nil {
				return ctxi.ErrorLog(errorcode.DBError, err, "db.CreateNoExist")
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
				return ctxi.ErrorLog(errorcode.DBError, err, "db.CreateContentTags")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &request.Id{Id: req.Id}, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(ctx context.Context, req *content.MomentListReq) (*content.MomentListRep, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, _ := auth(ctxi, true)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	contentDBDao := dao.GetDBDao(ctxi, db)

	total, moments, err := contentDBDao.GetMomentListDB(req)
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
	var ids, userIds []uint64

	for i := range moments {
		ids = append(ids, moments[i].Id)
		m[moments[i].Id] = moments[i]
		userIds = append(userIds, moments[i].UserId)
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
	exts, err := contentDBDao.GetContentExt(content.ContentMoment, ids)
	if err != nil {
		return nil, err
	}
	for i := range exts {
		if moment, ok := m[exts[i].RefId]; ok {
			moment.Ext = exts[i]
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
				moment.Action.Collects = append(moment.Action.Collects, collects[i].FavId)
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
	return &content.MomentListRep{
		Total: total,
		List:  moments,
		Users: users,
	}, nil
}

func (*MomentService) Delete(ctx context.Context, req *request.Id) (*empty.Empty, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	contentDBDao := dao.GetDBDao(ctxi, db)

	err = contentDBDao.DelByAuth(model.MomentTableName, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
