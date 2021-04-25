package service

import (
	"context"
	"fmt"
	"github.com/liov/hoper/go/v2/content/client"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"net/http"
	"unicode/utf8"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
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

func (*MomentService) Info(ctx context.Context, req *request.Object) (*content.Moment, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)

	db := dao.Dao.GetDB(ctxi.Logger)
	var moment content.Moment
	err = db.Table(model.MomentTableName).
		Where(`id = ?`, req.Id).First(&moment).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "First")
	}
	// tags
	contentTags, err := contentDao.GetContentTagDB(db, content.ContentMoment, []uint64{moment.Id})
	if err != nil {
		return nil, err
	}
	var tags = make([]*content.TinyTag, len(contentTags))
	for i := range contentTags {
		tags[i] = &contentTags[i].TinyTag
	}
	moment.Tags = tags

	//like
	if auth.Id != 0 {
		likes, err := contentDao.GetContentActionsDB(db, content.ActionLike, content.ContentMoment, []uint64{req.Id}, auth.Id)
		if err != nil {
			return nil, err
		}

		for i := range likes {
			if likes[i].Action == content.ActionLike {
				moment.LikeId = likes[i].Id
			}
			if likes[i].Action == content.ActionUnlike {
				moment.UnlikeId = likes[i].Id
			}
		}
		collects, err := contentDao.GetCollectsDB(db, content.ContentMoment, []uint64{req.Id}, auth.Id)
		if err != nil {
			return nil, err
		}
		for i := range collects {
			moment.Collects = append(moment.Collects, collects[i].FavId)
		}
	}
	// ext
	exts, err := contentDao.GetContentExtDB(db, content.ContentMoment, []uint64{moment.Id})
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

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Object, error) {

	if utf8.RuneCountInString(req.Content) < conf.Conf.Customize.Moment.MaxContentLen {
		return nil, errorcode.InvalidArgument.Message(fmt.Sprintf("文章内容不能小于%d个字", conf.Conf.Customize.Moment.MaxContentLen))
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

	req.UserId = auth.Id
	db := dao.Dao.GetDB(ctxi.Logger)
	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errorcode.ParamInvalid.Message("心情不存在")
		}*/
	var tags []model.TinyTag
	if len(req.Tags) > 0 {
		tags, err = contentDao.GetTagsDB(db, req.Tags)
		if err != nil {
			return nil, err
		}
	}

	req.UserId = auth.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		if req.Permission == 0 {
			req.Permission = content.ViewPermissionAll
		}
		err = tx.Table(model.MomentTableName).Create(req).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "tx.CreateReq")
		}
		err = tx.Table(model.ContentExtTableName).Create(&model.ContentExt{
			Type:  content.ContentMoment,
			RefId: req.Id,
		}).Error
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "tx.CreateReq")
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
		if err != errorcode.DBError {
			return nil, ctxi.ErrorLog(errorcode.DBError, err, "Transaction")
		}
		return nil, errorcode.DBError
	}
	return &request.Object{Id: req.Id}, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(ctx context.Context, req *content.MomentListReq) (*content.MomentListRep, error) {
	ctxi, span := contexti.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := auth(ctxi, true)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)

	db := dao.Dao.GetDB(ctxi.Logger)

	total, moments, err := contentDao.GetMomentListDB(db, req)
	if err != nil {
		return nil, err
	}
	var m = make(map[uint64]*content.Moment)
	var ids []uint64
	var userIds []uint64
	for i := range moments {
		ids = append(ids, moments[i].Id)
		m[moments[i].Id] = moments[i]
		userIds = append(userIds, moments[i].UserId)
		// 屏蔽字段
		momentMaskField(moments[i])
	}
	// tag
	tags, err := contentDao.GetContentTagDB(db, content.ContentMoment, ids)
	if err != nil {
		return nil, err
	}

	for i := range tags {
		if moment, ok := m[tags[i].RefId]; ok {
			moment.Tags = append(moment.Tags, &tags[i].TinyTag)
		}
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
		likes, err := contentDao.GetContentActionsDB(db, content.ActionLike, content.ContentMoment, ids, auth.Id)
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
		collects, err := contentDao.GetCollectsDB(db, content.ContentMoment, ids, auth.Id)
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
	return &content.MomentListRep{
		Total: total,
		List:  moments,
		Users: users,
	}, nil
}

func (*MomentService) Delete(ctx context.Context, req *request.Object) (*empty.Empty, error) {
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
	err = contentDao.DelByAuthDB(db, model.MomentTableName, req.Id, auth.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
