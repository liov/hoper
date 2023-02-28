package plugin

import (
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"strconv"
	"strings"
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type Builder struct {
	plugin   *protogen.Plugin
	params   map[string]string
	pathType pathType
}

func Generate(gen *protogen.Plugin) error {
	return nil
}

func New(opts protogen.Options, request *pluginpb.CodeGeneratorRequest) (*Builder, error) {
	plugin, err := opts.New(request)
	if err != nil {
		return nil, err
	}

	b := &Builder{
		plugin: plugin,
		params: parseParameter(request.GetParameter()),
	}

	for k, v := range b.params {
		switch k {

		case "paths":
			switch v {
			case "import":
				b.pathType = pathTypeImport
			case "source_relative":
				b.pathType = pathTypeSourceRelative
			}
		}
	}

	return b, nil
}

func parseParameter(param string) map[string]string {
	paramMap := make(map[string]string)

	for _, p := range strings.Split(param, ",") {
		if i := strings.Index(p, "="); i < 0 {
			paramMap[p] = ""
		} else {
			paramMap[p[0:i]] = p[i+1:]
		}
	}

	return paramMap
}

func (b *Builder) Generate() (*pluginpb.CodeGeneratorResponse, error) {

	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range b.plugin.Files {
		if !protoFile.Generate {
			continue
		}
		fileName := protoFile.GeneratedFilenamePrefix + ".pb.enum.go"
		g := b.plugin.NewGeneratedFile(fileName, ".")
		genFileMap[fileName] = g
		if !FileEnabledExtGen(protoFile) || len(protoFile.Enums) == 0 {
			continue
		}
		// third traverse: build associations
		for _, enum := range protoFile.Enums {
			if EnabledExtGen(enum) {
				genFileMap[fileName] = g
				break
			}
		}

	}

	for _, protoFile := range b.plugin.Files {
		fileName := protoFile.GeneratedFilenamePrefix + ".pb.enum.go"
		g, ok := genFileMap[fileName]
		if !ok {
			continue
		}

		g.P("package ", protoFile.GoPackageName)

		for _, enum := range protoFile.Enums {
			b.generate(protoFile, enum, g)
		}

	}

	return b.plugin.Response(), nil
}

func (b *Builder) generate(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	if EnabledEnumStringer(e) {
		b.generateString(f, e, g)
	}
	if EnabledEnumJsonMarshal(e) {
		b.generateJsonMarshal(f, e, g)
	}
	if EnabledEnumErrorCode(e) {
		b.generateErrorCode(f, e, g)
	}
	if EnabledEnumGqlGen(e) {
		b.generateGQLMarshal(f, e, g)
	}
}

func (b *Builder) generateString(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	g.P("func (x ", ccTypeName, ") String() string {")
	g.P()
	if len(e.Values) > 64 {
		g.P("return ", ccTypeName, "_name[x]")
	} else {
		g.P("switch x {")
		for _, ev := range e.Values {
			name := stringsi.CamelCase(ev.Desc.Name())
			//PrintComments(e.Comments, g)
			value := name
			g.P("case ", name, " :")
			if cn := GetEnumValueCN(ev); cn != "" {
				value = cn
			}
			g.P("return ", strconv.Quote(value))
		}
	}
	g.P("}")
	g.P("return ", strconv.Quote(""))
	g.P("}")
	g.P()
}

func (b *Builder) generateGQLMarshal(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	typ := "uint32"
	if typ1 := GetEnumType(e); typ1 != "" {
		typ = typ1
	}
	g.P("func (x ", ccTypeName, ") MarshalGQL(w ", generateImport("Writer", "io", g), ") {")
	g.P(`w.Write(`, generateImport("QuoteToBytes", "github.com/liov/hoper/server/go/lib/utils/strings", g), `(x.String()))`)
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalGQL(v interface{}) error {")
	g.P(`if i, ok := v.(`, typ, "); ok {")
	g.P(`*x = `, ccTypeName, `(i)`)
	g.P("return nil")
	g.P("}")
	g.P(`return `, generateImport("New", "errors", g), `("枚举值需要数字类型")`)
	g.P("}")
	g.P()
}

func (b *Builder) generateJsonMarshal(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	/*ccTypeName := stringsi.CamelCase(e.Desc.Name())
	if EnabledGoEnumValueMap(e) {
		g.P("func (x ", ccTypeName, ") MarshalJSON() ([]byte, error) {")
		g.P("return ", generateImport("QuoteToBytes", "github.com/liov/hoper/server/go/lib/utils/strings", g), "(x.String())", ", nil")
		g.P("}")
		g.P()
	}

		g.P("func (x *", ccTypeName, ") UnmarshalJSON(data []byte) error {")

		g.P("value, ok := ", ccTypeName, `_name[string(data)]`)
		g.P("if ok {")

		g.P("*x = ", ccTypeName, "(value)")
		g.P("return nil")

		g.P("}")
		g.P(`return `, generateImport("New", "errors", g), `("无效的`, ccTypeName, `")`)

		g.P("}")
		g.P()
	}*/
}

func (b *Builder) generateErrorCode(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	g.P("func (x ", ccTypeName, ") Error() string {")

	g.P(`return x.String()`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrRep() *", generateImport("ErrRep", "github.com/liov/hoper/server/go/lib/protobuf/errorcode", g), " {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") Message(msg string) error {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: msg}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrorLog(err error) error {")

	g.P(generateImport("Error", "github.com/liov/hoper/server/go/lib/utils/log", g), `(err)`)
	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") GRPCStatus() *", generateImport("Status", "google.golang.org/grpc/status", g), " {")

	g.P(`return `, `status.New(`, generateImport("Code", "google.golang.org/grpc/codes", g), `(x), x.String())`)

	g.P("}")
	g.P()
}

func generateImport(name string, importPath string, g *protogen.GeneratedFile) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}
