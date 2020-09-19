package any

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//为了api文档，不要使用这种返回

func GenAnys(details ...proto.Message) ([]*anypb.Any, error) {
	var anys []*anypb.Any
	for _, detail := range details {
		vany, err := anypb.New(detail)
		if err != nil {
			return nil, err
		}
		anys = append(anys, vany)
	}
	return anys, nil
}
