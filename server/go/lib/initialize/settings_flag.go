package initialize

import (
	"flag"
	"github.com/spf13/pflag"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// TODO: 优先级高于EnvConfig
type FlagConfig struct {
	PreSet
	Watch bool
	AfterSet
}

type PreSet struct {
	Env, ConfUrl string
}

type AfterSet struct {
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.StringVarP(&InitConfig.Env, "env", "e", DEVELOPMENT, "环境")

	if _, err := os.Stat(InitConfig.ConfUrl); os.IsNotExist(err) {
		InitConfig.ConfUrl = "./config/config.toml"
	}
	pflag.StringVarP(&InitConfig.ConfUrl, "conf", "c", InitConfig.ConfUrl, "配置文件路径,默认./config.toml或./config/config.toml")
	pflag.BoolVarP(&InitConfig.flag.Watch, "watch", "w", false, "是否监听配置文件")
	var agent string // socks5://localhost:1080
	pflag.StringVarP(&agent, "proxy", "p", "", "是否启用代理")

	pflag.Parse()
	if agent != "" {
		proxyURL, _ := url.Parse(agent)
		http.DefaultClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}
	}
}

func (init *initConfig) applyFlagConfig() {
	if !init.flag.Watch && init.ConfigCenterConfig != nil {
		init.ConfigCenterConfig.Watch = false
	}
}
