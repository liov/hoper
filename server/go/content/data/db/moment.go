package db

import (
	clausei "github.com/hopeio/gox/database/sql/gorm/clause"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/model"
	"github.com/liov/hoper/server/go/protobuf/content"
	"gorm.io/gorm/clause"
)

func (d *ContentDao) GetMomentList(req *content.MomentListReq) (int64, []*content.Moment, error) {
	ctxi := d.Context
	var moments []*content.Moment
	db := d.db.Table(model.TableNameMoment)
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, ctxi.RespErrorLog(errcode.DBError, err, "Count")
	}
	var clauses []clause.Expression
	clauses = append(clauses, clausei.PageExpr(int(req.PageNo), int(req.PageSize)))
	err = db.Clauses(clauses...).Order("created_at desc").Find(&moments).Error
	if err != nil {
		return 0, nil, ctxi.RespErrorLog(errcode.DBError, err, "GetMomentList")
	}
	return count, moments, nil
}
