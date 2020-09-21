package logrus

import (
	"fmt"
	"net/url"
)

func Logrus() {
	url := url.PathEscape("测试")
	fmt.Println(url)
}
