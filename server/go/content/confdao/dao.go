package confdao

import (
	"database/sql"
	"github.com/hopeio/tiga/initialize/basic_dao/gormdb/postgres"
	"github.com/hopeio/tiga/initialize/basic_dao/mail"
	"github.com/hopeio/tiga/initialize/basic_dao/pebble"
	initredis "github.com/hopeio/tiga/initialize/basic_dao/redis"
	"github.com/hopeio/tiga/initialize/basic_dao/ristretto"
	"github.com/hopeio/tiga/initialize/basic_dao/viper"
)

// 原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
// 其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   postgres.DB
	StdDB    *sql.DB
	PebbleDB pebble.DB
	// RedisPool Redis连接池
	Redis initredis.Redis
	Cache ristretto.Cache
	//elastic
	MailAuth mail.Mail `init:"config:mail"`
	Viper    viper.Viper
}

func (d *dao) Init() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:force_reload_after_create")
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")
	d.StdDB, _ = db.DB.DB()
}
