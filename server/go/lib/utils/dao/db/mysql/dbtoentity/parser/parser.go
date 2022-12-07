package parser

import (
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db"
	"github.com/liov/hoper/server/go/lib/utils/dao/db/mysql"
	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

func MysqlConvertByTable(sql string) {
	gen := NewMysqlGen(sql)
	dbi.ConvertByTable(gen, gen.Tables()[0])
}

type mysqlgen struct {
	*sqlparser.CreateTable
}

func NewMysqlGen(sql string) *mysqlgen {
	statement, _ := sqlparser.Parse(sql)
	return &mysqlgen{statement.(*sqlparser.CreateTable)}
}

func (m *mysqlgen) Tables() []string {
	return []string{m.DDL.NewName.Name.String()}
}

func (m *mysqlgen) Fields(tableName string) []*dbi.Field {
	var dbfields []*dbi.Field
	for i, column := range m.Columns {
		dbfields = append(dbfields, &dbi.Field{
			Field:   column.Name,
			Type:    column.Type,
			Comment: "",
			GoTYpe:  "",
		})
		for _, option := range column.Options {
			if option.Type == sqlparser.ColumnOptionComment {
				dbfields[i].Comment = option.Value[1 : len(option.Value)-1]
				break
			}
		}
	}
	return dbfields
}

func (m *mysqlgen) TypeToGoTYpe(typ string) string {
	return mysql.MySqlTypeToGoType(typ)
}
