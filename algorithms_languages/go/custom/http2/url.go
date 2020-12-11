package main

import (
	"fmt"
	url2 "net/url"
)

func main() {
	urlStr := fmt.Sprintf("http://%s/notifications/v2?appId=%s&cluster=%s",
		"hoper.xyz", "foo", "bar")
	url, _ := url2.Parse(urlStr)
	fmt.Println(url.RawQuery, "fq", url.ForceQuery, "fm", url.Fragment)
}
