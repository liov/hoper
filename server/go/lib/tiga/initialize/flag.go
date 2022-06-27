package initialize

import (
	"flag"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

//TODO: 优先级高于EnvConfig
type FlagConfig struct {
	Env, ConfUrl string
}

func flaginit() {
	if flag.Parsed() {
		return
	}
	flag.StringVar(&InitConfig.Env, "env", DEVELOPMENT, "环境")

	if _, err := os.Stat(InitConfig.ConfUrl); os.IsNotExist(err) {
		InitConfig.ConfUrl = "./config/config.toml"
	}
	flag.StringVar(&InitConfig.ConfUrl, "conf", InitConfig.ConfUrl, "配置文件路径,默认./config.toml或./config/config.toml")

	agent := flag.Bool("agent", false, "是否启用代理")

	flag.Parse()
	if *agent {
		proxyURL, _ := url.Parse("socks5://localhost:1080")
		http.DefaultClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}
	}
}
