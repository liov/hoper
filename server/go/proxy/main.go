package main

import (
	"github.com/hopeio/utils/log"
	"github.com/hopeio/utils/net/http/proxy"
)

func main() {
	log.Fatal(proxy.Director(":8080"))
}
