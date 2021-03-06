package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceGenerator struct {
	pkg     protoreflect.FullName
	service protoreflect.ServiceDescriptor
}

func (s serviceGenerator) Generate(f *codegen.File) {
	s.generateInterface(f)
}

func (s serviceGenerator) generateInterface(f *codegen.File) {
	commentGenerator{descriptor: s.service}.generateLeading(f, 0)
	f.P("export interface ", descriptorTypeName(s.service), " {")
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		if !supportedMethod(method) {
			return
		}
		commentGenerator{descriptor: method}.generateLeading(f, 1)
		input := typeFromMessage(s.pkg, method.Input())
		output := typeFromMessage(s.pkg, method.Output())
		f.P(t(1), method.Name(), "(", input.Reference(), "): Promise<", output.Reference(), ">")
	})
	f.P("}")
	f.P()
}

func supportedMethod(method protoreflect.MethodDescriptor) bool {
	return !method.IsStreamingClient() && !method.IsStreamingServer()
}
