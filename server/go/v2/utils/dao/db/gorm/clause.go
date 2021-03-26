package gormi

import "gorm.io/gorm/clause"

func Page(pageNo,pageSize int) clause.Limit{
	if pageSize == 0 || pageSize > 100 {
		pageSize = 100
	}
	if pageNo != 0 {
		return clause.Limit{Offset: (pageNo - 1) * pageSize, Limit: pageSize}
	}
	return clause.Limit{Limit: pageSize}
}
