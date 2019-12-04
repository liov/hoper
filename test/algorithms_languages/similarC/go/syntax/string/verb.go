package main

import (
	"fmt"
	"os"
)

type User struct {
	name string
	age  int
}

var valF float64 = 32.9983
var valI int = 89
var valS string = "Go is an open source programming language that makes it easy to build simple,  reliable,  and efficient software."
var valB bool = true

func main() {

	p := User{"John", 28}

	fmt.Printf("Printf struct %%v : %v\n", p)
	fmt.Printf("Printf struct %%+v : %+v\n", p)
	fmt.Printf("Printf struct %%#v : %#v\n", p)
	fmt.Printf("Printf struct %%T : %T\n", p)

	fmt.Printf("Printf struct %%p : %p\n", &p)

	fmt.Printf("Printf float64 %%v : %v\n", valF)
	fmt.Printf("Printf float64 %%+v : %+v\n", valF)
	fmt.Printf("Printf float64 %%#v : %#v\n", valF)
	fmt.Printf("Printf float64 %%T : %T\n", valF)
	fmt.Printf("Printf float64 %%f : %f\n", valF)
	fmt.Printf("Printf float64 %%4.3f : %4.3f\n", valF)
	fmt.Printf("Printf float64 %%8.3f : %8.3f\n", valF)
	fmt.Printf("Printf float64 %%-8.3f : %-8.3f\n", valF)
	fmt.Printf("Printf float64 %%e : %e\n", valF)
	fmt.Printf("Printf float64 %%E : %E\n", valF)

	fmt.Printf("Printf int %%v : %v\n", valI)
	fmt.Printf("Printf int %%+v : %+v\n", valI)
	fmt.Printf("Printf int %%#v : %#v\n", valI)
	fmt.Printf("Printf int %%T : %T\n", valI)
	fmt.Printf("Printf int %%d : %d\n", valI)
	fmt.Printf("Printf int %%8d : %8d\n", valI)
	fmt.Printf("Printf int %%-8d : %-8d\n", valI)
	fmt.Printf("Printf int %%b : %b\n", valI)
	fmt.Printf("Printf int %%c : %c\n", valI)
	fmt.Printf("Printf int %%o : %o\n", valI)
	fmt.Printf("Printf int %%U : %U\n", valI)
	fmt.Printf("Printf int %%q : %q\n", valI)
	fmt.Printf("Printf int %%x : %x\n", valI)

	fmt.Printf("Printf string %%v : %v\n", valS)
	fmt.Printf("Printf string %%+v : %+v\n", valS)
	fmt.Printf("Printf string %%#v : %#v\n", valS)
	fmt.Printf("Printf string %%T : %T\n", valS)
	fmt.Printf("Printf string %%x : %x\n", valS)
	fmt.Printf("Printf string %%X : %X\n", valS)
	fmt.Printf("Printf string %%s : %s\n", valS)
	fmt.Printf("Printf string %%200s : %200s\n", valS)
	fmt.Printf("Printf string %%-200s : %-200s\n", valS)
	fmt.Printf("Printf string %%q : %q\n", valS)

	fmt.Printf("Printf bool %%v : %v\n", valB)
	fmt.Printf("Printf bool %%+v : %+v\n", valB)
	fmt.Printf("Printf bool %%#v : %#v\n", valB)
	fmt.Printf("Printf bool %%T : %T\n", valB)
	fmt.Printf("Printf bool %%t : %t\n", valB)

	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)

	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}
