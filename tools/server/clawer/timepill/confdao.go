package timepill

import (
	"context"
	"encoding/base64"
	"github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	initializeredis "github.com/liov/hoper/server/go/lib/initialize/redis"
	"github.com/liov/hoper/server/go/lib/initialize/ristretto"
	"time"
)

type Config struct {
	TimePill Customize
}

func (c *Config) Init() {
	Token = "Basic " + base64.StdEncoding.EncodeToString([]byte(c.TimePill.User+":"+c.TimePill.Password))
}

type Customize struct {
	User        string
	Password    string
	PhotoPath   string
	PhotoPrefix string
	SearchHost  string
	PageSize    int
	Timer       time.Duration
}

type dao struct {
	Hoper     postgres.DB
	Redis     initializeredis.Redis
	Cache     ristretto.Cache
	UserCache ristretto.Cache
	//NsqP      insq.Producer `init:"notInject;config:nsq-producer"`
}

func (dao *dao) Init() {
}

func (dao *dao) Close() {
}

var (
	Dao   dao
	Conf  Config
	Token string
)

func (dao *dao) DBDao(ctx context.Context) *DBDao {
	return &DBDao{ctx: ctx, Hoper: dao.Hoper.DB}
}

func (dao *dao) RedisDao(ctx context.Context) *TimepillRedisDao {
	return &TimepillRedisDao{ctx: ctx, Redis: dao.Redis.Client}
}
