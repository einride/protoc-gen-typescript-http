package plugin

import (
	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	opts    Options
	pkg     protoreflect.FullName
	message protoreflect.MessageDescriptor
}

func (m messageGenerator) Generate(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)
	f.P("export type ", scopedDescriptorTypeName(m.pkg, m.message), " = {")
	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		fieldType := typeFromField(m.pkg, field)
		name := field.JSONName()
		if m.opts.UseProtoNames {
			name = field.TextName()
		}
		reference := fieldType.Reference()
		if m.opts.UseEnumNumbers && field.Kind() == protoreflect.EnumKind {
			reference = "number"
		}
		if field.ContainingOneof() == nil && !field.HasOptionalKeyword() {
			f.P(t(1), name, ": ", reference, " | undefined;")
		} else {
			f.P(t(1), name, "?: ", reference, ";")
		}
	})

	f.P("};")
	f.P()
}
