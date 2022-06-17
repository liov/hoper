package main

import (
	"flag"
	http_fs "github.com/actliboy/hoper/server/go/lib/utils/net/http/fs"
	http_fs_watch "github.com/actliboy/hoper/server/go/lib/utils/net/http/fs/watch"
	_ "github.com/fatedier/frp/assets/frpc"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/util/log"
	"os"
	"time"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "https://xxx", "url")
	flag.Parse()
	log.Info("url: %s", url)
	urlchan := make(chan string)
	watch := http_fs_watch.New(time.Second)
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
		if cfg.AssetsDir != "" && cfg.AssetsDir != url {
			urlchan <- cfg.AssetsDir
			cfg.AssetsDir = ""
		}
	}
	watch.Add(url, callback)
	go func() {
		for newurl := range urlchan {
			log.Info("url: %s", newurl)
			watch.Add(newurl, callback)
			watch.Remove(url)
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
