package output

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

func (th *LoggerCall) Write(b []byte) (n int, err error) {

	if len(th.HookURL) == 0 {
		return 0,errors.New("无效的URL")
	}

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
		return 0,err
	}
	body := strings.NewReader(string(bodyByte))
	var resp *http.Response
	for i := range th.HookURL {
		resp, err = http.Post(th.HookURL[i], "application/json", body)
		if err != nil {
			return 0,err
		}
		if resp.StatusCode != http.StatusOK {
			return 0,errors.New("上报出错")
		}
	}
	defer resp.Body.Close()
	return 0,nil
}
