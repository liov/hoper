package dbi

import (
	"bytes"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
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
	GoTYpe  string
}

func (f *Field) Generate() *ast.Field {
	field := stringsi.ConvertToCamelCase(f.Field)
	return &ast.Field{
		Doc: nil,
		Names: []*ast.Ident{
			{
				Name: field,
				Obj:  &ast.Object{Kind: ast.Var, Name: f.Field},
			},
		},
		Type:    &ast.Ident{Name: f.GoTYpe},
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

type ConvertInterface interface {
	Tables() []string
	Fields(tableName string) []*Field
	TypeToGoTYpe(typ string) string
}

func Convert(c ConvertInterface, filename string) {
	tables := c.Tables()
	decl := GetDecl()
	var buf bytes.Buffer
	buf.WriteString(FileTmpl)
	for i := range tables {
		buf.Write(genTable(c, tables[i], decl))
		buf.Write(TwoLine())
	}
	fs.Write(&buf, filename)
}

func ConvertByTable(c ConvertInterface, tableName string) {
	decl := GetDecl()
	var buf bytes.Buffer
	buf.WriteString(FileTmpl)
	buf.Write(genTable(c, tableName, decl))
	buf.Write(TwoLine())

	fs.Write(&buf, tableName+".go")
}

func genTable(c ConvertInterface, tableName string, decl *ast.GenDecl) []byte {
	node := decl.Specs[0].(*ast.TypeSpec)
	node.Name.Name = stringsi.ConvertToCamelCase(tableName)
	fields := node.Type.(*ast.StructType).Fields
	fields.List = nil
	dbfields := c.Fields(tableName)
	for j := range dbfields {
		if dbfields[j].GoTYpe == "" {
			dbfields[j].GoTYpe = c.TypeToGoTYpe(dbfields[j].Type)
		}
		fields.List = append(fields.List, dbfields[j].Generate())
	}
	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), decl)
	if err != nil {
		fmt.Println(err)
	}
	return b.Bytes()
}
