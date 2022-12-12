package stringsi

import (
	"log"
	"testing"
)

func TestUnquote(t *testing.T) {
	var s = []byte(`status:403 Forbidden{"ok":0,"errno":"100005","msg":"\u8bf7\u6c42\u8fc7\u4e8e\u9891\u7e41"}`)
	log.Println(ConvertUnicode(s))
	s = []byte(`"\u8bf7\u6c42\u8fc7\u4e8e\u9891\u7e41"`)
	log.Println(Unquote(s))
}
