package service

import (
	"context"
	"net/http"

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

func (*MomentService) Info(context.Context, *content.GetMomentReq) (*content.Moment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthInfo not implemented")
}

func (m *MomentService) Add(ctx context.Context, req *content.AddMomentReq) (*request.Empty, error) {
	ctxi, span := user.CtxFromContext(ctx).StartSpan("")
	defer span.End()
	user, err := ctxi.GetAuthInfo(AuthWithUpdate)
	if err != nil {
		return nil, err
	}
	err = dao.Dao.Limit(ctxi, &conf.Conf.Customize.Limit)
	if err != nil {
		return nil, err
	}
	log.Info(user)
	req.UserId = user.Id
	db := dao.Dao.GORMDB
	/*	var count int64
		db.Table(`mood`).Where(`name = ?`, req.MoodName).Count(&count)
		if count == 0 {
			return nil, errorcode.ParamInvalid.Message("心情不存在")
		}*/
	var tags []model.TinyTag
	db.Table("tag").Select("id,name").
		Where("name IN (?)", req.Tags).Find(&tags)

	req.UserId = user.Id
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(req).Error
		if err != nil {
			return ctxi.Log(err, "db.CreateReq", err.Error())
		}
		var contentTags []model.ContentTag
		var noExist []content.Tag
		for i := range req.Tags {
			// 性能可以优化
			for j := range tags {
				if req.Tags[i] == tags[j].Name {
					contentTags = append(contentTags, model.ContentTag{
						Type:  content.ContentMoment,
						RefId: req.Id,
						TagId: tags[i].Id,
					})
					break
				} else {
					noExist = append(noExist, content.Tag{Name: req.Tags[i], UserId: user.Id})
				}
			}
			if err = tx.Create(&noExist).Error; err != nil {
				return ctxi.Log(err, "db.CreateNoExist", err.Error())
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
		}
		return nil
	})
	if err != nil {
		return nil, ctxi.Log(errorcode.DBError, "db.Transaction", err.Error())
	}
	return nil, nil
}
func (*MomentService) Edit(context.Context, *content.AddMomentReq) (*request.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}

func (*MomentService) List(context.Context, *content.MomentListReq) (*content.MomentListRep, error) {

	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
