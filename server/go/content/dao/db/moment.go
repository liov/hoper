package db

import (
	"github.com/actliboy/hoper/server/go/content/model"
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/hopeio/pandora/protobuf/errorcode"
	dbi "github.com/hopeio/pandora/utils/dao/db/const"
	clausei "github.com/hopeio/pandora/utils/dao/db/gorm/clause"
	"gorm.io/gorm/clause"
)

func (d *ContentDBDao) GetMomentListDB(req *content.MomentListReq) (int64, []*content.Moment, error) {
	ctxi := d.Context
	var moments []*content.Moment
	db := d.db.Table(model.MomentTableName).Where(dbi.NotDeleted)
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Count")
	}
	var clauses []clause.Expression
	clauses = append(clauses, clausei.Page(int(req.PageNo), int(req.PageSize)))
	err = db.Clauses(clauses...).Order("created_at desc").Find(&moments).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "GetMomentListDB")
	}
	return count, moments, nil
}
