package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"github.com/einride/protoc-gen-typescript-http/internal/protowalk"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type packageGenerator struct {
	pkg   protoreflect.FullName
	files []protoreflect.FileDescriptor
}

func (p packageGenerator) Generate(f *codegen.File) {
	protowalk.WalkFiles(p.files, func(desc protoreflect.Descriptor) bool {
		if wkt, ok := WellKnownType(desc); ok {
			f.P(wkt.TypeDeclaration())
			return false
		}
		switch v := desc.(type) {
		case protoreflect.MessageDescriptor:
			if v.IsMapEntry() {
				return false
			}
			messageGenerator{pkg: p.pkg, message: v}.Generate(f)
		case protoreflect.EnumDescriptor:
			enumGenerator{pkg: p.pkg, enum: v}.Generate(f)
		case protoreflect.ServiceDescriptor:
			serviceGenerator{pkg: p.pkg, service: v}.Generate(f)
		}
		return true
	})
}
