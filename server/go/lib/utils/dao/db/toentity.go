package dbi

import (
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const Tmpl = `package entity
import "time"

type Example struct{
A int` + "`json:\"a\" explain:\"模板\"`" + `
B string
C time.Time
}
`

const FileTmpl = `package generate

import "time"

`
const TagTmpl = "`json:\"%s\" explain:\"%s\"`"

func NewLine() byte {
	return '\n'
}

func TwoLine() []byte {
	return []byte("\n\n")
}

func AddStruct(name string) []byte {
	return []byte(`type ` + name + ` struct{`)
}

func StructEnd(name string) []byte {
	return []byte(`}`)
}

type Field struct {
	Field   string
	Type    string
	Comment string
}

func (f *Field) goTYpe(dbType string) string {
	if dbType == MYSQL {
		if strings.Contains(f.Type, "int") {
			return "int"
		}
		if strings.Contains(f.Type, "varchar") || strings.Contains(f.Type, "text") {
			return "string"
		}
		if strings.Contains(f.Type, "timestamp") || strings.Contains(f.Type, "datetime") || strings.Contains(f.Type, "date") {
			return "time.Time"
		}
		if strings.Contains(f.Type, "float") || strings.Contains(f.Type, "double") || strings.Contains(f.Type, "decimal") {
			return "float64"
		}
		return "bool"
	}
	return ""
}

func (f *Field) Generate(dbType string) *ast.Field {
	field := stringsi.ConvertToCamelCase(f.Field)
	return &ast.Field{
		Doc: nil,
		Names: []*ast.Ident{
			{
				Name: field,
				Obj:  &ast.Object{Kind: ast.Var, Name: f.Field},
			},
		},
		Type:    &ast.Ident{Name: f.goTYpe(dbType)},
		Tag:     &ast.BasicLit{Kind: token.STRING, Value: "`" + `json:"` + stringsi.LowerFirst(field) + `" explain:"` + f.Comment + "\"`"},
		Comment: nil,
	}
}

func GetDecl() *ast.GenDecl {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "tmpl.go", Tmpl, parser.ParseComments)
	decl := f.Decls[1].(*ast.GenDecl)
	decl.Rparen = token.Pos(10)
	return decl
}

/*func generate() {
	file := ast.File{
		Doc:     nil,
		Package: 0,
		Name: &ast.Ident{
			NamePos: 0,
			Name:    "",
			Obj:     nil,
		},
		Decls:      nil,
		Scope:      nil,
		Imports:    nil,
		Unresolved: nil,
		Comments:   nil,
	}
}
*/
