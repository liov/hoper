package main

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/dingding"
	"time"
)

func Notify(c *Config) error {

	if c.DingToken == "" {
		return nil
	}

	msg := "\\n # 发布通知 " +
		" \\n ### 项目: " + c.Repo +
		" \\n ### 操作人: " + c.CommitAuthor +
		" \\n ### 参考: " + c.CommitRef +
		" \\n ### 分支: " + c.CommitBranch +
		" \\n ### 标签: " + c.CommitTag +
		" \\n ### 时间: " + fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")) +
		" \\n ### 提交: " + c.Commit +
		" \\n ### 提交信息: " + c.CommitMessage

	var err error
	if c.DingSecret != "" {
		err = dingding.SendRobotMarkDownMessageWithSecret(c.DingToken, c.DingSecret, "发布通知", msg)
	} else {
		err = dingding.SendRobotMarkDownMessage(c.DingToken, "发布通知", msg)
	}

	return err
}
