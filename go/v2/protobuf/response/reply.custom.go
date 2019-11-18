package response

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
)

func (m *AnyReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error)  {
	return json.Marshal(m)
}