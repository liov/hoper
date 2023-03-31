package dao

import (
	"github.com/hopeio/pandora/context/http_context"
	"github.com/hopeio/pandora/initialize/gormdb/postgres"
	"github.com/hopeio/pandora/initialize/pebble"
	initredis "github.com/hopeio/pandora/initialize/redis"
	"github.com/hopeio/pandora/initialize/ristretto"
	"github.com/hopeio/pandora/utils/log"
)

var Dao *dao = &dao{}

type uploadDao struct {
	*http_context.Context
}

func GetDao(ctx *http_context.Context) *uploadDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &uploadDao{ctx}
}

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
