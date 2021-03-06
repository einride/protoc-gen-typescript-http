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
	f.P("export type ", m.Type(), " = {")
	f.P("}")
	f.P()
}
