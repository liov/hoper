package schema

import "github.com/gorilla/schema"

var DefaultDecoder = schema.NewDecoder()

func init() {
	DefaultDecoder.SetAliasTag("json")
}
