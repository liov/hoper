package initialize

import (
	"crypto/tls"
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P3Mail() *smtp.Client {
	conf := MailConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &conf); !exist {
		return nil
	}

	addr:= conf.Host+ conf.Port
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Error(err)
	}

	client, err := smtp.NewClient(conn, conf.Host)
	if err != nil {
		log.Error(err)
	}
	auth := smtp.PlainAuth("", conf.User, conf.Password, conf.Host)
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				log.Error(err)
			}
		}
	}

	return client
}
