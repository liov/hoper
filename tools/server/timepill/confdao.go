package timepill

import (
	"context"
	"encoding/base64"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache_ristretto"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	insq "github.com/actliboy/hoper/server/go/lib/tiga/initialize/nsq"
	initializeredis "github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
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
	Hoper     db.DB
	Redis     initializeredis.Redis
	Cache     cache_ristretto.Cache
	UserCache cache_ristretto.Cache
	NsqP      insq.Producer `init:"config:nsq-producer"`
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
