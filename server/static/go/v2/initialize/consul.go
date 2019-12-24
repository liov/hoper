package initialize

import "github.com/spf13/viper"

func P1Consul() {
	var runtime_viper = viper.New()
	runtime_viper.WatchRemoteConfig()
}
