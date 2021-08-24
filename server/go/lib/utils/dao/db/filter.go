package dbi

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Filter struct {
	Field  string `json:"field"`
	Method string `json:"method"`
	Value  string `json:"value"`
}

type Filters []Filter

func (f Filters) Build() string {
	var conditions []string
	for _, filter := range f {
		filter.Field = strings.TrimSpace(filter.Field)
		filter.Method = strings.TrimSpace(filter.Method)
		filter.Value = strings.TrimSpace(filter.Value)

		if filter.Field == "" || filter.Method == "" || filter.Value == "" {
			continue
		}

		switch filter.Method {
		case ">", "<", "=", "!=", ">=", "<=":
			conditions = append(conditions, fmt.Sprintf("%s %s '%s'", filter.Field, filter.Method, filter.Value))
		case "IN", "NOT IN":
			valueSplit := strings.Split(filter.Value, ",")
			conditions = append(conditions, fmt.Sprintf("%s %s ('%s')", filter.Field, filter.Method, strings.Join(valueSplit, "','")))
		case "BETWEEN":
			valueSplit := strings.Split(filter.Value, ",")
			if len(valueSplit) < 2 {
				continue
			}
			conditions = append(conditions, fmt.Sprintf("%s %s %s AND %s", filter.Field, filter.Method, valueSplit[0], valueSplit[1]))
		}

	}

	if len(conditions) == 0 {
		return ""
	}
	return strings.Join(conditions, " AND ")
}

func (f Filters) BuildORM(odb *gorm.DB) *gorm.DB {
	var scopes []func(db *gorm.DB) *gorm.DB
	for _, filter := range f {
		filter.Field = strings.TrimSpace(filter.Field)
		filter.Method = strings.TrimSpace(filter.Method)
		filter.Value = strings.TrimSpace(filter.Value)

		if filter.Field == "" || filter.Method == "" || filter.Value == "" {
			continue
		}

		switch filter.Method {
		case ">", "<", "=", "!=", ">=", "<=":
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("%s %s ?", filter.Field, filter.Method), filter.Value)
			})

		case "in", "not in":
			valueSplit := strings.Split(filter.Value, ",")
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("%s %s (?)", filter.Field, filter.Method), valueSplit)
			})
		case "between":
			valueSplit := strings.Split(filter.Value, ",")
			scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("%s %s ? and ?", filter.Field, filter.Method), valueSplit[0], valueSplit[1])
			})
		}
	}
	return odb.Scopes(scopes...)
}
