package confdao

import (
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	initredis "github.com/hopeio/cherry/initialize/conf_dao/redis"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	GORMDB postgres.DB
	// RedisPool Redis连接池
	Redis initredis.Client
}

func (d *dao) InitBeforeInject() {

}

func (d *dao) InitAfterInjectConfig() {

}

func (d *dao) InitAfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

}
