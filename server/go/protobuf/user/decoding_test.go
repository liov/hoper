package user

import (
	"encoding/json"
	"testing"

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
