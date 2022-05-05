package gofound

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"net/http"
)

// IndexDoc 索引实体
type IndexDoc struct {
	Id       uint32                 `json:"id,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

type IndexDocInterface struct {
	Id       uint32      `json:"id,omitempty"`
	Text     string      `json:"text,omitempty"`
	Document interface{} `json:"document,omitempty"`
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc

	Keys []uint32 `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	Score float32 `json:"score,omitempty"` //得分
}

type Result struct {
	State bool `json:"state"`

	Message string `json:"message,omitempty"`

	Data interface{} `json:"data,omitempty"`
}

type RemoveIndexModel struct {
	Id uint32 `json:"id,omitempty"`
}

func foundCall(url, method string, param interface{}) {
	var res Result
	err := client.NewRequest(url, method, param).SetHeader("Authorization", Token).SetLogger(nil).Do(&res)
	if err != nil {
		log.Error(err)
	}
}

func AddIndex(id uint32, text string, document interface{}) {
	foundCall("/api/index", http.MethodGet, &IndexDocInterface{
		Id:       id,
		Text:     text,
		Document: document,
	})
}
