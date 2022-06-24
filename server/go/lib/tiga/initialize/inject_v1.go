package initialize

import (
	"github.com/actliboy/hoper/server/go/lib/utils/slices"
	"github.com/pelletier/go-toml"
	"reflect"
	"strings"
)

func (init *Init) UnmarshalAndSetV1(bytes []byte) {
	toml.Unmarshal(bytes, init.conf)
	init.CloseDao()
	init.inject()
}

// Customize
func (init *Init) inject() {
	var fieldNameDaoMap = make(map[string]interface{})
	setConfig(reflect.ValueOf(init.conf).Elem(), fieldNameDaoMap)
	init.conf.Init()
	if init.dao == nil {
		return
	}
	setDao(reflect.ValueOf(init.dao).Elem(), fieldNameDaoMap)
	init.dao.Init()
}

func setConfig(v reflect.Value, fieldNameDaoMap map[string]interface{}) {
	if !v.IsValid() {
		return
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setConfig(v.Field(i).Elem(), fieldNameDaoMap)
		case reflect.Struct:
			setConfig(v.Field(i), fieldNameDaoMap)
		}
		if v.Field(i).Addr().CanInterface() {
			inter := v.Field(i).Addr().Interface()
			if conf, ok := inter.(NeedInit); ok {
				conf.Init()
			}
			confName := strings.ToUpper(typ.Field(i).Name)
			if slices.StringContains(InitConfig.ConfigCenterConfig.NoInject, confName) {
				continue
			}
			if conf, ok := inter.(Generate); ok {
				ret := conf.Generate()
				fieldNameDaoMap[confName] = ret
			}
		}
	}
}

func setDao(v reflect.Value, fieldNameDaoMap map[string]interface{}) {

	if !v.IsValid() {
		return
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		confName := strings.ToUpper(typ.Field(i).Name)
		if slices.StringContains(InitConfig.ConfigCenterConfig.NoInject, confName) {
			continue
		}
		var dao interface{}
		var ok bool
		tagSettings := ParseTagSetting(typ.Field(i).Tag.Get(tag), ";")
		if tagSettings.NotInject {
			continue
		}
		if tagSettings.ConfigName != "" {
			dao, ok = fieldNameDaoMap[tagSettings.ConfigName]
		}
		if !ok {
			dao, ok = fieldNameDaoMap[confName]
		}
		if ok {
			daoValue := reflect.ValueOf(dao)
			if daoValue.Type().AssignableTo(v.Field(i).Type()) || daoValue.Type().Implements(v.Field(i).Type()) {
				v.Field(i).Set(daoValue)
			}
			continue
		}
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setDao(v.Field(i).Elem(), fieldNameDaoMap)
		case reflect.Struct:
			setDao(v.Field(i), fieldNameDaoMap)
		}
	}
}
