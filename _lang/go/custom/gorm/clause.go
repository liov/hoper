package main

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"
	"sync"
)

func main() {

	db, _ = gorm.Open(tests.DummyDialector{}, nil)
	user, _ := schema.Parse(&tests.User{}, &sync.Map{}, db.NamingStrategy)
	stmt := gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}
	clauses := []clause.Interface{clause.Select{}, clause.From{}, clause.Where{Exprs: []clause.Expression{clause.Eq{Column: clause.PrimaryColumn, Value: "1"}, clause.Gt{Column: "age", Value: 18}, clause.Or(clause.Neq{Column: "name", Value: "jinzhu"}, clause.Neq{Column: "id", Value: "1"})}}}

	for _, clause := range clauses {
		stmt.AddClause(clause)
	}

	stmt.Build("SELECT", "FROM", "WHERE")
	fmt.Println(stmt.SQL.String())
}
