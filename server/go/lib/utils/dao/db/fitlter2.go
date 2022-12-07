package dbi

import (
	"database/sql/driver"
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/dao/db/const"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Operation int

const (
	OperationPlace Operation = iota
	Equal
	NotEqual
	Greater
	Less
	Between
	GreaterOrEqual
	LessOrEqual
	IsNotNull
	IsNull
	In
	NotIn
	LIKE
)

func (m Operation) SQL() string {
	switch m {
	case Equal:
		return "= ?"
	case NotEqual:
		return "!= ?"
	case Greater:
		return "> ?"
	case Less:
		return "< ?"
	case Between:
		return "BETWEEN ? AND ?"
	case GreaterOrEqual:
		return ">= ?"
	case LessOrEqual:
		return "<= ?"
	case IsNull:
		return "IS NULL"
	case IsNotNull:
		return "IS NOT NULL"
	case In:
		return "IN (?)"
	case NotIn:
		return "NOT IN (?)"
	case LIKE:
		return "LIKE ?"
	}
	return ""
}

func (m Operation) String() string {
	switch m {
	case Equal:
		return " = "
	case In:
		return " IN "
	case Between:
		return " BETWEEN "
	case Greater:
		return " > "
	case Less:
		return " < "
	case NotEqual:
		return " != "

	case GreaterOrEqual:
		return " >= "
	case LessOrEqual:
		return " <= "
	case IsNull:
		return " IS NULL"
	case IsNotNull:
		return " IS NOT NULL"
	case NotIn:
		return " NOT IN "
	}
	return "="
}

type Expression struct {
	Field     string        `json:"field"`
	Operation Operation     `json:"method"`
	Value     []interface{} `json:"value"`
}

type Exprs []Expression

func (f Exprs) Build() string {
	var conditions []string
	for _, filter := range f {
		filter.Field = strings.TrimSpace(filter.Field)

		if filter.Field == "" || filter.Operation == 0 || len(filter.Value) == 0 {
			continue
		}

		switch filter.Operation {
		case Greater, Less, Equal, NotEqual, GreaterOrEqual, LessOrEqual:
			conditions = append(conditions, filter.Field+filter.Operation.String()+ConvertParams(filter.Value[0], "'"))
		case In, NotIn:
			var vars = make([]string, len(filter.Value))
			for idx, v := range filter.Value {
				vars[idx] = ConvertParams(v, "'")
			}
			conditions = append(conditions, filter.Field+filter.Operation.String()+"("+strings.Join(vars, ",")+")")
		case Between:
			if len(filter.Value) < 2 {
				continue
			}
			var vars = make([]string, len(filter.Value))
			for idx, v := range filter.Value {
				vars[idx] = ConvertParams(v, "'")
			}
			conditions = append(conditions, filter.Field+filter.Operation.String()+vars[0]+" AND "+vars[1])
		case LIKE:
			conditions = append(conditions, filter.Field+filter.Operation.String()+ConvertParams(filter.Value[0], "'"))
		case IsNull, IsNotNull:
			conditions = append(conditions, filter.Field+filter.Operation.String())
		}

	}

	if len(conditions) == 0 {
		return ""
	}
	return strings.Join(conditions, " AND ")
}

func ConvertParams(v interface{}, escaper string) string {
	switch v := v.(type) {
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		if v.IsZero() {
			return escaper + _const.TmFmtZero + escaper
		} else {
			return escaper + v.Format(_const.TmFmtWithMS) + escaper
		}
	case *time.Time:
		if v != nil {
			if v.IsZero() {
				return escaper + _const.TmFmtZero + escaper
			} else {
				return escaper + v.Format(_const.TmFmtWithMS) + escaper
			}
		} else {
			return _const.NullStr
		}
	case driver.Valuer:
		reflectValue := reflect.ValueOf(v)
		if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
			r, _ := v.Value()
			ConvertParams(r, escaper)
		} else {
			return _const.NullStr
		}
	case fmt.Stringer:
		reflectValue := reflect.ValueOf(v)
		if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
			return escaper + strings.Replace(fmt.Sprintf("%v", v), escaper, "\\"+escaper, -1) + escaper
		} else {
			return _const.NullStr
		}
	case []byte:
		if isPrintable(v) {
			return escaper + strings.Replace(string(v), escaper, "\\"+escaper, -1) + escaper
		} else {
			return escaper + "<binary>" + escaper
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return utils.ToString(v)
	case float64, float32:
		return fmt.Sprintf("%.6f", v)
	case string:
		return escaper + strings.Replace(v, escaper, "\\"+escaper, -1) + escaper
	default:
		rv := reflect.ValueOf(v)
		if v == nil || !rv.IsValid() || rv.Kind() == reflect.Ptr && rv.IsNil() {
			return _const.NullStr
		} else if valuer, ok := v.(driver.Valuer); ok {
			v, _ = valuer.Value()
			ConvertParams(v, escaper)
		} else if rv.Kind() == reflect.Ptr && !rv.IsZero() {
			ConvertParams(reflect.Indirect(rv).Interface(), escaper)
		} else {
			for _, t := range convertableTypes {
				if rv.Type().ConvertibleTo(t) {
					return ConvertParams(rv.Convert(t).Interface(), escaper)
				}
			}
			return escaper + strings.Replace(fmt.Sprint(v), escaper, "\\"+escaper, -1) + escaper
		}
	}
	return ""
}

var convertableTypes = []reflect.Type{reflect.TypeOf(time.Time{}), reflect.TypeOf(false), reflect.TypeOf([]byte{})}

func isPrintable(s []byte) bool {
	for _, r := range s {
		if !unicode.IsPrint(rune(r)) {
			return false
		}
	}
	return true
}

func (f Exprs) BuildORM(odb *gorm.DB) *gorm.DB {
	var scopes []func(db *gorm.DB) *gorm.DB
	for _, filter := range f {
		filter.Field = strings.TrimSpace(filter.Field)

		if filter.Field == "" || filter.Operation == 0 || len(filter.Value) == 0 {
			continue
		}

		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where(filter.Field+" "+filter.Operation.SQL(), filter.Value...)
		})
	}
	return odb.Scopes(scopes...)
}

func (f Exprs) BuildSQL() (string, []interface{}) {
	var builder strings.Builder
	var vars []interface{}
	for i, filter := range f {
		if filter.Field == "" || filter.Operation == 0 || len(filter.Value) == 0 {
			continue
		}
		builder.WriteString(filter.Field)
		builder.WriteByte(' ')
		builder.WriteString(filter.Operation.SQL())
		if i < len(f) {
			builder.WriteString(" AND")
		}
		vars = append(vars, filter.Value...)
	}
	return "", nil
}
