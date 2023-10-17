package plugin

import (
	"strconv"

	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type enumGenerator struct {
	opts Options
	pkg  protoreflect.FullName
	enum protoreflect.EnumDescriptor
}

func (e enumGenerator) Generate(f *codegen.File) {
	commentGenerator{opts: e.opts, descriptor: e.enum}.generateLeading(f, 0)
	if e.opts.UseEnumNumbers {
		f.P("export enum ", scopedDescriptorTypeName(e.pkg, e.enum), " {")

		rangeEnumValues(e.enum, func(value protoreflect.EnumValueDescriptor, last bool) {
			commentGenerator{opts: e.opts, descriptor: value}.generateLeading(f, 1)

			name := string(value.Name())
			name = TextToCase(name, e.opts.EnumFieldNaming)

			f.P(t(1), name, " = ", value.Number(), ",")
		})

		f.P("}")
	} else {
		f.P("export type ", scopedDescriptorTypeName(e.pkg, e.enum), " =")
		if e.enum.Values().Len() == 1 {
			commentGenerator{opts: e.opts, descriptor: e.enum.Values().Get(0)}.generateLeading(f, 1)
			f.P(t(1), strconv.Quote(string(e.enum.Values().Get(0).Name())), ";")
			return
		}
		rangeEnumValues(e.enum, func(value protoreflect.EnumValueDescriptor, last bool) {
			commentGenerator{opts: e.opts, descriptor: value}.generateLeading(f, 1)
			if last {
				f.P(t(1), "| ", strconv.Quote(string(value.Name())), ";")
			} else {
				f.P(t(1), "| ", strconv.Quote(string(value.Name())))
			}
		})
	}
}
