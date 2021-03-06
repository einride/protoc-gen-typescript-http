package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"github.com/einride/protoc-gen-typescript-http/internal/protowalk"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type typeGenerator struct {
	pkg   protoreflect.FullName
	files []protoreflect.FileDescriptor
}

func (t typeGenerator) Generate(f *codegen.File) {
	protowalk.WalkFiles(t.files, func(desc protoreflect.Descriptor) bool {
		if wkt, ok := WellKnownType(desc); ok {
			f.P(wkt.TypeDeclaration())
			return false
		}
		switch v := desc.(type) {
		case protoreflect.MessageDescriptor:
			if v.IsMapEntry() {
				return false
			}
			messageGenerator{pkg: t.pkg, message: v}.Generate(f)
		case protoreflect.EnumDescriptor:
			enumGenerator{pkg: t.pkg, enum: v}.Generate(f)
		}
		return true
	})
}
