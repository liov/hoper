package main

import (
	"github.com/bdragon300/tusgo"
	"github.com/hopeio/utils/log"
	"net/http"
	"net/url"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:8080/files/")
	cl := tusgo.NewClient(http.DefaultClient, baseURL)
	u := tusgo.Upload{}
	_, err := cl.GetUpload(&u, "http://localhost:8080/files/9473b67eb7b1af7c1be6a79abb228e19")
	if err != nil {
		return
	}
	log.Info(u)
}
