package dbtoentity

import (
	"bytes"
	"fmt"
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go/ast"
	"go/format"
	"go/token"
	"gorm.io/gorm"
	"strings"
)

func MysqlConvert(db *gorm.DB, filename string) {
	var tables []string
	db.Raw(`SHOW TABLES`).Scan(&tables)
	mysqlgen := mysqlgen{db: db, decl: dbi.GetDecl()}
	var buf bytes.Buffer
	buf.WriteString(dbi.FileTmpl)
	for i := range tables {
		buf.Write(mysqlgen.genTable(tables[i]))
		buf.Write(dbi.TwoLine())
	}
	fs.Write(&buf, filename)
}

func MysqlConvertByTable(db *gorm.DB, tableName string) {
	mysqlgen := mysqlgen{db: db, decl: dbi.GetDecl()}
	var buf bytes.Buffer
	buf.WriteString(dbi.FileTmpl)
	buf.Write(mysqlgen.genTable(tableName))
	buf.Write(dbi.TwoLine())

	fs.Write(&buf, tableName+".go")
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
	var dbfields []*dbi.Field
	m.db.Raw(`SHOW FULL COLUMNS FROM ` + tableName).Scan(&dbfields)
	for j := range dbfields {
		dbfields[j].GoTYpe = MysqlTypeToGoTYpe(dbfields[j].Type)
		fields.List = append(fields.List, dbfields[j].Generate())
	}
	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), m.decl)
	if err != nil {
		fmt.Println(err)
	}
	return b.Bytes()
}

func MysqlTypeToGoTYpe(typ string) string {
	if strings.Contains(typ, "int") {
		return "int"
	}
	if strings.Contains(typ, "varchar") || strings.Contains(typ, "text") {
		return "string"
	}
	if strings.Contains(typ, "timestamp") || strings.Contains(typ, "datetime") || strings.Contains(typ, "date") {
		return "time.Time"
	}
	if strings.Contains(typ, "float") || strings.Contains(typ, "double") || strings.Contains(typ, "decimal") {
		return "float64"
	}
	return "bool"
}
