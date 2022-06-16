package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	//VERSION is SDK version
	VERSION = "0.1"

	//ROOT is the root url
	ROOT = "https://oapi.dingtalk.com/"

	ContentTypeTextTmpl     = `{"msgtype":"text","text":{"content":"%s"}}`
	ContentTypeMarkdownTmpl = `{"msgtype":"markdown","markdown":{"title":"%s","text":"%s"}}`
)

type ContentType int

const (
	ContentTypeText ContentType = iota
	ContentTypeMarkdown
)

func (c ContentType) String() string {
	switch c {
	case ContentTypeText:
		return "text"
	case ContentTypeMarkdown:
		return "markdown"
	default:
		return ""
	}
}

func (c ContentType) Tmpl() string {
	switch c {
	case ContentTypeText:
		return ContentTypeTextTmpl
	case ContentTypeMarkdown:
		return ContentTypeMarkdownTmpl
	default:
		return ContentTypeTextTmpl
	}
}

func (c ContentType) Body(title, content string) string {
	switch c {
	case ContentTypeText:
		return fmt.Sprintf(c.Tmpl(), content)
	case ContentTypeMarkdown:
		return fmt.Sprintf(c.Tmpl(), title, content)
	default:
		return fmt.Sprintf(c.Tmpl(), content)
	}
}

func SendRobotMessage(accessToken, secret, title, content string, contentType ContentType) error {
	signUrl, err := RobotUrl(accessToken, secret)
	if err != nil {
		return err
	}
	body := strings.NewReader(contentType.Body(title, content))

	return client.Post(ROOT+signUrl, body).Do(nil)
}

func RobotUrl(accessToken, secret string) (string, error) {
	if accessToken == "" {
		return "", errors.New("token不能为为空")
	}
	if secret != "" {
		// 密钥加签处理
		now := time.Now().UnixNano() / int64(time.Millisecond)
		timestampStr := strconv.FormatInt(now, 10)
		h := hmac.New(sha256.New, []byte(secret))
		sum := h.Sum(nil)
		return fmt.Sprintf("robot/send?access_token=%s&timestamp=%s&sign=%s", accessToken, timestampStr, url.QueryEscape(base64.StdEncoding.EncodeToString(sum))), nil
	}
	return fmt.Sprintf("robot/send?access_token=%s", accessToken), nil
}

//SendRobotTextMessage can send a text message to a group chat
func SendRobotTextMessage(accessToken string, content string) error {
	return SendRobotMessage(accessToken, "", "", content, ContentTypeText)
}

func SendRobotTextMessageWithSecret(accessToken, secret, content string) error {
	if accessToken == "" {
		return errors.New("secret不能为空")
	}
	return SendRobotMessage(accessToken, secret, "", content, ContentTypeText)
}

func SendRobotMarkDownMessage(token, title, content string) error {
	return SendRobotMessage(token, "", title, content, ContentTypeMarkdown)
}

func SendRobotMarkDownMessageWithSecret(token, secret, title, content string) error {
	if len(secret) == 0 {
		return errors.New("secret不能为空")
	}
	return SendRobotMessage(token, secret, title, content, ContentTypeMarkdown)
}
