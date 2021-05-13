package request

import (
	"strconv"
)

func Error(err error) string {
	if e, ok := err.(*strconv.NumError); ok {
		return "strconv." + e.Func + ": " + "parsing " + e.Num + ": " + e.Err.Error()
	}
	return err.Error()
}
