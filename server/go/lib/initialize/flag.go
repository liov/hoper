package initialize

import (
	"github.com/spf13/pflag"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// TODO: 优先级高于EnvConfig
type FlagConfig struct {
	Env, ConfUrl string
}

func init() {
	pflag.StringVarP(&InitConfig.Env, "env", "e", DEVELOPMENT, "环境")

	if _, err := os.Stat(InitConfig.ConfUrl); os.IsNotExist(err) {
		InitConfig.ConfUrl = "./config/config.toml"
	}
	pflag.StringVarP(&InitConfig.ConfUrl, "conf", "c", InitConfig.ConfUrl, "配置文件路径,默认./config.toml或./config/config.toml")
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
