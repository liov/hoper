package verification

import (
	"errors"
	"net/http"

	"github.com/liov/hoper/go/v2/utils/net/http/client"
)

var LuosimaoErr = errors.New("人机识别验证失败")

type LuosimaoResult struct {
	Error int    `json:"error"`
	Res   string `json:"res"`
	Msg   string `json:"msg"`
}

func (l *LuosimaoResult) CheckError() error {
	if l.Res != "success" {
		return LuosimaoErr
	}
	return nil
}

// LuosimaoVerify 对前端的验证码进行验证
func LuosimaoVerify(reqURL, apiKey, response string) error {
	if reqURL == "" || apiKey == "" {
		// 没有配置LuosimaoAPIKey的话，就没有验证码功能
		return nil
	}
	if response == "" {
		return LuosimaoErr
	}

	req := struct {
		ApiKey   string `json:"api_key"`
		Response string `json:"response"`
	}{
		ApiKey:   apiKey,
		Response: response,
	}
	result :=new(LuosimaoResult)

	err := client.NewRequest(reqURL, http.MethodPost, &req).
		SetContentType(client.ContentTypeForm).HTTPRequest(result)
	if err != nil {
		return err
	}
	return nil
}
