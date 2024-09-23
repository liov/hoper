package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target := "http://example.com" // 目标服务器的地址
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
