package initialize

import (
	"github.com/spf13/pflag"
	"os"
	"reflect"
	"strings"
)

const (
	tag     = "init"
	exprTag = "expr"
)

type DaoTagSettings struct {
	NotInject  bool
	ConfigName string
}

type ConfigTagSettings struct {
	Flag string
	Env  string
}

func (c *ConfigTagSettings) Set(index int, value string) {
	switch index {
	case 0:
		c.Flag = strings.ToUpper(value)
	case 1:
		c.Env = strings.ToUpper(value)
	}
}

type TagSettings interface {
	Set(index int, value string)
}

var daotags = []string{"NOTINJECT", "CONFIG"}
var conftags = []string{"FLAG", "ENV"}

func (t *DaoTagSettings) Set(index int, value string) {
	switch index {
	case 0:
		t.NotInject = true
	case 1:
		t.ConfigName = strings.ToUpper(value)
	}
}

func ParseDaoTagSettings(str string) *DaoTagSettings {
	var settings DaoTagSettings
	ParseTagSetting(str, ";", &settings, daotags)
	return &settings
}

func ParseConfigTagSettings(str string) *ConfigTagSettings {
	var settings ConfigTagSettings
	ParseTagSetting(str, ";", &settings, conftags)
	return &settings
}

// ParseTagSetting default sep ;
func ParseTagSetting(str string, sep string, settings TagSettings, tags []string) {

	names := strings.Split(str, sep)
	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))
		var v string
		if len(values) >= 2 {
			v = strings.Join(values[1:], ":")
		} else if k != "" {
			v = "true"
		}
		for idx, key := range tags {
			if key == k {
				settings.Set(idx, v)
				break
			}
		}
	}

	return
}

func Unmarshal(v reflect.Value) error {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Ptr:
			Unmarshal(v.Field(i).Elem())
		case reflect.Struct:
			Unmarshal(v.Field(i))
		}
		tag := typ.Field(i).Tag.Get(tag)
		if tag != "" {
			settings := ParseConfigTagSettings(tag)
			switch field.Kind() {
			case reflect.String:
				field.Set(reflect.ValueOf(os.Getenv(settings.Env)))
				pflag.StringVarP(field.Addr().Interface().(*string), "", settings.Flag, field.Interface().(string), "")
			case reflect.Int:
				pflag.StringVarP(field.Addr().Interface().(*string), "", settings.Flag, field.Interface().(string), "")
			}

		}
	}
	return nil
}
