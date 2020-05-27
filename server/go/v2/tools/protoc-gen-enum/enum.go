package main

import "github.com/liov/hoper/go/v2/tools/protoc-gen-enum/command"

func main() {
	command.Write(command.Generate(command.Read()))
}
