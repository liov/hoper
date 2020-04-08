package any

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

//为了api文档，不要使用这种返回

func GenAnys(details ...proto.Message) ([]*any.Any, error) {
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
	return vany, nil
}

func GenGogoAnys(details ...proto.Message) ([]*types.Any, error) {
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
	return vany, nil
}

func ResolveAny(typeUrl string) (proto.Message, error) {
	// Only the part of typeUrl after the last slash is relevant.
	mname := typeUrl
	if slash := strings.LastIndex(mname, "/"); slash >= 0 {
		mname = mname[slash+1:]
	}
	mt := proto.MessageType(mname)
	if mt == nil {
		return nil, fmt.Errorf("unknown message type %q", mname)
	}
	return reflect.New(mt.Elem()).Interface().(proto.Message), nil
}
