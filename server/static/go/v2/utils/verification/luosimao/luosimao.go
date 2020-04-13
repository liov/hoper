package luosimao

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/liov/hoper/go/v2/utils/log"
)

var LuosimaoErr = errors.New("人机识别验证失败")

// LuosimaoVerify 对前端的验证码进行验证
func LuosimaoVerify(reqURL, apiKey, response string) error {
	if apiKey == "" {
		// 没有配置LuosimaoAPIKey的话，就没有验证码功能
		return nil
	}
	if response == "" {
		return errors.New("人机识别验证失败")
	}
	reqData := make(url.Values)
	reqData["api_key"] = []string{apiKey}
	reqData["response"] = []string{response}

	res, err := http.PostForm(reqURL, reqData)
	if err != nil {
		log.Error(err)
		return LuosimaoErr
	}

	defer res.Body.Close()

	resBody, readErr := ioutil.ReadAll(res.Body)

	if readErr != nil {
		log.Error(readErr)
		return LuosimaoErr
	}

	type LuosimaoResult struct {
		Error int    `json:"error"`
		Res   string `json:"res"`
		Msg   string `json:"msg"`
	}
	var luosimaoResult LuosimaoResult
	if err := json.Unmarshal(resBody, &luosimaoResult); err != nil {
		log.Error(err)
		return LuosimaoErr
	}
	if luosimaoResult.Res != "success" {
		log.Info("luosimaoResult.Res", luosimaoResult.Res)
		return LuosimaoErr
	}
	return nil
}
