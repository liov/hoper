package number

import "strconv"

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}
