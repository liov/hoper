package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"net/url"
	"strconv"

	"time"
)

func Notify(c *Config) error {

	if c.DingToken == "" {
		return nil
	}
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + c.DingToken
	msg := "\\n # 发布通知 " +
		" \\n ### 项目: " + c.Repo +
		" \\n ### 操作人: " + c.CommitAuthor +
		" \\n ### 参考: " + c.CommitRef +
		" \\n ### 分支: " + c.CommitBranch +
		" \\n ### 标签: " + c.CommitTag +
		" \\n ### 时间: " + fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")) +
		" \\n ### 提交: " + c.Commit +
		" \\n ### 提交信息: " + c.CommitMessage

	msg = `{
   "msgtype": "markdown",
   "markdown": {
		"title":"发布通知",
        "text": "` + msg + `"
   },
   "at": {
       "isAtAll": false
   }
}`
	var err error
	if c.DingSecret != "" {
		//通过加签接口发布
		err = SendSignDingTalkMessage(c.DingToken, c.DingSecret, msg)
	} else {
		err = client.Post(url, msg).Do(nil)
	}

	return err
}

// SendSignDingTalkMessage 发送加签钉钉消息
//  @param title: 消息标题
//  @param content: 消息主体内容，可以使markdown格式
//  @return err
func SendSignDingTalkMessage(token, secret, msg string) error {
	if len(secret) == 0 || len(token) == 0 {
		return errors.New("钉钉配置secret或token不能为空")
	}
	// 密钥加签处理
	now := time.Now().UnixNano() / int64(time.Millisecond)
	timestampStr := strconv.FormatInt(now, 10)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestampStr + "\n" + secret))
	sum := h.Sum(nil)
	signUrl := fmt.Sprintf("robot/send?access_token=%s&timestamp=%s&sign=%s", token, timestampStr, url.QueryEscape(base64.StdEncoding.EncodeToString(sum)))

	// 发送
	err := client.Post("https://oapi.dingtalk.com/"+signUrl, msg).Do(nil)
	return err
}
