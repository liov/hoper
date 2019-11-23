package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("%s\n","foo")
	s:=fmt.Sprintf("%s bar","foo")
	println(s)//内置print输出顺序并不能确定
/*	file,_:=os.Create("interview/foo.txt")
	defer file.Close()*/
	fmt.Fprintf(os.Stdout,"%s\n","foo")
}
