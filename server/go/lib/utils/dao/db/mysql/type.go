package mysql

import "strings"

func MySqlTypeToGoType(typ string) string {
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
