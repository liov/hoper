package confdao

import (
	"github.com/cockroachdb/pebble"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *gorm.DB
	PebbleDB *pebble.DB
	// RedisPool Redis连接池
	Redis *redis.Client
	Cache *ristretto.Cache
}

// CloseDao close the resource.
func (d *dao) Close() {
	if d.PebbleDB != nil {
		d.PebbleDB.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.GORMDB != nil {
		rawDB, _ := d.GORMDB.DB()
		rawDB.Close()
	}
}

func (d *dao) Init() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

}
