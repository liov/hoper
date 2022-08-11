package rpc

import (
	"crypto/md5"
	"fmt"
	"tools/bilibili/tool"
)

var appApi = &AppApi{}

type AppApi struct{}

// 客户端api
func (api *AppApi) GetPlayerInfoV2(cid, qn int) *VideoInfo {
	return GetZ[VideoInfo](GetPlayerUrlV2(cid, qn))
}

const _entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"

var (
	appKey, sec = tool.GetAppKey(_entropy)
)

func GetPlayerUrlV2(cid, qn int) string {
	var _paramsTemp = "appkey=%s&cid=%d&otype=json&qn=%d&quality=%d&type="
	var _playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
	params := fmt.Sprintf(_paramsTemp, appKey, cid, qn, qn)
	chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+sec)))
	return fmt.Sprintf(_playApiTemp, params, chksum)
}
