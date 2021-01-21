package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

func main() {
	name := "xxx"
	urlStr := `http://open.api.tianyancha.com/services/open/ic/baseinfo/2.0?` + url.Values{"name": []string{name}}.Encode()
	newReq, _ := http.NewRequest("GET", urlStr, nil)

	newReq.Header.Set("HeaderAuthorization", "token")
	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		panic("panic")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("panic")
	}

	viper.SetConfigType("json")
	err = viper.ReadConfig(bytes.NewBuffer(body))
	if err != nil {
		panic("panic")
	}

	if viper.GetFloat64("error_code") == 300000 {
		panic("panic")
	}

	if viper.GetFloat64("error_code") != 0 {
		log.Println(viper.Get("reason"))
	}

	creditCode := viper.Get("result.creditCode")
	if creditCode == nil {
		log.Println("panic")
	}
	log.Println(creditCode)
}
