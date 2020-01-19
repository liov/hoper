package generator

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func init() {
	proto.RegisterExtension(E_EnumvalueCN)
	proto.RegisterExtension(E_GoprotoEnumType)
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

var E_GoprotoEnumType = &proto.ExtensionDesc{
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
		v, err := proto.GetExtension(field.Options, E_EnumvalueCN)
		if err == nil && v.(*string) != nil {
			return *(v.(*string))
		}
	}
	return ""
}
