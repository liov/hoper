package dao

import (
	"github.com/cockroachdb/pebble"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"github.com/liov/hoper/go/v2/upload/config"
	gormi "github.com/liov/hoper/go/v2/utils/dao/db/gorm"
	"github.com/liov/hoper/go/v2/utils/log"
	"gorm.io/gorm"
)

var Dao *dao = &dao{}

type uploadDao struct {
	ctxi *user.Ctx
}

func GetDao(ctx *user.Ctx) *uploadDao {
	if ctx == nil{
		log.Fatal("ctx can't nil")
	}
	return &uploadDao{ctx}
}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *gorm.DB
	PebbleDB *pebble.DB
	// RedisPool Redis连接池
	Redis       *redis.Client
	Cache       *ristretto.Cache

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

func (d *dao) Custom() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")


}

func (d *dao) GetDB(log *log.Logger) *gorm.DB {
	if initialize.InitConfig.Env == initialize.DEVELOPMENT{
		return d.GORMDB
	}
	return d.GORMDB.Session(&gorm.Session{
		Logger: &gormi.SQLLogger{Logger: log.Logger,
			Config: &config.Conf.Database.Gorm.Logger,
		}})
}

