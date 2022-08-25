package rpc

import (
	"testing"
	"tools/bilibili/tool"
)

func TestAPI(t *testing.T) {
	t.Run("GetView", func(t *testing.T) {
		t.Log(api.GetView(tool.Bv2av("BV15G411h7gY")))
	})
	t.Run("GetFavLResourceList", func(t *testing.T) {
		t.Log(api.GetFavLResourceList(63181530, 1))
	})
	t.Run("GetVideoInfo", func(t *testing.T) {
		t.Log(api.GetPlayerInfo(471601857, 790731181, 120))
	})
	t.Run("GetVideoInfo2", func(t *testing.T) {
		t.Log(appApi.GetPlayerInfoV2(790731181, 120))
	})
}
