package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type EtcdConfig struct {
	Addr string
}

func (init *Init) P0EtcdOnce() {
	conf := &EtcdConfig{}
	if exist := reflecti.GetFieldValue(init.conf, conf); !exist {
		return
	}
	var runtimeViper = viper.New()

	runtimeViper.AddRemoteProvider("etcd", conf.Addr, InitKey)
	runtimeViper.SetConfigType("toml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "Env", "dotenv"

	// read from remote Config the first time.
	err := runtimeViper.ReadRemoteConfig()
	if err != nil {
		log.Error(err)

	}
	// unmarshal Config
	cCopy := init.conf
	//dCopy := init.dao
	runtimeViper.Unmarshal(cCopy)
	log.Debug(cCopy)
	// open a goroutine to watch remote changes forever
	//这段实现不够优雅
	/*	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Errorf("unable to read remote Config: %v", err)
				continue
			}
			conf :=runtime_viper.AllSettings()
			log.Debug(conf)
			// unmarshal new Config into our runtime Config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(cCopy)
			refresh(cCopy, dCopy)
			log.Debug(cCopy)
		}
	}()*/
}
