package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	message protoreflect.MessageDescriptor
}

func (m messageGenerator) Type() string {
	name := string(m.message.Name())
	// top level message
	if m.message.Parent() == m.message.ParentFile() {
		return name
	}
	return messageGenerator{message: m.message.Parent().(protoreflect.MessageDescriptor)}.Type() + "_" + name
}

func (m messageGenerator) Generate(f *codegen.File) {
	m.generateType(f)
	m.message.IsMapEntry()
	for i := 0; i < m.message.Messages().Len(); i++ {
		msg := m.message.Messages().Get(i)
		// maps are handled on field level
		if msg.IsMapEntry() {
			continue
		}
		messageGenerator{message: m.message.Messages().Get(i)}.Generate(f)
	}
	for i := 0; i < m.message.Enums().Len(); i++ {
		enumGenerator{enum: m.message.Enums().Get(i)}.Generate(f)
	}
}

func (m messageGenerator) generateType(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)
	f.P("export type ", m.Type(), " = {")
	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		fieldType := typeFromField(field)
		f.P(t(1), field.JSONName(), "?: ", fieldType.Reference(), ";")
	})
	f.P("};")
	f.P()
}
