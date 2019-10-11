package config

import (
	"flag"

	"github.com/liov/hoper/go/v2/initialize/dao"
	"github.com/liov/hoper/go/v2/utils/db/gormCallback"
)

func CustomInit() {
	flag.Parse()
	db()
}

func db() {
	db := dao.Dao.DB
	db.Callback().Create().Remove("gorm:update_time_stamp")
	db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Create().Replace("gorm:save_before_associations", saveBeforeAssociationsCallback)
	//db.Callback().Create().Replace("gorm:save_after_associations", saveAfterAssociationsCallback)
	db.Callback().Delete().Replace("gorm:delete", gormCallback.DeleteCallback)
}
