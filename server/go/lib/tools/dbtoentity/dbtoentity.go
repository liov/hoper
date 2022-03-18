package dbtoentity

import (
	"bytes"
	"fmt"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

const tmpl = `package entity
import "time"

type Example struct{
A int` + "`json:\"a\" explain:\"模板\"`" + `
B string
C time.Time
}
`

const fileTmpl = `package generate

import "time"

`
const tagTmpl = "`json:\"%s\" explain:\"%s\"`"

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

func MysqlConvert(db *gorm.DB) {
	var tables []string
	db.Raw(`SHOW TABLES`).Scan(&tables)
	mysqlgen := mysqlgen{db: db, decl: GetDecl()}
	var buf bytes.Buffer
	buf.WriteString(fileTmpl)
	for i := range tables {
		buf.Write(mysqlgen.genTable(tables[i]))
		buf.Write(TwoLine())
	}
	write(&buf, "generate.go")
}

func MysqlConvertByTable(db *gorm.DB, tableName string) {
	var tables []string
	db.Raw(`SHOW TABLES`).Scan(&tables)
	mysqlgen := mysqlgen{db: db, decl: GetDecl()}
	var buf bytes.Buffer
	buf.WriteString(fileTmpl)
	buf.Write(mysqlgen.genTable(tableName))
	buf.Write(TwoLine())
	write(&buf, tableName+".go")
}

func (f *Field) goTYpe() string {
	if strings.Contains(f.Type, "int") {
		return "int"
	}
	if strings.Contains(f.Type, "varchar") || strings.Contains(f.Type, "text") {
		return "string"
	}
	if strings.Contains(f.Type, "timestamp") || strings.Contains(f.Type, "datetime") {
		return "time.Time"
	}
	if strings.Contains(f.Type, "float") || strings.Contains(f.Type, "double") || strings.Contains(f.Type, "decimal") {
		return "float64"
	}
	return "bool"
}

func (f *Field) generate() *ast.Field {
	field := stringsi.ConvertToCamelCase(f.Field)
	return &ast.Field{
		Doc: nil,
		Names: []*ast.Ident{
			{
				Name: field,
				Obj:  &ast.Object{Kind: ast.Var, Name: f.Field},
			},
		},
		Type:    &ast.Ident{Name: f.goTYpe()},
		Tag:     &ast.BasicLit{Kind: token.STRING, Value: "`" + `json:"` + stringsi.LowerFirst(field) + `" explain:"` + f.Comment + "\"`"},
		Comment: nil,
	}
}

func GetDecl() *ast.GenDecl {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "tmpl.go", tmpl, parser.ParseComments)
	decl := f.Decls[1].(*ast.GenDecl)
	decl.Rparen = token.Pos(10)
	return decl
}

type mysqlgen struct {
	db   *gorm.DB
	decl *ast.GenDecl
}

func (m *mysqlgen) genTable(tableName string) []byte {
	node := m.decl.Specs[0].(*ast.TypeSpec)
	node.Name.Name = stringsi.ConvertToCamelCase(tableName)
	fields := node.Type.(*ast.StructType).Fields
	fields.List = nil
	var dbfields []*Field
	m.db.Raw(`SHOW FULL COLUMNS FROM ` + tableName).Scan(&dbfields)
	for j := range dbfields {
		fields.List = append(fields.List, dbfields[j].generate())
	}
	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), m.decl)
	if err != nil {
		fmt.Println(err)
	}
	return b.Bytes()
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

func write(buf *bytes.Buffer, filename string) {
	pwd, _ := os.Getwd()
	entity := filepath.Join(pwd, "entity")
	_, err := os.Stat(entity)
	if os.IsNotExist(err) {
		os.Mkdir(entity, 0666)
	}
	f, _ := os.Create(filepath.Join(entity, filename))
	defer f.Close()
	f.Write(buf.Bytes())
}
