package generator

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
)

func init() {
	proto.RegisterExtension(E_EnumvalueCN)
	proto.RegisterExtension(E_EnumType)
	proto.RegisterExtension(E_GenEnumValueMap)
	proto.RegisterExtension(E_EnumNumOrder)
}

var E_EnumvalueCN = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumValueOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         66002,
	Name:          "enumproto.enumvalue_cn",
	Tag:           "bytes,66002,opt,name=enumvalue_cn",
	Filename:      "enum.proto",
}

func IsEnumValueCN(field *google_protobuf.EnumValueDescriptorProto) bool {
	name := GetEnumValueCN(field)
	if len(name) > 0 {
		return true
	}
	return false
}

func GetEnumValueCN(field *google_protobuf.EnumValueDescriptorProto) string {
	if field == nil {
		return ""
	}
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, E_EnumvalueCN)
		if err == nil && v.(*string) != nil {
			return *(v.(*string))
		}
	}
	return ""
}

var E_EnumType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         62025,
	Name:          "enumproto.enum_customtype",
	Tag:           "bytes,62025,opt,name=enum_customtype",
	Filename:      "enum.proto",
}

func IsEnumType(field *google_protobuf.EnumDescriptorProto) bool {
	name := GetEnumType(field)
	if len(name) > 0 {
		return true
	}
	return false
}

func GetEnumType(field *google_protobuf.EnumDescriptorProto) string {
	if field == nil {
		return ""
	}
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, E_EnumType)
		if err == nil && v.(*string) != nil {
			return *(v.(*string))
		}
	}
	return ""
}

var E_GenEnumValueMap = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62026,
	Name:          "enumproto.enum_genvaluemap",
	Tag:           "varint,62026,opt,name=enum_genvaluemap",
	Filename:      "enum.proto",
}

func EnabledGoEnumValueMap(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, E_GenEnumValueMap, false)
}

func TurnOffGoEnumValueMap(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(E_GenEnumValueMap, false)(enum)
}

var E_EnumNumOrder = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62027,
	Name:          "enumproto.enum_numorder",
	Tag:           "varint,62027,opt,name=enum_numorder",
	Filename:      "enum.proto",
}

func EnabledEnumNumOrder(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, E_EnumNumOrder, false)
}

func TurnOffEnumNumOrder(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(E_EnumNumOrder, false)(enum)
}
