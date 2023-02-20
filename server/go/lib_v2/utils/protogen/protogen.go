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

func GetOption[T any](desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def T) T {
	if desc == nil {
		log.Println("desc is nil")
		return def
	}
	if !proto.HasExtension(desc.Options(), xt) {
		return def
	}

	v, ok := proto.GetExtension(desc.Options(), xt).(T)

	if !ok {
		return def
	}
	rv := reflect.ValueOf(v)

	if rv.IsValid() {
		return v
	}
	return def
}

func GetBaseTypeOption[T any](desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def T) T {
	if desc == nil {
		return def
	}

	v, ok := proto.GetExtension(desc.Options(), xt).(T)
	if !ok {
		return def
	}
	return v
}

func GetStructTypeOption[T any](desc protoreflect.Descriptor, xt protoreflect.ExtensionType, def *T) *T {
	if desc == nil {
		return def
	}

	v, ok := proto.GetExtension(desc.Options(), xt).(*T)
	if !ok && v != nil {
		return def
	}
	return v
}

func SetExtension[T any](desc protoreflect.Descriptor, extension *protoimpl.ExtensionInfo, value T) {
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
