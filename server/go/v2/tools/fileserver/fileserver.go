package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("D:/hoper/client/flutter/build/app/outputs/flutter-apk")))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err)
	}
}
