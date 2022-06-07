package initialize

import (
	"strings"
)

const (
	exprTag    = "expr"
	configTag  = "config"
	ConfigName = "CONFIG"
	IsInject   = "NOTINJECT"
)

type TagSettings struct {
	NotInject  bool
	ConfigName string
}

var tags = []string{IsInject, ConfigName}

func (t *TagSettings) Set(index int, value string) {
	switch index {
	case 0:
		t.NotInject = value == "true"
	case 1:
		t.ConfigName = strings.ToUpper(value)
	}
}

func ParseTagSetting(str string, sep string) TagSettings {
	var settings TagSettings
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
		if settings.NotInject {
			return settings
		}
	}

	return settings
}
