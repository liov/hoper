package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"io"
	"log"
	"strconv"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	reader, err := pro.R(pro.Conf.Pro.CommonUrl + strconv.Itoa(510780))
	if err != nil {
		log.Fatal(err)
	}
	data, _ := io.ReadAll(reader)
	log.Println(string(data))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
	}

	fmt.Println(pro.ParseHtml(doc))
}
