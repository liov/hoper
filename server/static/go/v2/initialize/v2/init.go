package initialize

import (
	"flag"
	"reflect"
	"strings"
	"unsafe"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/tools/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
)

type EnvConfig struct {
	Tenant string
}

type BasicNacosConfig struct {
	NacosAddr   string
	NacosGroup  string
	NacosDataId string
}

type InitConfig struct {
	Dev, Test, Prod *EnvConfig
	BasicNacosConfig
	initialize.Init
}

type Init struct {
	envConfig   *EnvConfig
	nacosConfig *nacos.Config
	initialize.Init
}

func NewInitWithLoadConfig(conf initialize.Config, dao initialize.Dao) *Init {
	initConfig := InitConfig{}
	log.Info(int64(uintptr(unsafe.Pointer(&initConfig))))
	err := configor.New(&configor.Config{Debug: true}).
		Load(&initConfig, initialize.ConfUrl) //"./Config/Config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	init := NewInit(conf, dao)

	init.BasicConfig = initConfig.BasicConfig
	init.NoInit = initConfig.NoInit

	value := reflect.ValueOf(&initConfig).Elem()
	typ := reflect.TypeOf(&initConfig).Elem()
	var tmpConfig = EnvConfig{}
	tmpTyp := reflect.TypeOf(&tmpConfig)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*			tmpConfig = value.Field(i).Interface().(*nacos.Config)
						//真·深度复制
						data,_:=json.Marshal(tmpConfig)
						if err:=json.Unmarshal(data,init.envConfig);err!=nil{
							log.Fatal(err)
						}*/
			//会被回收,也可能是被移动了？
			tmpConfig = *value.Field(i).Interface().(*EnvConfig)
			init.envConfig = &tmpConfig
			break
		}
	}
	init.nacosConfig = &nacos.Config{
		Addr:   initConfig.NacosAddr,
		Tenant: init.envConfig.Tenant,
		Group:  initConfig.NacosGroup,
		DataId: initConfig.NacosDataId,
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
		if dao != nil {
			dao.Close()
		}
		init.SetDao()
	})

	return func() {
		if dao != nil {
			dao.Close()
		}
		log.Sync()
	}
}

func (init *Init) config() *nacos.Client {
	nacosClient := init.nacosConfig.NewClient()
	nacosClient.GetConfigAllInfoHandle(init.Unmarshal)
	init.SetDao()
	return nacosClient
}
