package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	message protoreflect.MessageDescriptor
}

func (m messageGenerator) Type() string {
	return string(m.message.Name())
}

func (m messageGenerator) Generate(f *codegen.File) {
	m.generateType(f)
}

func (m messageGenerator) generateType(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)
	f.P("export type ", m.Type(), " = {")
	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		f.P(t(1), fieldName(field), "?: ", fieldType(field), fieldCardinality(field), ";")
	})
	f.P("};")
	f.P()
}
