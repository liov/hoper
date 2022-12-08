package initialize

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize/conf_center"
	ilocal "github.com/liov/hoper/server/go/lib/initialize/conf_center/local"
	"github.com/liov/hoper/server/go/lib/utils/configor/local"
	"github.com/liov/hoper/server/go/lib/utils/errors/multierr"
	"os"
	"reflect"
	"strings"

	"github.com/liov/hoper/server/go/lib/utils/log"
)

// 约定大于配置
var (
	// 此处不是真正的初始化，只是为了让配置中心能够读取到配置
	InitConfig = &initConfig{
		Env: DEVELOPMENT, ConfUrl: "./config.toml",
	}
)

type Env string

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
	InitKey     = "initialize"
)

type ConfigCenterConfig struct {
	NoInject []string
	conf_center.ConfigCenterConfig
}

type BasicConfig struct {
	Module string
}

type FileConfig struct {
	BasicConfig
	Dev, Test, Prod *ConfigCenterConfig
}

type initConfig struct {
	Env, ConfUrl string
	BasicConfig
	ConfigCenterConfig *ConfigCenterConfig
	confM              map[string]interface{}
	conf               NeedInit
	dao                Dao
	//closes     []interface{}
	deferf      []func()
	initialized bool
}

func Start(conf Config, dao Dao, notinit ...string) func(deferf ...func()) {
	//逃逸到堆上了
	init := NewInit(conf, dao)
	init.LoadConfig(notinit...)
	init.initialized = true
	return func(deferf ...func()) {
		for _, f := range deferf {
			f()
		}
		for _, f := range init.deferf {
			f()
		}

	}
}

func (init *initConfig) LoadConfig(notinit ...string) *initConfig {
	onceConfig := FileConfig{}
	if _, err := os.Stat(init.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置文件在同一目录下或在config目录下或指定配置文件")
	}
	err := local.Load(&onceConfig, init.ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	fmt.Printf("Load config from: %s\n", init.ConfUrl)

	init.BasicConfig = onceConfig.BasicConfig

	value := reflect.ValueOf(&onceConfig).Elem()
	typ := reflect.TypeOf(&onceConfig).Elem()

	tmpTyp := reflect.TypeOf(&ConfigCenterConfig{})
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*tmpConfig = value.Field(i).Interface().(*nacos.Config)
			//真·深度复制
			data,_:=json.Marshal(tmpConfig)
			if err:=json.Unmarshal(data,init.ConfigCenterConfig);err!=nil{
				log.Fatal(err)
			}*/

			if value.Field(i).IsValid() && !value.Field(i).IsNil() {
				//会被回收,也可能是被移动了？
				init.ConfigCenterConfig = &(*value.Field(i).Interface().(*ConfigCenterConfig))
			} else {
				// 单配置文件
				init.ConfigCenterConfig = &ConfigCenterConfig{
					ConfigCenterConfig: conf_center.ConfigCenterConfig{
						ConfigType: "local",
						Watch:      true,
						Local: &ilocal.Local{
							Config:     local.Config{},
							ConfigPath: init.ConfUrl,
							ReloadType: "fsnotify",
						},
					},
				}
			}
			break
		}
	}

	for i := range init.ConfigCenterConfig.NoInject {
		init.ConfigCenterConfig.NoInject[i] = strings.ToUpper(init.ConfigCenterConfig.NoInject[i])
	}
	for i := range notinit {
		init.ConfigCenterConfig.NoInject = append(init.ConfigCenterConfig.NoInject, strings.ToUpper(notinit[i]))
	}

	cfgcenter := init.ConfigCenterConfig.ConfigCenter(init.Module, init.Env != PRODUCT)

	err = cfgcenter.HandleConfig(init.UnmarshalAndSet)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}

	return init
}

func NewInit(conf Config, dao Dao) *initConfig {
	init := &initConfig{
		Env: InitConfig.Env, ConfUrl: InitConfig.ConfUrl,
		confM: map[string]interface{}{},
		conf:  conf, dao: dao,
		deferf: []func(){
			func() { closeDao(dao) },
			log.Sync,
		},
	}
	InitConfig = init
	return init
}

func (init *initConfig) SetInit(conf Config, dao Dao) {
	init.conf = conf
	init.dao = dao
}

func (init *initConfig) RegisterDeferFunc(deferf ...func()) {
	init.deferf = append(init.deferf, deferf...)
}

func (init *initConfig) Config() Config {
	return init.conf
}

func (init *initConfig) closeDao() {
	if !init.initialized {
		return
	}
	err := closeDao(init.dao)
	if err != nil {
		log.Error(err)
	}
}

func closeDao(dao any) error {
	if dao == nil {
		return nil
	}
	if err, ok := closeDaoHelper(dao); ok {
		return err
	}
	var err multierr.MultiError
	daoValue := reflect.ValueOf(dao).Elem()
	for i := 0; i < daoValue.NumField(); i++ {
		/*	closer := daoValue.Field(i).MethodByName("Close")
			if closer.IsValid() {
				closer.Call(nil)
			}*/
		fieldV := daoValue.Field(i)
		if fieldV.Type().Kind() == reflect.Struct {
			fieldV = daoValue.Field(i).Addr()
		}
		field := fieldV.Interface()
		if err1, ok := closeDaoHelper(field); ok && err1 != nil {
			err.Append(err1)
		}
	}
	if err.HasErrors() {
		return &err
	}
	return nil
}

func closeDaoHelper(dao any) (error, bool) {
	if dao == nil {
		return nil, true
	}
	if closer, ok := dao.(DaoFieldCloser); ok {
		return closer.Close(), true
	}
	if closer, ok := dao.(DaoFieldCloser1); ok {
		closer.Close()
		return nil, true
	}
	return nil, false
}
