package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init()  {
	extra.SupportPrivateFields()
}

var Json = jsoniter.ConfigDefault

