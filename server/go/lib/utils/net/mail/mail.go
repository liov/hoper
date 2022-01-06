package mail

import (
	"bytes"
	"crypto/tls"
	"net"
	"net/smtp"

	templatei "github.com/actliboy/hoper/server/go/lib/utils/def/template"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

//550,Mailbox not found or access denied.是因为收件邮箱不存在
type Mail struct {
	FromName, From, Subject, ContentType, Content string
	To                                            []string
}

const msg = `{{define "mail"}}To: {{join .To ",\n\t"}}
From: {{.FromName}} <{{.From}}>
Subject: {{.Subject}}
Content-Type: {{if .ContentType}}{{.ContentType}}{{- else}}text/html; charset=UTF-8{{end}}

{{.Content}}{{end}}
`

func init() {
	templatei.Parse(msg)
}

func (m *Mail) GenMsg() []byte {
	var buf = new(bytes.Buffer)
	err := templatei.Execute(buf, "mail", m)
	if err != nil {
		log.Error("executing template:", err)
	}
	return buf.Bytes()
}

func (m *Mail) SendMailTLS(addr string, auth smtp.Auth) error {
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
	if err := client.Mail(m.From); err != nil {
		return err
	}
	for _, addr := range m.To {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(m.GenMsg())
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
