package main

import (
	"log"
	"net/http"

	_ "github.com/liov/hoper/go/v2/httptpl/httptpl2/internal/service"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
)

func main() {
	/*	f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		trace.Start(f)
		defer trace.Stop()*/
	router := pick.NewEasyRouter(false, "httptpl")
	router.ServeFiles("/static", "E:/")
	log.Println("visit http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
