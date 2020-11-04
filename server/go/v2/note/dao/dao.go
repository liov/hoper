package dao

import (
	"database/sql"
	"net/smtp"

	"github.com/bluele/gcache"
	"github.com/etcd-io/bbolt"
	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/utils/apollo"
	"gorm.io/gorm"
)

//原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
//其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB *gorm.DB
	StdDB  *sql.DB
	Bolt   *bbolt.DB
	// RedisPool Redis连接池
	Redis       *redis.Pool
	RedisExpire int32
	Cache       gcache.Cache
	McExpire    int32
	//elastic
	MailAuth smtp.Auth
	Apollo   *apollo.Client
}

// CloseDao close the resource.
func (d *dao) Close() {
	if d.Bolt != nil {
		d.Bolt.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.StdDB != nil {
		d.StdDB.Close()
	}
	if d.Apollo != nil {
		d.Apollo.Close()
	}
}

func (d *dao) Custom() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:force_reload_after_create")
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")
	d.StdDB, _ = db.DB()
}
