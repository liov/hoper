package main

import (
	"github.com/liov/hoper/go/v2/utils/dao/db/get"
	"tools/pro"
)

func main() {
	get.GetDB().Migrator().CreateTable(&pro.Post{}, &pro.InvalidPost{})
}
