package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"tools/pro"
)

func main() {
	reader, _ := pro.Request(http.DefaultClient, fmt.Sprintf(pro.CommonUrl, "375509"))
	data, _ := ioutil.ReadAll(reader)
	log.Println(string(data))
}
