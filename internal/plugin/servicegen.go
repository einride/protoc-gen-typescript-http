package plugin

import (
	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceGenerator struct {
	pkg        protoreflect.FullName
	genHandler bool
	service    protoreflect.ServiceDescriptor
}

func (s serviceGenerator) Generate(f *codegen.File) {
	s.generateInterface(f)
	if s.genHandler {
		s.generateHandler(f)
	}
	s.generateClient(f)
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
		f.P(t(1), method.Name(), "(request: ", input.Reference(), "): Promise<", output.Reference(), ">")
	})
	f.P("}")
	f.P()
}

func (s serviceGenerator) generateHandler(f *codegen.File) {
	f.P("type requestHandler = (path: string, method: string, body: string | null) => Promise<unknown>")
	f.P()
}

func (s serviceGenerator) generateClient(f *codegen.File) {
	f.P(
		"export function create",
		descriptorTypeName(s.service),
		"Client(handler: requestHandler): ",
		descriptorTypeName(s.service),
		" {",
	)
	f.P(t(1), "return {")
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		s.generateMethod(f, method)
	})
	f.P(t(1), "}")
	f.P("}")
}

func (s serviceGenerator) generateMethod(f *codegen.File, method protoreflect.MethodDescriptor) {
	outputType := typeFromMessage(s.pkg, method.Output())
	f.P(t(2), method.Name(), "(request) {")
	f.P(t(3), "return handler(\"\", \"\", null) as Promise<", outputType.Reference(), ">")
	f.P(t(2), "},")
}

func supportedMethod(method protoreflect.MethodDescriptor) bool {
	return !method.IsStreamingClient() && !method.IsStreamingServer()
}
