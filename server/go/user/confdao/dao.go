package confdao

import (
	"database/sql"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/postgres"
	"github.com/hopeio/cherry/initialize/conf_dao/mail"
	"github.com/hopeio/cherry/initialize/conf_dao/redis"
	"github.com/hopeio/cherry/initialize/conf_dao/ristretto"
	"github.com/liov/hoper/server/go/protobuf/user"
	"gorm.io/gorm"
)

// 原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
// 其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB postgres.DB
	StdDB  *sql.DB
	// RedisPool Redis连接池
	Redis redis.Client
	Cache ristretto.Cache[string, any]
	//elastic
	Mail mail.Mail `init:"config:mail"`
}

func (d *dao) InitBeforeInject() {

}

func (d *dao) InitAfterInjectConfig() {
}

func (d *dao) InitAfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB.DB()
}

func Migrator(db *gorm.DB) {
	db.Migrator().AutoMigrate(&user.User{}, &user.Resume{}, &user.UserActionLog{}, &user.UserBannedLog{}, &user.UserDeviceInfo{}, &user.UserScoreLog{}, &user.UserExt{})
}
