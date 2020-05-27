package slices

func Intersection(a []uint64, b []uint64) []uint64 {
	if len(a) > len(b) {
		return intersection(a, b)
	}
	return intersection(b, a)
}

func intersection(a []uint64, b []uint64) []uint64 {
	var ret []uint64
	for _, x := range a {
		if In(x, b) {
			ret = append(ret, x)
		}
	}
	return ret
}

func OrderArrayIntersection(a []uint64, b []uint64) []uint64 {
	var ret []uint64
	if len(a) == 0 || len(b) == 0 {
		return nil
	}
	var idx int
	for _, x := range a {
		if x > b[len(b)-1] {
			return ret
		}
		for j := idx; idx < len(b)-1; j++ {
			if a[len(a)-1] < b[idx] {
				return ret
			}
			if x == b[idx] {
				ret = append(ret, x)
				idx = j
			}
		}
	}
	return ret
}

func In(a uint64, b []uint64) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}

func Union(a []uint64, b []uint64) []uint64 {
	var m = make(map[uint64]struct{}, len(a)+len(b))
	for _, x := range a {
		m[x] = struct{}{}
	}
	for _, x := range b {
		m[x] = struct{}{}
	}
	var ret = make([]uint64, len(m))
	for k, _ := range m {
		ret = append(ret, k)
	}
	return ret
}

func Difference(a []uint64, b []uint64) []uint64 {
	var ret []uint64
	if len(a) > len(b) {
		for _, x := range a {
			if !In(x, b) {
				ret = append(ret, x)
			}
		}
	} else {
		for _, x := range b {
			if !In(x, a) {
				ret = append(ret, x)
			}
		}
	}
	return ret
}

type CmpKey interface {
	CmpKey() uint64
}
