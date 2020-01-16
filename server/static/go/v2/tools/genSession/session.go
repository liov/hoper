package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type Session struct {
	UserID       int    `json:"userId"`
	UserName     string `json:"userName"`
	UserRealName string `json:"userRealName"`
	ClientIp     string `json:"clientIp"`
}

func main() {
	sess := &Session{UserID: 699}
	data, _ := json.Marshal(sess)
	log.Println("X-Session-Data", base64.StdEncoding.EncodeToString(data))
}
