package pro

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/cache_ristretto"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	v8 "github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic/v8"
	initializeredis "github.com/actliboy/hoper/server/go/lib/tiga/initialize/redis"
	"os"
	"time"
)

type Config struct {
	Pro Customize
}

func (c *Config) Init() {
	c.Pro.Interval = c.Pro.Interval * time.Millisecond
	c.Pro.Timer = c.Pro.Timer * time.Second
	CommonDirLen = len(c.Pro.CommonDir)

	/*	SetClient(http.DefaultClient,30,`socks5://localhost:8080`)
		SetClient(picClient,30,`socks5://localhost:8080`)*/
}

type Customize struct {
	CommonUrl string
	Loop      int
	CommonDir string
	Ext       string
	Timer     time.Duration
	Interval  time.Duration
}

type TimepillDao struct {
	Hoper db.DB
	Redis initializeredis.Redis
	Cache cache_ristretto.Cache
}

type dao struct {
	DB    db.DB `init:"config:hoper"`
	Redis initializeredis.Redis
	Cache cache_ristretto.Cache
	Es8   v8.Es `init:"config:elasticsearch8"`
}

func (dao *dao) Init() {
}

func (dao *dao) Close() {
}

var (
	Dao          dao
	Conf         Config
	CommonDirLen int
)

const Sep = string(os.PathSeparator)
