package main

import (
	"github.com/liov/hoper/go/v2/tailmon/pick"
	"log"
	"net/http"

	"github.com/liov/hoper/go/v2/tailmon/pick/_example/service"
	_ "github.com/liov/hoper/go/v2/tailmon/pick/_example/service"
)

func main() {
	router := pick.New(false, "httptpl")
	pick.RegisterService(&service.UserService{}, &service.TestService{}, &service.StaticService{})
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
