package dbtoentity

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize"
	initmysql "github.com/liov/hoper/server/go/lib/initialize/gormdb/mysql"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

type dao struct {
	MysqlTest initmysql.DB
}

func (d *dao) Close() {
}

func (d dao) Init() {

}

type config struct{}

func (d config) Init() {

}

var Dao dao
var Conf config

func TestDBToEntity(t *testing.T) {
	defer initialize.Start(&Conf, &Dao)()
	MysqlConvert(Dao.MysqlTest.DB)
}

func TestTableToEntity(t *testing.T) {
	defer initialize.Start(&Conf, &Dao)()
	MysqlConvertByTable(Dao.MysqlTest.DB, "sku_competition")
}

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "tmpl.go", dbi.Tmpl, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	var b strings.Builder
	err = format.Node(&b, fset, f)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(b.String())
	b.Reset()
	ty := f.Decls[1].(*ast.GenDecl)
	node := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
	node.Name.Name = "A"
	fileds := node.Type.(*ast.StructType).Fields
	fileds.List = append(fileds.List, &ast.Field{
		Doc: nil,
		Names: []*ast.Ident{
			{
				Name: "D",
				Obj:  &ast.Object{Kind: ast.Var, Name: "D"},
			},
		},
		Type:    &ast.Ident{Name: "time.Time"},
		Tag:     &ast.BasicLit{Kind: token.STRING, Value: `json:"d"`},
		Comment: nil,
	})
	err = format.Node(&b, fset, ty)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(b.String())
}
