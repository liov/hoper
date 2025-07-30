package global

import (
	"database/sql"
	"github.com/hopeio/initialize/dao/gormdb/postgres"
	initredis "github.com/hopeio/initialize/dao/redis"
	"github.com/hopeio/gox/log"
	"github.com/liov/hoper/server/go/protobuf/content"
)

// 原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
// 其实init主要就是配置文件数据库连接，可以理解为init放进dao

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB postgres.DB
	StdDB  *sql.DB
	// RedisPool Redis连接池
	Redis initredis.Client
	//elastic
}

func (d *dao) BeforeInject() {

}

func (d *dao) AfterInjectConfig() {

}

func (d *dao) AfterInject() {
	d.GORMDB.Conf.NamingStrategy.TablePrefix = "content."
	d.GORMDB.NamingStrategy = d.GORMDB.Conf.NamingStrategy
	err := d.GORMDB.Exec(`CREATE SCHEMA IF NOT EXISTS "content"`).Error
	if err != nil {
		log.Fatal(err)
	}
	err = d.GORMDB.Migrator().AutoMigrate(&content.Content{}, &content.Container{},
		&content.ContentAttr{}, &content.Comment{},
		&content.Favorite{}, &content.FavFollow{}, &content.UserStatistics{},
		&content.Statistics{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
