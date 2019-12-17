package mail

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenMsg(t *testing.T) {
	var msg = Message{
		Name:    "liov",
		Mail:    "liov@github.com",
		Subject: "测试",
		Content: "邮件",
		ToMail:  strings.Join([]string{"test1@mail.com", "test2@mail.com"}, ","),
	}
	bytes := GenMsg(&msg)
	fmt.Println(string(bytes))
}
