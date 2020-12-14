package v1

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"github.com/liov/hoper/go/v2/utils/log/output"
	"go.uber.org/zap"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func LogFormatter(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			source          = values[1]
		)

		messages = []interface{}{source}

		if level == "sql" {
			// duration
			messages = append(messages, fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
						} else {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
						}
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages = append(messages, sql)
			messages = append(messages, fmt.Sprintf("%v", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned"))
		} else {
			messages = append(messages, values[2:]...)
		}
	}

	return
}

type loggerV1 struct {
	logger *zap.Logger
}

func (entry *loggerV1) Print(v ...interface{}) {
	msgs := LogFormatter(v...)
	if v[0].(string) == "sql" {
		entry.logger.Info("sql",
			zap.Time("@timestamp", time.Now()),
			zap.Object("@fields", &output.Kibana{
				Interface:   msgs[0].(string),
				ProcessTime: msgs[1].(time.Duration),
				Param:       msgs[2].(string),
				Result:      msgs[3].(string),
				Other:       "sql",
			}))
	}
	if v[0].(string) == "log" {
		entry.logger.Error("sqllog",
			zap.Time("@timestamp", time.Now()),
			zap.Any("@fields", &output.Kibana{
				Interface: msgs[0].(string),
				Result:    msgs[1].(error).Error(),
				Other:     "sql",
			}))
	}
}
