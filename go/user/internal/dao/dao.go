package dao

import (
	"github.com/bluele/gcache"
	"github.com/etcd-io/bbolt"
	"github.com/garyburd/redigo/redis"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
)

var Dao *dao = &dao{}

// dao dao.
type dao struct {
	// DB 数据库连接
	db   *gorm.DB
	bolt *bbolt.DB
	// RedisPool Redis连接池
	redis *redis.Pool
	// MongoDB 数据库连接
	mongo       *mgo.Database
	redisExpire int32
	cache       gcache.Cache
	mcExpire    int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Close close the resource.
func (d *dao) Close() {
	if d.bolt != nil {
		d.bolt.Close()
	}
	d.redis.Close()
	d.db.Close()
	if d.mongo != nil {
		d.mongo.Session.Close()
	}
}

func SetDB(db *gorm.DB) {
	Dao.db = db
}



func SetCache(c gcache.Cache) {
	Dao.cache = c
}

func SetRedis(redb *redis.Pool) {
	Dao.redis = redb
}
