package mail

import (
	"crypto/tls"
	"net"
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/log"
)

func SendMailTLS(addr string, auth smtp.Auth, from string,to []string, message []byte) error {
	client, err := createSMTPClient(addr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer client.Close()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	client.Quit()
	return nil
}

func createSMTPClient(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
