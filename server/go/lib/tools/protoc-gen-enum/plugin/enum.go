package plugin

import (
	"github.com/liov/hoper/server/go/lib/protobuf/utils/enum"
	protogeni "github.com/liov/hoper/server/go/lib/utils/protogen"
	"google.golang.org/protobuf/compiler/protogen"
)

func FileEnabledExtGen(f *protogen.File) bool {
	return protogeni.GetOption(f.Desc, enum.E_EnumExtGenAll, true).(bool)
}

func EnabledExtGen(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumExtGen, true).(bool)
}

func GetEnumValueCN(ev *protogen.EnumValue) string {
	return protogeni.GetOption(ev.Desc, enum.E_EnumvalueCn, "").(string)
}

func GetEnumType(e *protogen.Enum) string {

	return protogeni.GetOption(e.Desc, enum.E_EnumCustomtype, "int32").(string)
}

func EnabledGoEnumValueMap(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumGenvaluemap, true).(bool)
}

func TurnOffGoEnumValueMap(e *protogen.Enum) {
	protogeni.SetExtension(e.Desc, enum.E_EnumGenvaluemap, false)
}

func EnabledEnumNumOrder(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumNumorder, false).(bool)
}

func TurnOffEnumNumOrder(e *protogen.Enum) {
	protogeni.SetExtension(e.Desc, enum.E_EnumNumorder, false)
}

func EnabledEnumJsonMarshal(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumJsonmarshal, true).(bool)
}

func TurnOffEnumJsonMarshal(e *protogen.Enum) {
	protogeni.SetExtension(e.Desc, enum.E_EnumJsonmarshal, false)
}

func EnabledEnumErrorCode(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumErrorcode, false).(bool)
}

func TurnOffEnumErrorCode(e *protogen.Enum) {
	protogeni.SetExtension(e.Desc, enum.E_EnumErrorcode, false)
}

func EnabledEnumGqlGen(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumGqlgen, true).(bool)
}

func TurnOffEnumGqlGen(e *protogen.Enum) {
	protogeni.SetExtension(e.Desc, enum.E_EnumGqlgen, false)
}

func EnabledFileEnumGqlGen(f *protogen.File) bool {
	return protogeni.GetOption(f.Desc, enum.E_EnumGqlgenAll, false).(bool)
}

func TurnOffFileEnumGqlGen(f *protogen.File) {
	protogeni.SetExtension(f.Desc, enum.E_EnumGqlgenAll, false)
}

func EnabledGoEnumPrefix(f *protogen.File, e *protogen.Enum) bool {
	if protogeni.GetOption(f.Desc, enum.E_EnumPrefixAll, false).(bool) {
		return true
	}
	return protogeni.GetOption(e.Desc, enum.E_EnumPrefix, true).(bool)
}

func EnabledEnumStringer(e *protogen.Enum) bool {
	return protogeni.GetOption(e.Desc, enum.E_EnumStringer, true).(bool)
}
