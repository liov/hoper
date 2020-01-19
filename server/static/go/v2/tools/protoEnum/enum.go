package main

import "github.com/liov/hoper/go/v2/tools/protoEnum/command"

func main() {
	command.Write(command.Generate(command.Read()))
}
