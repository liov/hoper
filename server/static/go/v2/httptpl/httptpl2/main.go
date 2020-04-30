package main

import (
	"log"
	"net/http"

	_ "github.com/liov/hoper/go/v2/httptpl/httptpl2/internal/service"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func main() {
	router := pick.NewEasyRouter(true, "httptpl")
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
