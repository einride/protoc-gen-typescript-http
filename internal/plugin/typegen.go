package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type typeGenerator struct {
	files []protoreflect.FileDescriptor
}

func (t typeGenerator) Generate(f *codegen.File) {
	t.generateMessages(f)
	t.generateEnums(f)
	t.generateWKTs(f)
}

func (t typeGenerator) generateMessages(f *codegen.File) {
	for _, file := range t.files {
		messages := file.Messages()
		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			messageGenerator{message: message}.Generate(f)
		}
	}
}

func (t typeGenerator) generateEnums(f *codegen.File) {
	for _, file := range t.files {
		enums := file.Enums()
		for i := 0; i < enums.Len(); i++ {
			enum := enums.Get(i)
			enumGenerator{enum: enum}.Generate(f)
		}
	}
}

func (t typeGenerator) generateWKTs(f *codegen.File) {
	wkts := t.collectWKTs()
	for _, wkt := range wkts {
		f.P(wkt.TypeDeclaration())
	}
}

func (t typeGenerator) collectWKTs() []WellKnown {
	wkts := make([]WellKnown, 0, 10)
	seen := make(map[WellKnown]struct{})
	collectWKT := func(wkt WellKnown) {
		if _, ok := seen[wkt]; ok {
			return
		}
		seen[wkt] = struct{}{}
		wkts = append(wkts, wkt)
	}
	collectField := func(field protoreflect.FieldDescriptor) {
		if field.IsMap() {
			return
		}
		switch field.Kind() {
		case protoreflect.MessageKind:
			if wkt, ok := WellKnownType(field.Message()); ok {
				collectWKT(wkt)
			}
		case protoreflect.EnumKind:
			if wkt, ok := WellKnownType(field.Enum()); ok {
				collectWKT(wkt)
			}
		}
	}
	collectMessage := func(message protoreflect.MessageDescriptor) {
		rangeFields(message, func(field protoreflect.FieldDescriptor) {
			switch {
			case field.IsMap():
				collectField(field.MapValue())
			default:
				collectField(field)
			}
		})
	}
	for _, file := range t.files {
		rangeFileMessages(file, collectMessage)
	}
	return wkts
}
