package json

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/json-iterator/go/extra"
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
		reflect.DeepEqual(string(data), `{"a":1,"b":"str","c":null}`)
	})
}

func TestJson2(t *testing.T) {
	data := []byte(`{"getUser":{"details":{"name":"","id":1,"gender":男,"phone":""}}}`)
	var j = &graphql.Response{
		Data: data,
	}
	b, err := json.Marshal(j)
	if err != nil {
		log.Println(err)
	}
	log.Println(b)
	var j2 = json.RawMessage(data)
	b, err = json.Marshal(j2)
	if err != nil {
		log.Println(err)
	}
	log.Println(b)
}
