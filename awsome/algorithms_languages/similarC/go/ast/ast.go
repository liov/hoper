package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// src是我们要打印AST的输入。
	src := `
package main
var (
	Action_name = map[int32]string{
		0: "Signup",
		1: "Active",
		2: "RestPassword",
		3: "EditPassword",
		4: "CreateResume",
		5: "EditResume",
		6: "DELETEResume",
	}
)
`

	// 通过解析src来创建AST。
	fset := token.NewFileSet() // 职位相对于fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// 打印AST。
	ast.Print(fset, f)

}
