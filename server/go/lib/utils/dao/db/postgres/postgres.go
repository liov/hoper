package postgres

import "strings"

func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "ERROR: duplicate key")
}
