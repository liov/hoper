package initialize

import (
	"time"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ConsulConfig struct {
	Addr string
}

func (init *Init) P0ConsulOnce() {
	conf := &ConsulConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return
	}
	var runtime_viper = viper.New()

	runtime_viper.AddRemoteProvider("consul", conf.Addr, InitKey)
	runtime_viper.SetConfigType("toml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		log.Error(err)

	}
	// unmarshal config
	cCopy := init.conf
	dCopy := init.dao
	runtime_viper.Unmarshal(cCopy)
	log.Debug(cCopy)
	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Errorf("unable to read remote config: %v", err)
				continue
			}

			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(cCopy)
			Refresh(cCopy, dCopy)
			log.Debug(cCopy)
		}
	}()
}
