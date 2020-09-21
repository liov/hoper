package main

func one(a, b []byte) []byte {
	var tmp = make([]byte, 256)
	for _, c1 := range b {
		if tmp[int(c1)] != 1 {
			tmp[int(c1)] = 1
		}
	}
	var c []byte
	for _, c2 := range a {
		if tmp[int(c2)] == 0 {
			c = append(c, c2)
		}
	}
	return c
}

func two(a []int) {
	//for i:=0;i<len(a)/2;
}

func three(a, b []int) float32 {
	m := len(a)
	n := len(b)
	if m > n {
		a, b = b, a
	}
	half := (m + n) / 2
	e := (m+n)&1 == 0
	var i, j int
	var left = 0
	var right = len(a)
	for {
		i = (left + right) / 2
		j = half - i
		if i < right && a[i] < b[j-1] {
			left = i + 1
		} else if i > left && a[i-1] > b[j] {
			right = i - 1
		} else {
			var x, y int
			if i == len(a) {
				x = b[j]
			} else if j == len(b) {
				x = a[i]
			} else {
				if b[j] > a[i-1] {
					x = a[i-1]
				} else {
					x = b[j]
				}
			}
			if !e {
				return float32(x)
			}
			if i == 0 {
				y = b[j-1]
			} else if j == 0 {
				y = a[i-1]
			} else {
				if b[j] > a[i-1] {
					y = b[j]
				} else {
					y = a[i-1]
				}
			}
			return float32((x + y) / 2)
		}
	}
}

func main() {
	println(three([]int{1, 3, 5, 7}, []int{2, 4, 6}))
}

type point struct {
	x, y int
}

func four(a [][]int) int {
	var get bool
	var count int
	var m = make(map[point]int)
	for x := 0; x < len(a); x++ {
		for y := 0; y < len(a[0]); y++ {
			if a[x][y] == 2 {
				if x+1 < len(a) && y+1 < len(a[0]) && a[x+1][y+1] == 1 {
					a[x+1][y+1] = 2
					get = true
				}
				if x-1 >= 0 && y-1 >= 0 && a[x-1][y-1] == 1 {
					a[x-1][y-1] = 2
					get = true
				}

				if x-1 >= 0 && y+1 < len(a[0]) && a[x-1][y+1] == 1 {
					a[x-1][y+1] = 2
					get = true
				}
				if y-1 >= 0 && x+1 < len(a) && a[x-1][y+1] == 1 {
					a[x+1][y-1] = 2
					get = true
				}
			} else if a[x][y] == 1 {
				m[point{x, y}] = a[x][y]
			}
		}
	}
}

func pass(a [][]int, x, y int) int {
	if x == len(a)-1 || y == len(a[0])-1 {
		return 0
	}
	var get bool
	var count int
	if a[x][y] == 2 {
		if a[x+1][y+1] == 1 {
			a[x+1][y+1] = 2
			get = true
			count += pass(a, x+1, y+1)
		}
		if x-1 >= 0 && y-1 >= 0 && a[x-1][y-1] == 1 {
			a[x-1][y-1] = 2
			get = true
			count += pass(a, x+1, y+1)
		}
		if x-1 >= 0 && a[x-1][y+1] == 1 {
			a[x-1][y+1] = 2
			count += pass(a, x+1, y+1)
			get = true
		}
		if y-1 >= 0 && a[x-1][y+1] == 1 {
			a[x+1][y-1] = 2
			count += pass(a, x+1, y+1)
			get = true
		}
	}
	if get {
		count += 1
	}
	return count
}
