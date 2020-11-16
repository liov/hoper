package slices

import "github.com/liov/hoper/go/v2/utils/def"

// BinarySearch 二分查找
func BinarySearch(arr []def.CmpKey, x def.CmpKey) int {
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
