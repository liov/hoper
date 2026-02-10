package db

import (
	clausex "github.com/hopeio/gox/database/sql/gorm/clause"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"gorm.io/gorm/clause"
)

func (d *ContentDao) GetMomentList(req *content.MomentListReq) (int64, []*content.Moment, error) {
	var moments []*content.Moment
	db := d.Table(model.TableNameMoment)
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, errcode.DBError.Wrap(err)
	}
	var clauses []clause.Expression
	clauses = append(clauses, clausex.PaginationExpr(req.PageNo, req.PageSize))
	err = db.Clauses(clauses...).Order("created_at desc").Find(&moments).Error
	if err != nil {
		return 0, nil, errcode.DBError.Wrap(err)
	}
	return count, moments, nil
}
