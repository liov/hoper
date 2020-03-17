package v2

import (
	"flag"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/tools/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Init struct {
	Dev  nacos.Config
	Test nacos.Config
	Prod nacos.Config
	initialize.Init
}

func NewInit(conf initialize.Config, dao initialize.Dao) *Init {
	init := &Init{Init: *initialize.NewInit(conf, dao)}
	err := configor.New(&configor.Config{Debug: false}).
		Load(init, initialize.ConfUrl) //"./Config/Config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	return init
}

func Start(conf initialize.Config, dao initialize.Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	init := NewInit(conf, dao)
	nacosClient := init.config()

	go nacosClient.Listener(func(bytes []byte) {
		init := NewInit(conf, dao)
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
	var nacosClient *nacos.Client
	switch init.Env {
	case initialize.DEVELOPMENT:
		nacosClient = init.Dev.NewClient()
	case initialize.TEST:
		nacosClient = init.Test.NewClient()
	case initialize.PRODUCT:
		nacosClient = init.Prod.NewClient()
	}
	if nacosClient == nil {
		log.Fatal("配置错误")
	}
	nacosClient.GetConfigAllInfoHandle(init.Unmarshal)
	init.SetDao()
	return nacosClient
}
