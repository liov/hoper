package pro

import (
	"github.com/actliboy/hoper/server/go/lib/initialize/cache/ristretto"
	initpostgres "github.com/actliboy/hoper/server/go/lib/initialize/db/postgres"
	v8 "github.com/actliboy/hoper/server/go/lib/initialize/elastic/v8"
	initializeredis "github.com/actliboy/hoper/server/go/lib/initialize/redis"
	"os"
	"time"
)

type Config struct {
	Pro Customize
}

func (c *Config) Init() {
	c.Pro.Interval = c.Pro.Interval * time.Millisecond
	c.Pro.Timer = c.Pro.Timer * time.Second
	os.MkdirAll(c.Pro.CommonDir, 0666)
	if c.Pro.CommonDir[len(c.Pro.CommonDir)-1] != '/' {
		c.Pro.CommonDir += "/"
	}
	CommonDirLen = len(c.Pro.CommonDir)

	/*	WithClient(http.DefaultClient,30,`socks5://localhost:8080`)
		WithClient(picClient,30,`socks5://localhost:8080`)*/
}

type Customize struct {
	CommonUrl string
	Loop      int
	CommonDir string
	Ext       string
	Timer     time.Duration
	Interval  time.Duration
	WorkCount uint
	StopTime  time.Duration
}

type TimepillDao struct {
	Hoper initpostgres.DB
	Redis initializeredis.Redis
	Cache ristretto.Cache
}

type dao struct {
	DB    initpostgres.DB `init:"config:hoper"`
	Redis initializeredis.Redis
	Cache ristretto.Cache
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
