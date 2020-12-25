package main

import (
	"log"
	"net/http"

	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	_ "github.com/liov/hoper/go/v2/utils/net/http/pick/_example/service"
)

func main() {
	router := pick.New(false, "httptpl")
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
