package initialize

import (
	"flag"
	"os"
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
	NacosAddr  string
	NacosGroup string
	//NacosDataId string DataId与Moudle相同
}

type InitConfig struct {
	Dev, Test, Prod *EnvConfig
	BasicNacosConfig
	initialize.Init
}

var BasicConfig *Init

type Init struct {
	EnvConfig   *EnvConfig
	NacosConfig *nacos.Config
	initialize.Init
}

func NewInitWithLoadConfig(conf initialize.Config, dao initialize.Dao) *Init {
	initConfig := InitConfig{}
	log.Debug(int64(uintptr(unsafe.Pointer(&initConfig))))
	if _, err := os.Stat(initialize.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置目录在同一目录下")
	}
	err := configor.New(&configor.Config{Debug: true}).
		Load(&initConfig, initialize.ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	init := NewInit(conf, dao)

	init.BasicConfig = initConfig.BasicConfig
	init.Env = initialize.Env
	init.NoInit = initConfig.NoInit

	value := reflect.ValueOf(&initConfig).Elem()
	typ := reflect.TypeOf(&initConfig).Elem()
	var tmpConfig = EnvConfig{}
	tmpTyp := reflect.TypeOf(&tmpConfig)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*				tmpConfig = value.Field(i).Interface().(*nacos.Config)
							//真·深度复制
							data,_:=json.Marshal(tmpConfig)
							if err:=json.Unmarshal(data,init.EnvConfig);err!=nil{
								log.Fatal(err)
							}*/
			//会被回收,也可能是被移动了？
			tmpConfig = *value.Field(i).Interface().(*EnvConfig)
			init.EnvConfig = &tmpConfig
			break
		}
	}
	init.NacosConfig = &nacos.Config{
		Addr:   initConfig.NacosAddr,
		Tenant: init.EnvConfig.Tenant,
		Group:  initConfig.NacosGroup,
		DataId: initConfig.Module,
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
	BasicConfig = init
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
	nacosClient := init.NacosConfig.NewClient()
	nacosClient.GetConfigAllInfoHandle(init.Unmarshal)
	init.SetDao()
	return nacosClient
}
