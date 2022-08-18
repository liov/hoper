package initialize

import (
	"github.com/actliboy/hoper/server/go/lib/utils/slices"
	"github.com/pelletier/go-toml"
	"log"

	"reflect"
	"strings"
)

/*
注入这种

	type DB struct {
		*gorm.DB `init:"entity"`
		Conf     DatabaseConfig `init:"config"`
	}

	type config struct {
		//自定义的配置
		Customize serverConfig
		Server    server.ServerConfig
		Log       log.LogConfig
		Viper     *viper.Viper
	}

	type dao struct {
		// GORMDB 数据库连接
		GORMDB   db.DB
	}
*/
type daoField struct {
	Entity reflect.Value
	Config reflect.Value
}

const (
	EntityField = "ENTITY"
	ConfigField = "CONFIG"
)

func (init *Init) UnmarshalAndSetV3(bytes []byte) {
	tmp := map[string]interface{}{}
	err := toml.Unmarshal(bytes, &tmp)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range tmp {
		init.confM[strings.ToUpper(k)] = v
	}
	init.closeDao()
	init.inject3()
}

// Customize
func (init *Init) inject3() {
	setConfig2(init.conf, init.confM)
	init.conf.Init()
	if init.dao == nil {
		return
	}
	setDao3(reflect.ValueOf(init.dao).Elem(), init.confM)
	init.dao.Init()
}

func setDao3(v reflect.Value, confM map[string]any) {

	if !v.IsValid() {
		return
	}
	typ := v.Type()
	generateTyp := reflect.TypeOf((*Generate)(nil)).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Addr().CanInterface() {
			if field.Kind() == reflect.Ptr && (!field.IsValid() || field.IsNil()) {
				field.Set(reflect.New(field.Type().Elem()))
			}
			fieldtyp := field.Type()
			if fieldtyp == reflect.TypeOf(NeedInitPlaceholder{}) {
				continue
			}
			confName := strings.ToUpper(typ.Field(i).Name)
			if slices.StringContains(InitConfig.ConfigCenterConfig.NoInject, confName) {
				continue
			}

			var daoField daoField
			for j := 0; j < fieldtyp.NumField(); j++ {
				subfield := fieldtyp.Field(j)
				if strings.ToUpper(subfield.Name) == EntityField || strings.ToUpper(subfield.Tag.Get(tag)) == EntityField {
					daoField.Entity = field.Field(j)
					continue
				}
				if strings.ToUpper(subfield.Name) == ConfigField || strings.ToUpper(subfield.Tag.Get(tag)) == ConfigField || subfield.Type.Implements(generateTyp) {
					daoField.Config = field.Field(j)
				}
			}

			if daoField.Config.IsValid() {
				tagSettings := ParseDaoTagSettings(typ.Field(i).Tag.Get(tag))
				if tagSettings.NotInject {
					continue
				}
				if tagSettings.ConfigName != "" {
					confName = tagSettings.ConfigName
				}
				conf := daoField.Config.Addr().Interface()
				/*
					如果conf设置的是指针，且没有初始化，会有问题，这里初始化会报不可寻址，似乎不能返回interface{}
					valueConf := reflect.ValueOf(conf)
					if valueConf.Kind() == reflect.Ptr && (!valueConf.IsValid() || valueConf.IsNil()) {
						valueConf.Set(reflect.New(valueConf.Type().Elem()))
					}*/
				injectconf(conf, confName, confM)
				if conf1, ok := conf.(NeedInit); ok {
					conf1.Init()
				}
				if conf1, ok := conf.(Generate); ok {
					daoField.Entity.Set(reflect.ValueOf(conf1.Generate()))
				}
			}
		}
	}
}
