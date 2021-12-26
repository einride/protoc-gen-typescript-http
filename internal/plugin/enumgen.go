package plugin

import (
	"strconv"

	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type enumGenerator struct {
	pkg  protoreflect.FullName
	enum protoreflect.EnumDescriptor
}

func (e enumGenerator) Generate(f *codegen.File) {
	commentGenerator{descriptor: e.enum}.generateLeading(f, 0)
	f.P("export type ", scopedDescriptorTypeName(e.pkg, e.enum), " =")
	if e.enum.Values().Len() == 1 {
		commentGenerator{descriptor: e.enum.Values().Get(0)}.generateLeading(f, 1)
		f.P(t(1), strconv.Quote(string(e.enum.Values().Get(0).Name())), ";")
		return
	}
	rangeEnumValues(e.enum, func(value protoreflect.EnumValueDescriptor, last bool) {
		commentGenerator{descriptor: value}.generateLeading(f, 1)
		if last {
			f.P(t(1), "| ", strconv.Quote(string(value.Name())), ";")
		} else {
			f.P(t(1), "| ", strconv.Quote(string(value.Name())))
		}
	})
}
