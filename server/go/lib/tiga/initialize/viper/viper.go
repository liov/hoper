package viper

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/reflect"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ViperConfig struct {
	Remote   bool
	Provider string
	Endpoint string
	Path     string
}

func (conf *ViperConfig) Init() {

}

func (conf *ViperConfig) Generate() interface{} {
	return conf.generate()
}

func (conf *ViperConfig) generate() *viper.Viper {
	iconf := initialize.InitConfig.Config()
	if exist := reflecti.GetFieldValue(iconf, conf); !exist {
		return nil
	}
	var runtimeViper = viper.GetViper()

	runtimeViper.SetConfigType("toml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "Env", "dotenv"
	if conf.Remote {
		runtimeViper.AddRemoteProvider(conf.Provider, conf.Endpoint, initialize.InitKey)
		// read from remote Config the first time.
		err := runtimeViper.ReadRemoteConfig()
		if err != nil {
			log.Error(err)
		}
		runtimeViper.WatchRemoteConfig()
	} else {
		runtimeViper.AddConfigPath(conf.Path)
		runtimeViper.ReadInConfig()
		runtimeViper.WatchConfig()
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
	Conf ViperConfig
}

func (v *Viper) Config() interface{} {
	return &v.Conf
}

func (v *Viper) SetEntity(entity interface{}) {
	if client, ok := entity.(*viper.Viper); ok {
		v.Viper = client
	}
}
