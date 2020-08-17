package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("E:/code/home/hoper/client/flutter/build/app/outputs/flutter-apk")))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err)
	}
}
