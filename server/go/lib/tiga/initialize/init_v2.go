package initialize

import (
	"github.com/actliboy/hoper/server/go/lib/utils/reflect/mtos"
	"github.com/actliboy/hoper/server/go/lib/utils/slices"
	"github.com/pelletier/go-toml"
	"reflect"
	"strings"
)

func (init *Init) UnmarshalAndSetV2(bytes []byte) {
	toml.Unmarshal(bytes, &init.confM)
	init.CloseDao()
	init.injectM()
}

// Customize
func (init *Init) injectM() {
	v := reflect.ValueOf(init.conf).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setConfigM(v.Field(i).Interface(), v.Type().Field(i).Name, init.confM)
		case reflect.Struct:
			setConfigM(v.Field(i).Addr().Interface(), v.Type().Field(i).Name, init.confM)
		}
	}
	init.conf.Init()
	if init.dao == nil {
		return
	}
	setDaoM(reflect.ValueOf(init.dao).Elem(), init.confM)
	init.dao.Init()
}

func setConfigM(field interface{}, fieldName string, confM map[string]interface{}) {
	for k, v := range confM {
		if strings.ToUpper(k) == strings.ToUpper(fieldName) {
			mtos.Decode(v, field)
		}
	}
}

func setDaoM(v reflect.Value, confM map[string]interface{}) {

	if !v.IsValid() {
		return
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Addr().CanInterface() {
			if field.Kind() == reflect.Ptr && !field.IsValid() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			inter := field.Addr().Interface()
			if slices.StringContains(InitConfig.NoInject, strings.ToUpper(typ.Field(i).Name)) {
				continue
			}
			if daofield, ok := inter.(DaoField); ok {
				conf := daofield.Config()
				var confName string
				if confName, ok = typ.Field(i).Tag.Lookup("config"); ok {
					confName = strings.ToUpper(confName)
				} else {
					confName = strings.ToUpper(typ.Field(i).Name)
				}
				setConfigM(conf, confName, confM)
				if conf1, ok := conf.(NeedInit); ok {
					conf1.Init()
				}
				if conf1, ok := conf.(Generate); ok {
					daofield.SetEntity(conf1.Generate())
				}
			}
		}
	}
}
