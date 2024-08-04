package confdao

import (
	"github.com/hopeio/initialize/conf_dao/gormdb/postgres"
	"github.com/hopeio/utils/log"
	"github.com/liov/hoper/server/go/protobuf/common"
	"github.com/liov/hoper/server/go/protobuf/content"
)

// 原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
// 其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB postgres.DB
	// RedisPool Redis连接池
	//Mail mail.Mail `init:"config:mail"`
	//elastic
}

func (d *dao) InitBeforeInject() {

}

func (d *dao) InitAfterInjectConfig() {

}

func (d *dao) InitAfterInject() {
	d.GORMDB.NamingStrategy = d.GORMDB.Conf.NamingStrategy
	err := d.GORMDB.Exec(`CREATE SCHEMA IF NOT EXISTS "content"`).Error
	if err != nil {
		log.Fatal(err)
	}
	err = d.GORMDB.Migrator().AutoMigrate(&common.Tag{},
		&common.TagGroup{}, &common.TagTagGroup{}, &content.AttrAttrGroup{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
