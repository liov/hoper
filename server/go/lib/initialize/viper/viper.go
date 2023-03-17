package viper

import (
	initialize2 "github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/log"
	reflecti "github.com/liov/hoper/server/go/lib/utils/reflect"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type Config struct {
	Remote     bool
	Watch      bool
	Provider   string
	Endpoint   string
	Path       string
	ConfigType string
}

func (conf *Config) Init() {

}

func (conf *Config) Build() *viper.Viper {
	iconf := initialize2.InitConfig.Config()
	if exist := reflecti.GetFieldValue(iconf, conf); !exist {
		return nil
	}
	var runtimeViper = viper.GetViper()

	if conf.ConfigType == "" {
		conf.ConfigType = "toml"
	}

	runtimeViper.SetConfigType(conf.ConfigType) // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "Env", "dotenv"
	if conf.Remote {
		runtimeViper.AddRemoteProvider(conf.Provider, conf.Endpoint, initialize2.InitKey)
		// read from remote Config the first time.
		err := runtimeViper.ReadRemoteConfig()
		if err != nil {
			log.Error(err)
		}
		if conf.Watch {
			runtimeViper.WatchRemoteConfig()
		}

	} else {
		runtimeViper.AddConfigPath(conf.Path)
		runtimeViper.ReadInConfig()
		if conf.Watch {
			runtimeViper.WatchConfig()
		}
	}

	// unmarshal Config
	cCopy := iconf
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
			vconf :=runtime_viper.AllSettings()
			log.Debug(vconf)
			// unmarshal new Config into our runtime Config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(cCopy)
			refresh(cCopy, dCopy)
			log.Debug(cCopy)
		}
	}()*/
	return runtimeViper
}

type Viper struct {
	*viper.Viper
	Conf Config
}

func (v *Viper) Config() any {
	return &v.Conf
}

func (v *Viper) SetEntity() {
	v.Viper = v.Conf.Build()
}

func (v *Viper) Close() error {
	return nil
}
