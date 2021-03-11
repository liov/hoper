package service

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/model"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MomentService struct {
	content.UnimplementedMomentServiceServer
}

func (m *MomentService) Service() (describe, prefix string, middleware []http.HandlerFunc) {
	return "瞬间相关", "/api/moment", nil
}

func (m *MomentService) Info(ctx context.Context, req *content.GetMomentReq) (*content.Moment, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	_, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, err
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	var moment content.Moment
	err = db.Table(model.MomentTableName).
		Where(`id = ?`, req.Id).First(&moment).Error
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "First", err.Error())
	}
	tags, err := contentDao.GetTagsByRefIdDB(db, content.ContentMoment, moment.Id)
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "dbDao.GetTagContent", err.Error())
	}
	moment.Tags = tags

	return &moment, nil
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Empty, error) {

	if utf8.RuneCountInString(req.Content) < conf.Conf.Customize.Moment.MaxContentLen {
		return nil, errorcode.ParamInvalid.Message(fmt.Sprintf("文章内容不能小于%d个字",conf.Conf.Customize.Moment.MaxContentLen))
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

	req.UserId = auth.Id
	db := dao.Dao.GetDB(ctxi.Logger)
	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errorcode.ParamInvalid.Message("心情不存在")
		}*/
	tags, err := contentDao.GetTagsDB(db, req.Tags)
	if err != nil {
		return nil, ctxi.Log(errorcode.RedisErr, "GetTagsDB", err.Error())
	}
	req.UserId = auth.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		if req.Permission == 0 {
			req.Permission = content.ViewPermissionAll
		}
		err = tx.Create(req).Error
		if err != nil {
			return ctxi.Log(err, "db.CreateReq", err.Error())
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
				return ctxi.Log(err, "db.CreateNoExist", err.Error())
			}
		}
		if len(noExist) > 1 {
			if err = tx.Create(&noExist).Error; err != nil {
				return ctxi.Log(err, "db.CreateNoExist", err.Error())
			}
		}
		for i := range noExist {
			contentTags = append(contentTags, model.ContentTag{
				Type:  content.ContentMoment,
				RefId: req.Id,
				TagId: noExist[i].Id,
			})
		}
		if err = tx.Create(&contentTags).Error; err != nil {
			return ctxi.Log(err, "db.CreateContentTags", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errorcode.DBError
	}
	return nil, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*request.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(ctx context.Context, req *content.MomentListReq) (*content.MomentListRep, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, ctxi.Log(errorcode.RedisErr, "LimitRedis", err.Error())
	}
	log.Info(auth)

	db := dao.Dao.GetDB(ctxi.Logger)

	var moments []*content.Moment
	moments, err = contentDao.GetMomentListDB(db, req)
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "dbDao.GetMomentList", err.Error())
	}
	var ids []uint64
	for i := range moments {
		ids = append(ids, moments[i].Id)
	}
	tags, err := contentDao.GetTagContentDB(db, content.ContentMoment, ids)
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "dbDao.GetTagContent", err.Error())
	}
	var m = make(map[uint64][]*content.TinyTag)
	for i := range tags {
		m[tags[i].RefId] = append(m[tags[i].RefId], &tags[i].TinyTag)
	}
	for i := range moments {
		moments[i].Tags = m[moments[i].Id]
	}
	return &content.MomentListRep{
		Count: 0,
		List:  moments,
	}, nil
}

func (*MomentService) Delete(ctx context.Context, req *content.GetMomentReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	auth, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	contentDao := dao.GetDao(ctxi)
	err = contentDao.LimitRedis(dao.Dao.Redis, &conf.Conf.Customize.Moment.Limit)
	if err != nil {
		return nil, ctxi.Log(errorcode.RedisErr, "LimitRedis", err.Error())
	}
	db := dao.Dao.GetDB(ctxi.Logger)
	err = contentDao.DeleteMomentDB(db, req.Id, auth.Id)
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "DeleteMomentDB", err.Error())
	}
	return nil, nil
}
