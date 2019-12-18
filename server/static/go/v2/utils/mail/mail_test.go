package mail

import (
	"fmt"
	"testing"
)

func TestGenMsg(t *testing.T) {
	var msg = Mail{
		FromName: "liov",
		From:     "liov@github.com",
		Subject:  "测试",
		Content:  "邮件",
		To:       []string{"test1@mail.com", "test2@mail.com"},
	}
	bytes := GenMsg(&msg)
	fmt.Println(string(bytes))
}
