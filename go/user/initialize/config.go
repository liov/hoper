package initialize

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/user/internal/config"
	"github.com/liov/hoper/go/utls/log"
)

func (i *Init) config() {
	confUrl := flag.String("c", "./config/config.toml", "配置文件路径")
	err := configor.New(&configor.Config{Debug: false}).Load(i, *confUrl)
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).Load(config.Config, *confUrl)
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}

	if runtime.GOOS == "windows" {
		config.Config.Server.LuosimaoAPIKey = ""
		config.Config.Redis.Password = ""
		config.Config.Server.Env = DEVELOPMENT
	} else {
		flag.StringVar(&config.Config.Database.Password, "p", config.Config.Database.Password, "password")
		flag.StringVar(&config.Config.Server.MailPassword, "mp", config.Config.Server.MailPassword, "password")
		flag.Parse()
		config.Config.Redis.Password = config.Config.Database.Password
		config.Config.Server.Env = i.Env
	}

	config.Config.Server.UploadMaxSize = config.Config.Server.UploadMaxSize * 1024 * 1024
	config.Config.Server.ReadTimeout = config.Config.Server.ReadTimeout * time.Second
	config.Config.Server.WriteTimeout = config.Config.Server.WriteTimeout * time.Second
	config.Config.Redis.IdleTimeout = config.Config.Redis.IdleTimeout * time.Second
}
