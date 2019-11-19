package model

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
)

func (m *User) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(m)
}