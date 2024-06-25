package plugin

import (
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageGenerator struct {
	pkg     protoreflect.FullName
	message protoreflect.MessageDescriptor
}

func (m messageGenerator) Generate(f *codegen.File) {
	commentGenerator{descriptor: m.message}.generateLeading(f, 0)
	f.P("export type ", scopedDescriptorTypeName(m.pkg, m.message), " = {")
	rangeFields(m.message, func(field protoreflect.FieldDescriptor) {
		commentGenerator{descriptor: field}.generateLeading(f, 1)
		optional := hasOptionalAnnotation(field) || field.HasOptionalKeyword()
		fieldType := typeFromField(m.pkg, field)
		if field.ContainingOneof() == nil && !optional {
			f.P(t(1), field.JSONName(), ": ", fieldType.Reference(), " | undefined;")
		} else {
			f.P(t(1), field.JSONName(), "?: ", fieldType.Reference(), ";")
		}
	})

	f.P("};")
	f.P()
}

func hasOptionalAnnotation(field protoreflect.FieldDescriptor) bool {
	return fieldbehavior.Has(field, annotations.FieldBehavior_OPTIONAL)
}
