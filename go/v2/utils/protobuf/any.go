package protobuf

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func  GenAnys(details ...proto.Message) ([]*any.Any, error) {
	var anys []*any.Any
	for _, detail := range details {
		vany, err := ptypes.MarshalAny(detail)
		if err != nil {
			return nil, err
		}
		anys = append(anys, vany)
	}
	return anys, nil
}

func GenAny(detail proto.Message) (*any.Any, error) {
	vany, err := ptypes.MarshalAny(detail)
	if err != nil {
		return nil, err
	}
	return vany,nil
}

func  GenGogoAnys(details ...proto.Message) ([]*types.Any, error) {
	var anys []*types.Any
	for _, detail := range details {
		vany, err := types.MarshalAny(detail)
		if err != nil {
			return nil, err
		}
		anys = append(anys, vany)
	}
	return anys, nil
}

func GenGogoAny(detail proto.Message) (*types.Any, error) {
	vany, err := types.MarshalAny(detail)
	if err != nil {
		return nil, err
	}
	return vany,nil
}