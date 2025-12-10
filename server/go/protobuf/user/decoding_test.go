package user

import (
	"encoding/json/v2"
	"testing"

	"github.com/hopeio/protobuf/model"
	"google.golang.org/protobuf/encoding/protojson"
)

func BenchmarkProtoJson(b *testing.B) {
	gen := User{}
	for i := 0; i < b.N; i++ {
		protojson.Marshal(&gen)
	}
}

func BenchmarkStdJson(b *testing.B) {
	gen := User{}
	for i := 0; i < b.N; i++ {
		json.Marshal(&gen)
	}
}

func TestProtoJson(t *testing.T) {
	gen := User{Basic: &model.Model{Id: 1}}
	data, err := protojson.Marshal(&gen)
	t.Log(string(data), err)
}

func TestJson(t *testing.T) {
	gen := User{Basic: &model.Model{Id: 1}}
	data, err := json.Marshal(&gen)
	t.Log(string(data), err)
}
