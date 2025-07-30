package global

import (
	"github.com/hopeio/initialize/dao/gormdb/postgres"
	"github.com/hopeio/initialize/dao/mail"
	"github.com/hopeio/initialize/dao/redis"
	"github.com/hopeio/initialize/dao/ristretto"
	"github.com/hopeio/gox/log"
	"github.com/liov/hoper/server/go/protobuf/user"
)

// 原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
// 其实init主要就是配置文件数据库连接，可以理解为init放进dao

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB postgres.DB
	// RedisPool Redis连接池
	Redis redis.Client
	Cache ristretto.Cache[string, any]
	//elastic
	Mail mail.Mail `init:"config:mail"`
}

func (d *dao) BeforeInject() {
}

func (d *dao) AfterInjectConfig() {

}
func (d *dao) AfterInject() {
	if d.GORMDB.DB != nil {
		//d.GORMDB.Conf.NamingStrategy.TablePrefix = "user."
		d.GORMDB.NamingStrategy = d.GORMDB.Conf.NamingStrategy
		err := d.GORMDB.Exec(`CREATE SCHEMA IF NOT EXISTS "user"`).Error
		if err != nil {
			log.Fatal(err)
		}
		err = d.GORMDB.Migrator().AutoMigrate(&user.User{}, &user.Resume{}, &user.ActionLog{}, &user.BannedLog{}, &user.Device{}, &user.ScoreLog{}, &user.UserExt{}, user.Oauth{})
		if err != nil {
			log.Fatal(err)
		}
	}
}
