package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	reader, err := pro.Request(http.DefaultClient, pro.Conf.Pro.CommonUrl+strconv.Itoa(510707))
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(reader)
	log.Println(string(data))
	defer reader.Close()
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	fmt.Println(pro.ParseHtml(doc))
}
