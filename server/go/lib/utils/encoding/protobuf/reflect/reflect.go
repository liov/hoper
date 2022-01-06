package reflect

import (
	cuzproto "github.com/actliboy/hoper/server/go/lib/protobuf/utils/policy"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Redact(pb proto.Message) {
	m := pb.ProtoReflect()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		if proto.GetExtension(opts, cuzproto.E_NonSensitive).(bool) {
			return true
		}
		m.Clear(fd)
		return true
	})
}
