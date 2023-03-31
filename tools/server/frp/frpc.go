package main

import (
	"flag"
	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/util/log"
	http_fs "github.com/hopeio/pandora/utils/net/http/fs"
	http_fs_watch "github.com/hopeio/pandora/utils/net/http/fs/watch"
	"os"
	"time"
)

type Config struct {
	Interval int64
	Url      string
}

func main() {
	var conf Config
	flag.StringVar(&conf.Url, "url", "https://xxx", "url")
	flag.Int64Var(&conf.Interval, "i", 30, "url")
	flag.Parse()
	log.Info("url: %s", conf.Url)
	cfgchan := make(chan Config)
	watch := http_fs_watch.New(time.Second * time.Duration(conf.Interval))
	var svr *client.Service
	callback := func(hfile *http_fs.FileInfo) {
		if svr != nil {
			svr.Close()
			time.Sleep(time.Second)
		}
		cfgFilePath := "./frpc.ini"
		file, err := os.Create(cfgFilePath)
		if err != nil {
			log.Error("open config file error: %v", err)
		}
		_, err = file.Write(hfile.Binary)
		if err != nil {
			log.Error("write file error: %v", err)
		}
		err = file.Close()
		if err != nil {
			log.Error("close file error: %v", err)
		}
		cfg, pxyCfgs, visitorCfgs, err := config.ParseClientConfig(cfgFilePath)
		if err != nil {
			log.Error("close file error: %v", err)
		}
		svr, err = startService(cfg, pxyCfgs, visitorCfgs, cfgFilePath)
		go svr.Run()
		if (cfg.AssetsDir != "" && cfg.AssetsDir != conf.Url) || cfg.HeartbeatInterval != conf.Interval {
			cfgchan <- Config{
				Interval: cfg.HeartbeatInterval,
				Url:      cfg.AssetsDir,
			}
			cfg.AssetsDir = ""
			cfg.HeartbeatInterval = 30
		}
	}
	watch.AddGet(conf.Url, callback)
	go func() {
		for newconf := range cfgchan {
			if newconf.Url != conf.Url {
				log.Info("url: %s", newconf.Url)
				watch.AddGet(newconf.Url, callback)
				watch.Remove(conf.Url)
				conf.Url = newconf.Url
			}
			if newconf.Interval != conf.Interval {
				watch.Update(time.Second * time.Duration(conf.Interval))
				conf.Interval = newconf.Interval
			}
		}
	}()
	select {}
}

func startService(
	cfg config.ClientCommonConf,
	pxyCfgs map[string]config.ProxyConf,
	visitorCfgs map[string]config.VisitorConf,
	cfgFile string,
) (svr *client.Service, err error) {

	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel,
		cfg.LogMaxDays, cfg.DisableLogColor)

	if cfgFile != "" {
		log.Trace("start frpc service for config file [%s]", cfgFile)
		defer log.Trace("frpc service for config file [%s] stopped", cfgFile)
	}
	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if errRet != nil {
		err = errRet
		return
	}

	return
}
