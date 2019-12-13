package json

import (
	"encoding/json"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/magiconair/properties/assert"
)

type Foo struct {
	a int
	b string
	c json.RawMessage
}

func TestJson(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		extra.SupportPrivateFields()
		foo := Foo{a: 1, b: "str"}
		data, _ := Json.Marshal(foo)
		t.Log(string(data))
		var f Foo
		Json.Unmarshal(data, &f)
		t.Log(f)
		assert.Equal(t, string(data), `{"a":1,"b":"str","c":null}`)
	})
}
