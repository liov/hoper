package dao

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	"github.com/liov/hoper/server/go/lib/initialize/mail"
	"github.com/liov/hoper/server/go/lib/initialize/pebble"
	initredis "github.com/liov/hoper/server/go/lib/initialize/redis"
	"github.com/liov/hoper/server/go/lib/initialize/ristretto"
	"github.com/liov/hoper/server/go/mod/content/dao/db"
	rdao "github.com/liov/hoper/server/go/mod/content/dao/redis"

	"gorm.io/gorm"
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

func GetDBDao(ctx *http_context.Context, d *gorm.DB) *db.ContentDBDao {
	return db.GetDao(ctx, d)
}

func GetRedisDao(ctx *http_context.Context, r redis.Cmdable) *rdao.ContentRedisDao {
	return rdao.GetDao(ctx, r)
}
