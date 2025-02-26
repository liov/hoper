package global

import (
	"github.com/hopeio/initialize/dao/gormdb/postgres"
	initredis "github.com/hopeio/initialize/dao/redis"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	GORMDB postgres.DB
	// RedisPool Redis连接池
	Redis initredis.Client
}

func (d *dao) BeforeInject() {

}

func (d *dao) AfterInjectConfig() {

}

func (d *dao) AfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

}
