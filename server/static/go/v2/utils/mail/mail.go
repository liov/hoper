package mail

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/smtp"
	"text/template"

	"github.com/liov/hoper/go/v2/utils/log"
)

type Message struct {
	Name, Mail, Subject, ContentType, Content, ToMail string
}

const msg = `To: {{.ToMail}}
From: {{.Name}} <{{.Mail}}>
Subject: {{.Subject}}
Content-Type: {{if .ContentType}}{{.ContentType}}{{- else}}text/html; charset=UTF-8{{end}}
{{.Content}}
`

func GenMsg(m *Message) []byte {
	t := template.Must(template.New("msg").Parse(msg))
	var buf = new(bytes.Buffer)
	err := t.Execute(buf, m)
	if err != nil {
		log.Error("executing template:", err)
	}
	return buf.Bytes()
}

func SendMailTLS(addr string, auth smtp.Auth, from string, to []string, message []byte) error {
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
