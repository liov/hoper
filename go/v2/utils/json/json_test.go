package json

import (
	"testing"

	"github.com/json-iterator/go/extra"
)

type Foo struct {
	a int
	b string
}

func TestJson(t *testing.T)  {
	t.Run("test", func(t *testing.T) {
		extra.SupportPrivateFields()
		foo:=Foo{a:1,b:"str"}
		data,_:= Json.Marshal(foo)
		t.Log(string(data))
		var f Foo
		Json.Unmarshal(data,&f)
		t.Log(f)
	})
}
