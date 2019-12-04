package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// Serve static files from current working directory
//  it will start at port 8080 if port is being used it will try next one
//GOOS=windows GOARCH=amd64 go build -o serve.exe serve.go

func main() {
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("err opening port")
	}
	fmt.Printf("Listening at %s\n", addr)
	log.Fatal(http.Serve(listener, logger(http.FileServer(http.Dir(".")))))

}

func logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
