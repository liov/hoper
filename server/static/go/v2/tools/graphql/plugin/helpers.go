package plugin

import (
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	cuzproto "github.com/liov/hoper/go/v2/protobuf/utils/proto/gogo"

	"strings"
	"unicode"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

func getEnum(file *descriptor.FileDescriptorProto, typeName string) *descriptor.EnumDescriptorProto {
	for _, enum := range file.GetEnumType() {
		if enum.GetName() == typeName {
			return enum
		}
	}

	for _, msg := range file.GetMessageType() {
		if nes := getNestedEnum(file, msg, strings.TrimPrefix(typeName, msg.GetName()+".")); nes != nil {
			return nes
		}
	}
	return nil
}

func getNestedEnum(file *descriptor.FileDescriptorProto, msg *descriptor.DescriptorProto, typeName string) *descriptor.EnumDescriptorProto {
	for _, enum := range msg.GetEnumType() {
		if enum.GetName() == typeName {
			return enum
		}
	}

	for _, nes := range msg.GetNestedType() {
		if res := getNestedEnum(file, nes, strings.TrimPrefix(typeName, nes.GetName()+".")); res != nil {
			return res
		}
	}
	return nil
}

func resolveRequired(field *descriptor.FieldDescriptorProto) bool {
	if v := GetGqlFieldOptions(field); v != nil {
		return v.GetRequired()
	}
	return false
}

func ToLowerFirst(s string) string {
	if len(s) > 0 {
		return string(unicode.ToLower(rune(s[0]))) + s[1:]
	}
	return ""
}

func GetGqlFieldOptions(field *descriptor.FieldDescriptorProto) *cuzproto.Field {
	if field.Options != nil {
		v, err := proto.GetExtension(field.Options, cuzproto.E_Field)
		if err == nil && v.(*cuzproto.Field) != nil {
			return v.(*cuzproto.Field)
		}
	}
	return nil
}

// Match parsing algorithm from Generator.CommandLineParameters
func Params(gen *generator.Generator) map[string]string {
	params := make(map[string]string)

	for _, parameter := range strings.Split(gen.Request.GetParameter(), ",") {
		kvp := strings.SplitN(parameter, "=", 2)
		if len(kvp) != 2 {
			continue
		}

		params[kvp[0]] = kvp[1]
	}

	return params
}

var BasicType = []string{}

func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

func splitCPackageType(ctype string) (packageName string, typ string) {
	ss := strings.Split(ctype, ".")
	if len(ss) == 1 {
		return "", ctype
	}
	packageName = strings.Join(ss[0:len(ss)-1], ".")
	typeName := ss[len(ss)-1]
	importStr := strings.Map(badToUnderscore, packageName)
	typ = importStr + "." + typeName
	return packageName, typ
}

func getCastType(field *descriptor.FieldDescriptorProto) (packageName string, typ string, err error) {
	if field.Options != nil {
		var v interface{}
		v, err = proto.GetExtension(field.Options, gogoproto.E_Casttype)
		if err == nil && v.(*string) != nil {
			ctype := *(v.(*string))
			packageName, typ = splitCPackageType(ctype)
			return packageName, typ, nil
		}
	}
	return "", "", err
}
