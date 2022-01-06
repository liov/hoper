package main

import "github.com/actliboy/hoper/server/go/lib/tools/protoc-gen-enum/command"

func main() {
	command.Write(command.Generate(command.Read()))
}
