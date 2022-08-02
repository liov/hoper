package api

import (
	"testing"
	"tools/bilibili/parser"
)

func TestAPI(t *testing.T) {
	t.Run("GetView", func(t *testing.T) {
		api.GetView(parser.Bv2av("BV15G411h7gY"))
	})
	t.Run("GetFavList", func(t *testing.T) {
		api.GetFavList(1)
	})
}
