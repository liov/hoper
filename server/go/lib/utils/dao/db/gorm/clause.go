package gormi

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"gorm.io/gorm/clause"
)

func Page(pageNo, pageSize int) clause.Limit {
	if pageSize == 0 || pageSize > 100 {
		pageSize = 100
	}
	if pageNo != 0 {
		return clause.Limit{Offset: (pageNo - 1) * pageSize, Limit: pageSize}
	}
	return clause.Limit{Limit: pageSize}
}

func NewClause(field string, op dbi.Operation, args ...interface{}) clause.Expression {
	switch op {
	case dbi.Equal:
		return clause.Eq{
			Column: field,
			Value:  args[0],
		}
	case dbi.In:
		return clause.IN{
			Column: field,
			Values: args,
		}
	case dbi.Between:
		return clause.Expr{
			SQL:  field + " BETWEEN ? AND ?",
			Vars: args,
		}
	case dbi.Greater:
		return clause.Gt{
			Column: field,
			Value:  args[0],
		}
	case dbi.Less:
		return clause.Lt{
			Column: field,
			Value:  args[0],
		}
	case dbi.LIKE:
		return clause.Like{
			Column: field,
			Value:  args[0],
		}
	case dbi.GreaterOrEqual:
		return clause.Gte{
			Column: field,
			Value:  args[0],
		}
	case dbi.LessOrEqual:
		return clause.Lte{
			Column: field,
			Value:  args[0],
		}
	case dbi.NotIn:
		return clause.NotConditions{Exprs: []clause.Expression{clause.IN{
			Column: field,
			Values: args,
		}}}
	case dbi.NotEqual:
		return clause.Neq{
			Column: field,
			Value:  args[0],
		}
	case dbi.IsNull:
		return clause.Expr{
			SQL:  field + " IS NULL",
			Vars: nil,
		}
	case dbi.IsNotNull:
		return clause.Expr{
			SQL:  field + " IS NOT NULL",
			Vars: nil,
		}
	}
	return clause.Expr{
		SQL:  field,
		Vars: args,
	}
}
