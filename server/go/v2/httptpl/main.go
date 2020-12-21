package main

import (
	"log"
	"net/http"

	"github.com/liov/hoper/go/v2/httptpl/conf"
	_ "github.com/liov/hoper/go/v2/httptpl/service"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func main() {
	defer initialize.Start(conf.Conf, nil)()
	router := pick.New(false, "httptpl")
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
