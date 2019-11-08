package time2

import (
	"fmt"
	"testing"
)

func TestTime2(t *testing.T) {
	var tm Time = 1572838282583
	fmt.Println(tm.Time())
}
