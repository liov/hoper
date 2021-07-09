package number

import (
	"fmt"
	"strconv"
)

func TwoDecimalPlaces(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
