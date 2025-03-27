package main

import (
	"github.com/hopeio/utils/log"
	"github.com/hopeio/utils/net/http/proxy"
)

func main() {
	log.Fatal(proxy.DirectorServer(":8080"))
}
