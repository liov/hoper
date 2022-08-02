package clausei

import (
	dbi "github.com/actliboy/hoper/server/go/lib/utils/dao/db"
	clause2 "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	_type "github.com/actliboy/hoper/server/go/lib/utils/generics/dao/db/type"
	"gorm.io/gorm/clause"
)

type RangeReq[T _type.Ordered] request.RangeReq[T]

func (req *RangeReq[T]) Clause() clause.Expression {
	if req == nil || req.RangeField == "" {
		return new(clause2.EmptyClause)
	}
	// 泛型还很不好用，这种方式代替原来的interface{}
	zeroPtr := new(T)
	zero := *zeroPtr
	operation := dbi.Between
	if req.RangeEnd == zero && req.RangeStart != zero {
		operation = dbi.Greater
		if req.Include {
			operation = dbi.GreaterOrEqual
		}
		return clause2.NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart == zero && req.RangeEnd != zero {
		operation = dbi.Less
		if req.Include {
			operation = dbi.LessOrEqual
		}
		return clause2.NewWhereClause(req.RangeField, operation, req.RangeStart)
	}
	if req.RangeStart != zero && req.RangeEnd != zero {
		if req.Include {
			return clause2.NewWhereClause(req.RangeField, operation, req.RangeStart, req.RangeEnd)
		} else {
			return clause.Where{Exprs: []clause.Expression{clause2.NewWhereClause(req.RangeField, dbi.Greater, req.RangeStart), clause2.NewWhereClause(req.RangeField, dbi.Less, req.RangeStart)}}
		}
	}
	return new(clause2.EmptyClause)
}
