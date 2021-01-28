package main

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/runtime/protoimpl"
)

type Foo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func TestSize(t *testing.T) {
	t.Log(reflect.TypeOf(Foo{}).Size())
	t.Log(reflect.TypeOf(protoimpl.MessageState{}).Size())
	t.Log(reflect.TypeOf(protoimpl.SizeCache(1)).Size())
	t.Log(reflect.TypeOf(protoimpl.UnknownFields{}).Size())
	t.Log(reflect.TypeOf("1").Size())
}
