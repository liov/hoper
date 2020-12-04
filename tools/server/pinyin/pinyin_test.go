package py

import (
	"testing"

	"github.com/mozillazg/go-pinyin"
)

func TestPinyin(t *testing.T) {
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "N",
			want: `n`,
		},
		{
			name: "【",
			want: `【`,
		},
		{
			name: "[",
			want: `[`,
		},
		{
			name: ",",
			want: `,`,
		},
		{
			name: "。",
			want: `。`,
		},
		{
			name: "中",
			want: `z`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := FistLetter(test.name, a); got != test.want {
				t.Fatalf(" got: %s\nwant: %s\n", got, test.want)
			}
		})
	}

}
