package confdao

import (
	"database/sql"
	"github.com/hopeio/initialize/conf_dao/gormdb/postgres"
	initredis "github.com/hopeio/initialize/conf_dao/redis"
	"github.com/hopeio/utils/log"
	"github.com/liov/hoper/server/go/protobuf/content"
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
	Redis initredis.Client
	//elastic
}

func (d *dao) InitBeforeInject() {

}

func (d *dao) InitAfterInjectConfig() {

}

func (d *dao) InitAfterInject() {
	d.GORMDB.Conf.NamingStrategy.TablePrefix = "content."
	d.GORMDB.NamingStrategy = d.GORMDB.Conf.NamingStrategy
	err := d.GORMDB.Exec(`CREATE SCHEMA IF NOT EXISTS "content"`).Error
	if err != nil {
		log.Fatal(err)
	}
	err = d.GORMDB.Migrator().AutoMigrate(&content.Moment{}, &content.Diary{}, &content.Attribute{}, &content.Tag{},
		&content.UserAction{}, &content.Article{}, &content.Category{}, &content.ContentAttr{}, &content.Comment{}, &content.AttrGroup{},
		&content.Favorite{}, &content.FavFollow{}, &content.DiaryBook{}, &content.UserStatics{}, &content.Note{}, &content.TagGroup{},
		&content.ContentExt{}, &content.TagTagGroup{}, &content.AttrAttrGroup{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
