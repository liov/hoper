package user

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
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
