package api

import (
	"testing"
	"tools/bilibili/tool"
)

func TestAPI(t *testing.T) {
	t.Run("GetView", func(t *testing.T) {
		api.GetView(tool.Bv2av("BV15G411h7gY"))
	})
	t.Run("GetFavList", func(t *testing.T) {
		api.GetFavList(1)
	})
}
