package generator

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	cuzproto "github.com/liov/hoper/go/v2/protobuf/utils/proto/gogo"
)

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
		v, err := proto.GetExtension(field.Options, cuzproto.E_EnumvalueCn)
		if err == nil && v.(*string) != nil {
			return *(v.(*string))
		}
	}
	return ""
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
		v, err := proto.GetExtension(field.Options, cuzproto.E_EnumCustomtype)
		if err == nil && v.(*string) != nil {
			return *(v.(*string))
		}
	}
	return ""
}

func EnabledGoEnumValueMap(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, cuzproto.E_EnumGenvaluemap, false)
}

func TurnOffGoEnumValueMap(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(cuzproto.E_EnumGenvaluemap, false)(enum)
}

func EnabledEnumNumOrder(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, cuzproto.E_EnumNumorder, false)
}

func TurnOffEnumNumOrder(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(cuzproto.E_EnumNumorder, false)(enum)
}

func EnabledEnumJsonMarshal(file *google_protobuf.FileDescriptorProto, enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, cuzproto.E_EnumJsonmarshal, false)
}

func TurnOffEnumJsonMarshal(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(cuzproto.E_EnumJsonmarshal, false)(enum)
}

func EnabledEnumErrorCode(enum *google_protobuf.EnumDescriptorProto) bool {
	return proto.GetBoolExtension(enum.Options, cuzproto.E_EnumErrorcode, false)
}

func TurnOffEnumErrorCode(enum *descriptor.EnumDescriptorProto) {
	vanity.SetBoolEnumOption(cuzproto.E_EnumErrorcode, false)(enum)
}
