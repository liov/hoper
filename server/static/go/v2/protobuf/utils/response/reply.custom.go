package response

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/liov/hoper/go/v2/utils/protobuf/any"
)

func (m *BytesReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(m)
}

func (m *StringReply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(m)
}

type GoReply struct {
	Code    int32
	Message string
	Details proto.Message
}

func (m *Reply) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	var err error
	reply := GoReply{Code: m.Code, Message: m.Message}
	reply.Details, err = any.ResolveAny(m.Details.TypeUrl)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(m.Details.Value, reply.Details)
	if err != nil {
		return nil, err
	}
	return json.Marshal(reply)
}
