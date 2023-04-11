package confdao

import (
	"github.com/hopeio/pandora/initialize/gormdb/postgres"
	"github.com/hopeio/pandora/initialize/pebble"
	initredis "github.com/hopeio/pandora/initialize/redis"
	"github.com/hopeio/pandora/initialize/ristretto"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	GORMDB   postgres.DB
	PebbleDB pebble.DB
	// RedisPool Redis连接池
	Redis initredis.Redis
	Cache ristretto.Cache
}

func (d *dao) Init() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

}
