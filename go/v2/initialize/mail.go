package initialize

import (
	"crypto/tls"
	"net"
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P3Mail() {
	mailConf := MailConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &mailConf); !exist {
		return
	}

	addr:=mailConf.Host+mailConf.Port
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Error(err)
	}
	host, _, _ := net.SplitHostPort(addr)

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Error(err)
	}
	auth := smtp.PlainAuth("", mailConf.User, mailConf.Password, mailConf.Host)
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				log.Error(err)
			}
		}
	}

	reflect3.SetFieldValue(i.dao,client)
}
