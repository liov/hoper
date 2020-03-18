package v2

import (
	"flag"
	"unsafe"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/tools/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

type EnvConfig struct {
	nacos.Config
	NoInit []string
}

type InitConfig struct {
	Dev, Test, Prod nacos.Config
	initialize.Init
}

type Init struct {
	envConfig nacos.Config
	initialize.Init
}

func NewInitWithLoadConfig(conf initialize.Config, dao initialize.Dao) *Init {
	initConfig := &InitConfig{}
	log.Info(int64(uintptr(unsafe.Pointer(initConfig))))
	err := configor.New(&configor.Config{Debug: true}).
		Load(initConfig, initialize.ConfUrl) //"./Config/Config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	init := NewInit(conf, dao)
	log.Info(int64(uintptr(unsafe.Pointer(init))))
	init.BasicConfig = initConfig.BasicConfig
	init.NoInit = initConfig.NoInit
	switch initConfig.Env {
	case initialize.DEVELOPMENT:
		init.envConfig = initConfig.Dev
	case initialize.TEST:
		init.envConfig = initConfig.Test
	case initialize.PRODUCT:
		init.envConfig = initConfig.Prod
	}
	return init
}

func NewInit(conf initialize.Config, dao initialize.Dao) *Init {
	return &Init{Init: *initialize.NewInit(conf, dao)}
}

func Start(conf initialize.Config, dao initialize.Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	//逃逸到堆上了
	init := NewInitWithLoadConfig(conf, dao)
	nacosClient := init.config()

	go nacosClient.Listener(func(bytes []byte) {
		init.Unmarshal(bytes)
		dao.Close()
		init.SetDao()
	})

	return func() {
		if dao != nil {
			dao.Close()
		}
		log.Sync()
		/*for _, f := range closes {
			res := reflect.ValueOf(f).Caller(nil)
			if len(res) > 0 && res[0].IsValid() {
				log.Error(res[0].Interface())
			}
		}*/
	}
}

func (init *Init) config() *nacos.Client {
	nacosClient := init.envConfig.NewClient()

	if nacosClient == nil {
		log.Fatal("配置错误")
	}
	nacosClient.GetConfigAllInfoHandle(init.Unmarshal)

	init.SetDao()
	return nacosClient
}
