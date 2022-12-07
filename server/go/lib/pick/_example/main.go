package main

import (
	pick2 "github.com/liov/hoper/server/go/lib/pick"
	service2 "github.com/liov/hoper/server/go/lib/pick/_example/service"
	"log"
	"net/http"

	_ "github.com/liov/hoper/server/go/lib/pick/_example/service"
)

func main() {
	pick2.RegisterService(&service2.UserService{}, &service2.TestService{}, &service2.StaticService{})
	router := pick2.New(false, "httptpl")
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
