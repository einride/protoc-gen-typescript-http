package plugin

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func rangeFields(message protoreflect.MessageDescriptor, f func(field protoreflect.FieldDescriptor)) {
	for i := 0; i < message.Fields().Len(); i++ {
		f(message.Fields().Get(i))
	}
}

func t(n int) string {
	return strings.Repeat("\t", n)
}
