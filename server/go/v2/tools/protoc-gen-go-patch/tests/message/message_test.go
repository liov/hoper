package message

import (
	"testing"

	"github.com/liov/protopatch2/tests"
)

func TestBasicMessage(t *testing.T) {
	tests.ValidateMessage(t, &BasicMessage{})
}

func TestOneofMessage(t *testing.T) {
	m := &OneofMessage{}
	tests.ValidateMessage(t, m)
	var _ isOneofMessage_Contents = &OneofMessage_Id{}
	var _ isOneofMessage_Contents = &OneofMessage_Name{}
	var _ int32 = m.GetId()
	var _ string = m.GetName()
}

func TestNestedMessage(t *testing.T) {
	tests.ValidateMessage(t, &OuterMessage{})
	tests.ValidateMessage(t, &OuterMessage_InnerMessage{})
}

func TestRenamedMessage(t *testing.T) {
	tests.ValidateMessage(t, &RenamedMessage{})
}

func TestRenamedOneofMessage(t *testing.T) {
	m := &RenamedOneofMessage{}
	tests.ValidateMessage(t, m)
	var _ isRenamedOneofMessage_Contents = &RenamedOneofMessage_Id{}
	var _ isRenamedOneofMessage_Contents = &RenamedOneofMessage_Name{}
	var _ int32 = m.GetId()
	var _ string = m.GetName()
}

func TestRenamedOuterMessage(t *testing.T) {
	tests.ValidateMessage(t, &RenamedOuterMessage{})
	tests.ValidateMessage(t, &RenamedOuterMessage_InnerMessage{})
}

func TestRenamedInnerMessage(t *testing.T) {
	tests.ValidateMessage(t, &OuterMessageWithRenamedInnerMessage{})
	tests.ValidateMessage(t, &RenamedInnerMessage{})
}

func TestMessageWithRenamedField(t *testing.T) {
	m := &MessageWithRenamedField{}
	tests.ValidateMessage(t, m)
	var _ int32 = m.ID
	var _ int32 = m.GetID()
}

func TestMessageWithStructTags(t *testing.T) {
	m := &MessageWithTags{}
	tests.ValidateTag(t, m, "Value", "test", "value")
}

func TestNestedMessageWithStructTags(t *testing.T) {
	m := &OuterMessageWithTags_InnerMessage{}
	tests.ValidateTag(t, m, "Value", "test", "value")
}

func TestMessageWithOverrideTag(t *testing.T) {
	m := &MessageWithOverrideTag{}
	tests.ValidateTag(t, m, "Value1", "json", "value1")
	tests.ValidateTag(t, m, "Value2", "json", "-")
	tests.ValidateTag(t, m, "Value3", "json", "value")
	tests.ValidateTag(t, m, "Protobuf", "json", "protobuf,omitempty")
	tests.ValidateTag(t, m, "Protobuf", "protobuf", "bytes,4,opt,name=proto,proto3")
}

func TestNestedMessageWithOverrideTag(t *testing.T) {
	m := &OuterMessageOverrideTag_InnerMessage{}
	tests.ValidateTag(t, m, "Value1", "json", "value1")
	tests.ValidateTag(t, m, "Value2", "json", "-")
	tests.ValidateTag(t, m, "Value3", "json", "value")
	tests.ValidateTag(t, m, "Protobuf", "json", "protobuf,omitempty")
	tests.ValidateTag(t, m, "Protobuf", "protobuf", "bytes,4,opt,name=proto,proto3")
}
