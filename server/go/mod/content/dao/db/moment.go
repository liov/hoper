package db

import (
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	clausei "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/actliboy/hoper/server/go/mod/content/model"
	"github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"gorm.io/gorm/clause"
)

func (d *ContentDBDao) GetMomentListDB(req *content.MomentListReq) (int64, []*content.Moment, error) {
	ctxi := d.Ctx
	var moments []*content.Moment
	db := d.db.Table(model.MomentTableName).Where(dbi.PostgreNotDeleted)
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
