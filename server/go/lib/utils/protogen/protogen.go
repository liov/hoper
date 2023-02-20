package protogeni

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"log"
	"reflect"
)

// 这泛型真别扭
type BaseType interface {
}

func GetOption(desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def any) any {
	if desc == nil {
		log.Println("desc is nil")
		return def
	}
	if !proto.HasExtension(desc.Options(), xt) {
		return def
	}

	v := proto.GetExtension(desc.Options(), xt)

	rv := reflect.ValueOf(v)

	if rv.IsValid() {
		return v
	}

	return def
}

func GetBaseTypeOption(desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def any) any {
	if desc == nil {
		return def
	}
	return proto.GetExtension(desc.Options(), xt)
}

func GetStructTypeOption(desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def any) any {
	if desc == nil {
		return def
	}

	return proto.GetExtension(desc.Options(), xt)
}

func SetExtension(desc protoreflect.Descriptor, extension *protoimpl.ExtensionInfo, value any) {
	if !proto.HasExtension(desc.Options(), extension) {
		return
	}
	proto.SetExtension(desc.Options(), extension, value)
}

func GenerateImport(name string, importPath string, g *protogen.GeneratedFile) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}
