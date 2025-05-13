package global

import (
	"github.com/hopeio/initialize"
	"github.com/hopeio/initialize/dao/gormdb/postgres"
	"github.com/hopeio/initialize/dao/redis"
	"github.com/hopeio/initialize/dao/ristretto"
)

// dao dao.
type dao struct {
	initialize.EmbeddedPresets
	// GORMDB 数据库连接
	GORMDB postgres.DB
	// RedisPool Redis连接池
	Redis redis.Client
	Cache ristretto.Cache[int, any]
}

func (d *dao) AfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

}
