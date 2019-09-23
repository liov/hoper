package log

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap/zapcore"
)

type LoggerCall struct {
	HookLevel zapcore.Level
	HookURL   []string //上报地址
	AtMan     string   //@人
	CallMan   []string //手机号
}

func (th *LoggerCall) Fire(e *zapcore.Entry) error {
	if e.Level < th.HookLevel {
		return nil
	}
	if len(th.HookURL) == 0 {
		return errors.New("无效的URL")
	}
	b := e.Message
	hookBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "系统报警",
			"text":  string(b) + th.AtMan,
		},
		"at": map[string]interface{}{
			"atMobiles": th.CallMan,
			"isAtAll":   false,
		},
	}
	bodyByte, err := json.Marshal(hookBody)
	if nil != err {
		return err
	}
	body := strings.NewReader(string(bodyByte))
	var resp *http.Response
	for i := range th.HookURL {
		resp, err = http.Post(th.HookURL[i], "application/json", body)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(e.Message)
		}
	}
	defer resp.Body.Close()
	return nil
}
