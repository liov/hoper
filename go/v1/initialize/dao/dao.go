package dao

import (
	"github.com/bluele/gcache"
	"github.com/etcd-io/bbolt"
	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
)

//原本是个单独模块，但是考虑到数据库必须初始化，所以合进来了
//其实init主要就是配置文件数据库连接，可以理解为init放进dao
var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// DB 数据库连接
	DB   *gorm.DB
	Bolt *bbolt.DB
	// RedisPool Redis连接池
	Redis *redis.Pool
	// MongoDB 数据库连接
	Mongo       *mgo.Database
	RedisExpire int32
	Cache       gcache.Cache
	McExpire    int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Close close the resource.
func (d *dao) Close() {
	if d.Bolt != nil {
		d.Bolt.Close()
	}
	d.Redis.Close()
	d.DB.Close()
	if d.Mongo != nil {
		d.Mongo.Session.Close()
	}
}

func SetDB(db *gorm.DB) {
	Dao.DB = db
}

func SetCache(c gcache.Cache) {
	Dao.Cache = c
}

func SetRedis(redb *redis.Pool) {
	Dao.Redis = redb
}
