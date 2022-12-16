package dao

import (
	contexti "github.com/liov/hoper/server/go/lib/context"
	"github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	"github.com/liov/hoper/server/go/lib/initialize/pebble"
	initredis "github.com/liov/hoper/server/go/lib/initialize/redis"
	"github.com/liov/hoper/server/go/lib/initialize/ristretto"
	"github.com/liov/hoper/server/go/lib/utils/log"
)

var Dao *dao = &dao{}

type uploadDao struct {
	*contexti.Ctx
}

func GetDao(ctx *contexti.Ctx) *uploadDao {
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
