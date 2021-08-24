package main

import "github.com/liov/hoper/server/go/lib/tools/protoc-gen-enum/command"

func main() {
	command.Write(command.Generate(command.Read()))
}
