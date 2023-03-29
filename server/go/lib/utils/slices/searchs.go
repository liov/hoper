package slices

import (
	"github.com/liov/hoper/server/go/lib/utils/def/interface"
	"golang.org/x/exp/constraints"
)

// BinarySearch 二分查找
func BinarySearch[V constraints.Ordered](arr []_interface.CmpKey[V], x _interface.CmpKey[V]) int {
	l, r := 0, len(arr)-1
	for l <= r {
		mid := (l + r) / 2
		if arr[mid].CmpKey() == x.CmpKey() {
			return mid
		} else if x.CmpKey() > arr[mid].CmpKey() {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}
