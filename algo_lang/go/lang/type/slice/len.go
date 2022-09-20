package main

func main() {
	var x *struct {
		s [][32]int
	}
	println(len(x.s[99]))
}
