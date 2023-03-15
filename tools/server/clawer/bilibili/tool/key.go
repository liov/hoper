package tool

import (
	"github.com/liov/hoper/server/go/lib/utils/slices"
	"strings"
)

func GetAppKey(entropy string) (appkey, sec string) {
	revEntropy := slices.ReverseRunes([]rune(entropy))
	for i := range revEntropy {
		revEntropy[i] = revEntropy[i] + 2
	}
	ret := strings.Split(string(revEntropy), ":")

	return ret[0], ret[1]
}
