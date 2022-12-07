package initialize

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/reflect/mtos"
	"github.com/liov/hoper/server/go/lib/utils/slices"
	"github.com/pelletier/go-toml"
	"reflect"
	"strings"
)

/*
注入这种类型的dao
type DB struct {
	*gorm.DB
	Conf     DatabaseConfig
}

func (db *DB) Config() initialize.Generate {
	return &db.Conf
}

func (db *DB) SetEntity(entity interface{}) {
	if gormdb, ok := entity.(*gorm.DB); ok {
		db.DB = gormdb
	}
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

func (init *initConfig) UnmarshalAndSet(bytes []byte) {
	tmp := map[string]interface{}{}
	err := toml.Unmarshal(bytes, &tmp)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range tmp {
		init.confM[strings.ToUpper(k)] = v
	}
	init.closeDao()
	init.inject()
}

// Customize
func (init *initConfig) inject() {
	if init.conf != nil {
		setConfig(reflect.ValueOf(init.conf).Elem(), init.confM)
		init.conf.Init()
	}

	if init.dao != nil {
		setDao(reflect.ValueOf(init.dao).Elem(), init.confM)
		init.dao.Init()
	}
}

func setConfig(v reflect.Value, confM map[string]interface{}) {
	for i := 0; i < v.NumField(); i++ {
		filed := v.Field(i)
		switch filed.Kind() {
		case reflect.Ptr:
			injectconf(filed.Interface(), strings.ToUpper(v.Type().Field(i).Name), confM)
		case reflect.Struct:
			injectconf(filed.Addr().Interface(), strings.ToUpper(v.Type().Field(i).Name), confM)
		}
		if filed.Addr().CanInterface() {
			inter := v.Field(i).Addr().Interface()
			if conf, ok := inter.(NeedInit); ok {
				conf.Init()
			}
		}
	}
}

func injectconf(conf any, confName string, confM map[string]any) bool {
	filedv, ok := confM[confName]
	if ok {
		config := &mtos.DecoderConfig{
			Metadata: nil,
			Squash:   true,
			Result:   conf,
		}

		decoder, err := mtos.NewDecoder(config)
		if err != nil {
			return false
		}
		decoder.Decode(filedv)
	}
	return ok
}

type daoField struct {
	Entity               reflect.Value
	Config               reflect.Value
	entitySet, configSet bool
}

const (
	EntityField = "ENTITY"
	ConfigField = "CONFIG"
)

func setDao(v reflect.Value, confM map[string]any) {

	if !v.IsValid() {
		return
	}
	typ := v.Type()
	generateTyp := reflect.TypeOf((*Generate)(nil)).Elem()
	needInitPlaceholderTyp := reflect.TypeOf(NeedInitPlaceholder{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Addr().CanInterface() {
			fieldtyp := field.Type()
			if field.Kind() != reflect.Struct {
				log.Info(field.String(), "指针类型，忽略注入")
				continue
			}
			if fieldtyp == needInitPlaceholderTyp {
				continue
			}

			/*			if field.Kind() == reflect.Ptr && (!field.IsValid() || field.IsNil()) {
						field.Set(reflect.New(field.Type().Elem()))
					}*/

			confName := strings.ToUpper(typ.Field(i).Name)
			if slices.StringContains(InitConfig.ConfigCenterConfig.NoInject, confName) {
				continue
			}
			// 根据标签获取配置和要注入的类型
			var daoField daoField
			for j := 0; j < fieldtyp.NumField(); j++ {
				subfield := fieldtyp.Field(j)
				if strings.ToUpper(subfield.Name) == EntityField || strings.ToUpper(subfield.Tag.Get(tag)) == EntityField {
					daoField.Entity = field.Field(j)
					daoField.entitySet = true
					continue
				}
				// TODO: 定义新的Generate方法, 判断config Generate方法返回值方法类型为entity的类型
				if strings.ToUpper(subfield.Name) == ConfigField || strings.ToUpper(subfield.Tag.Get(tag)) == ConfigField {
					if subfield.Type.Implements(generateTyp) {
						daoField.Config = field.Field(j)
						daoField.configSet = true
					} else if field.Field(j).Addr().Type().Implements(generateTyp) {
						daoField.Config = field.Field(j).Addr()
						daoField.configSet = true
					}
				}
				if daoField.entitySet && daoField.configSet {
					break
				}
			}

			if daoField.entitySet && daoField.configSet && daoField.Config.IsValid() {
				// TODO:
				// daoField.Config.MethodByName("Build").Type().Out(0) == daoField.Entity.Type()
				tagSettings := ParseDaoTagSettings(typ.Field(i).Tag.Get(tag))
				if tagSettings.NotInject {
					continue
				}
				if tagSettings.ConfigName != "" {
					confName = tagSettings.ConfigName
				}
				conf := daoField.Config.Interface()
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
				continue
			}

			// 根据接口实现获取配置和要注入的类型
			inter := field.Addr().Interface()

			if daofield, ok := inter.(DaoField); ok {
				tagSettings := ParseDaoTagSettings(typ.Field(i).Tag.Get(tag))
				if tagSettings.NotInject {
					continue
				}
				if tagSettings.ConfigName != "" {
					confName = tagSettings.ConfigName
				}
				conf := daofield.Config()
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
				daofield.SetEntity(conf.Generate())
			}
		}
	}
}
