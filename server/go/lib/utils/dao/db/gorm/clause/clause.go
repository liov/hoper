package clausei

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"gorm.io/gorm/clause"
)

func Page(pageNo, pageSize int) clause.Limit {
	if pageSize == 0 {
		pageSize = 100
	}
	if pageNo > 1 {
		return clause.Limit{Offset: (pageNo - 1) * pageSize, Limit: pageSize}
	}
	return clause.Limit{Limit: pageSize}
}

func NewWhereClause(field string, op dbi.Operation, args ...interface{}) clause.Expression {
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

func DateBetween(column, dateStart, dateEnd string) clause.Expression {
	return NewWhereClause(column, dbi.Between, dateStart, dateEnd)
}

func Sort(column string, typ request.SortType) clause.Expression {
	var desc bool
	if typ == request.SortTypeDESC {
		desc = true
	}
	return clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Name: column, Raw: true}, Desc: desc}}}
}

type ListReq request.ListReq

func (req *ListReq) Clause() []clause.Expression {
	return []clause.Expression{Sort(req.SortField, req.SortType), Page(req.PageNo, req.PageSize)}
}

type RangeReq request.RangeReq

func (req *RangeReq) Clause() clause.Expression {
	operation := dbi.Between
	if req.RangeEnd == nil && req.RangeStart != nil {
		operation = dbi.Greater
		if req.Include {
			operation = dbi.GreaterOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart == nil && req.RangeEnd != nil {
		operation = dbi.Less
		if req.Include {
			operation = dbi.LessOrEqual
		}
		return NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart != nil && req.RangeEnd != nil {
		if req.Include {
			return NewWhereClause(req.RangeField, operation, req.RangeStart, req.RangeEnd)
		} else {
			return clause.Where{Exprs: []clause.Expression{NewWhereClause(req.RangeField, dbi.Greater, req.RangeStart), NewWhereClause(req.RangeField, dbi.Less, req.RangeStart)}}
		}
	}
	return nil
}
