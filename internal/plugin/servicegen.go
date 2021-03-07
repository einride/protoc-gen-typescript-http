package plugin

import (
	"strconv"
	"strings"

	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"github.com/einride/protoc-gen-typescript-http/internal/httprule"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceGenerator struct {
	pkg        protoreflect.FullName
	genHandler bool
	service    protoreflect.ServiceDescriptor
}

func (s serviceGenerator) Generate(f *codegen.File) error {
	s.generateInterface(f)
	if s.genHandler {
		s.generateHandler(f)
	}
	if err := s.generateClient(f); err != nil {
		return err
	}
	return nil
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

func (s serviceGenerator) generateClient(f *codegen.File) error {
	f.P(
		"export function create",
		descriptorTypeName(s.service),
		"Client(handler: requestHandler): ",
		descriptorTypeName(s.service),
		" {",
	)
	f.P(t(1), "return {")
	var methodErr error
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		if err := s.generateMethod(f, method); err != nil {
			methodErr = err
		}
	})
	if methodErr != nil {
		return methodErr
	}
	f.P(t(1), "}")
	f.P("}")
	return nil
}

func (s serviceGenerator) generateMethod(f *codegen.File, method protoreflect.MethodDescriptor) error {
	outputType := typeFromMessage(s.pkg, method.Output())
	r, ok := httprule.Get(method)
	if !ok {
		return nil
	}
	rule, err := httprule.ParseRule(r)
	if err != nil {
		return err
	}
	f.P(t(2), method.Name(), "(request) {")
	s.generateMethodPathValidation(f, method.Input(), rule)
	s.generateMethodPath(f, method.Input(), rule)
	s.generateMethodBody(f, method.Input(), rule)
	s.generateMethodQuery(f, method.Input(), rule)
	f.P(t(3), "let uri = path")
	f.P(t(3), "if (hasQuery) {")
	f.P(t(4), "uri += \"?\" + query.toString()")
	f.P(t(3), "}")
	f.P(t(3), "return handler(uri, ", strconv.Quote(rule.Method), ", body) as Promise<", outputType.Reference(), ">")
	f.P(t(2), "},")
	return nil
}

func (s serviceGenerator) generateMethodPathValidation(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule *httprule.Rule,
) {
	for _, seg := range rule.Template.Segments {
		if seg.Kind != httprule.SegmentKindVariable {
			continue
		}
		fp := seg.Variable.FieldPath
		nullPath := nullPropagationPath(fp, input)
		protoPath := strings.Join(fp, ".")
		errMsg := "missing required field request." + protoPath
		f.P(t(3), "if (!request.", nullPath, ") {")
		f.P(t(4), "throw new Error(", strconv.Quote(errMsg), ")")
		f.P(t(3), "}")
	}
}

func (s serviceGenerator) generateMethodPath(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule *httprule.Rule,
) {
	pathParts := make([]string, 0, len(rule.Template.Segments))
	for _, seg := range rule.Template.Segments {
		switch seg.Kind {
		case httprule.SegmentKindVariable:
			fieldPath := jsonPath(seg.Variable.FieldPath, input)
			pathParts = append(pathParts, "${request."+fieldPath+"}")
		case httprule.SegmentKindLiteral:
			pathParts = append(pathParts, seg.Literal)
		case httprule.SegmentKindMatchSingle: // TODO: Double check this and following case
			pathParts = append(pathParts, "*")
		case httprule.SegmentKindMatchMultiple:
			pathParts = append(pathParts, "**")
		}
	}
	path := strings.Join(pathParts, "/")
	if rule.Template.Verb != "" {
		path += ":" + rule.Template.Verb
	}
	f.P(t(3), "const path = `", path, "`")
}

func (s serviceGenerator) generateMethodBody(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule *httprule.Rule,
) {
	switch {
	case rule.Body == "":
		f.P(t(3), "const body = null;")
	case rule.Body == "*":
		f.P(t(3), "const body = JSON.stringify(request);")
	default:
		nullPath := nullPropagationPath(httprule.FieldPath{rule.Body}, input)
		f.P(t(3), "const body = JSON.stringify(request?.", nullPath, " ?? {})")
	}
}

func (s serviceGenerator) generateMethodQuery(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule *httprule.Rule,
) {
	f.P(t(3), "const query = new URLSearchParams();")
	// nothing in query
	if rule.Body == "*" {
		f.P(t(3), "const hasQuery = false;")
		return
	}
	// fields covered by path
	pathCovered := make(map[string]struct{})
	f.P(t(3), "let hasQuery = false;")
	walkJsonLeafFields(input, func(path httprule.FieldPath, field protoreflect.FieldDescriptor) {
		if _, ok := pathCovered[path.String()]; ok {
			return
		}
		if rule.Body != "" && path[0] == rule.Body {
			return
		}
		nullPath := nullPropagationPath(path, input)
		jp := jsonPath(path, input)
		f.P(t(3), "if (request.", nullPath, ") {")
		f.P(t(4), "hasQuery = true;")
		switch {
		case field.IsList():
			f.P(t(4), "for (const x of request.", jp, ") {")
			f.P(t(5), "query.append(", strconv.Quote(jp), ", x.toString());")
			f.P(t(4), "}")
		default:
			f.P(t(4), "query.set(", strconv.Quote(jp), ", request.", jp, ".toString());")
		}
		f.P(t(3), "}")
	})
}

func supportedMethod(method protoreflect.MethodDescriptor) bool {
	_, ok := httprule.Get(method)
	return ok && !method.IsStreamingClient() && !method.IsStreamingServer()
}

func jsonPath(path httprule.FieldPath, message protoreflect.MessageDescriptor) string {
	return strings.Join(jsonPathSegments(path, message), ".")
}

func nullPropagationPath(path httprule.FieldPath, message protoreflect.MessageDescriptor) string {
	return strings.Join(jsonPathSegments(path, message), "?.")
}

func jsonPathSegments(path httprule.FieldPath, message protoreflect.MessageDescriptor) []string {
	segs := make([]string, len(path))
	for i, p := range path {
		field := message.Fields().ByName(protoreflect.Name(p))
		segs[i] = field.JSONName()
		if i < len(path) {
			message = field.Message()
		}
	}
	return segs
}
