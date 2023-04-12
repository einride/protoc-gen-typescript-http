package plugin

import (
	"fmt"
	"strconv"
	"strings"

	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"go.einride.tech/protoc-gen-typescript-http/internal/httprule"
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
	return s.generateClient(f)
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
		f.P(t(1), method.Name(), "(request: ", input.Reference(), "): Promise<", output.Reference(), ">;")
	})
	f.P("}")
	f.P()
}

func (s serviceGenerator) generateHandler(f *codegen.File) {
	f.P("type RequestType = {")
	f.P(t(1), "path: string;")
	f.P(t(1), "method: string;")
	f.P(t(1), "body: string | null;")
	f.P("};")
	f.P()
	f.P("type RequestHandler = (request: RequestType, meta: { service: string, method: string }) => Promise<unknown>;")
	f.P()
}

func (s serviceGenerator) generateClient(f *codegen.File) error {
	f.P(
		"export function create",
		descriptorTypeName(s.service),
		"Client(",
		"\n",
		t(1),
		"handler: RequestHandler",
		"\n",
		"): ",
		descriptorTypeName(s.service),
		" {",
	)
	f.P(t(1), "return {")
	var methodErr error
	rangeMethods(s.service.Methods(), func(method protoreflect.MethodDescriptor) {
		if err := s.generateMethod(f, method); err != nil {
			methodErr = fmt.Errorf("generate method %s: %w", method.Name(), err)
		}
	})
	if methodErr != nil {
		return methodErr
	}
	f.P(t(1), "};")
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
		return fmt.Errorf("parse http rule: %w", err)
	}
	f.P(t(2), method.Name(), "(request) { // eslint-disable-line @typescript-eslint/no-unused-vars")
	s.generateMethodPathValidation(f, method.Input(), rule)
	s.generateMethodPath(f, method.Input(), rule)
	s.generateMethodBody(f, method.Input(), rule)
	s.generateMethodQuery(f, method.Input(), rule)
	f.P(t(3), "let uri = path;")
	f.P(t(3), "if (queryParams.length > 0) {")
	f.P(t(4), "uri += `?${queryParams.join(\"&\")}`")
	f.P(t(3), "}")
	f.P(t(3), "return handler({")
	f.P(t(4), "path: uri,")
	f.P(t(4), "method: ", strconv.Quote(rule.Method), ",")
	f.P(t(4), "body,")
	f.P(t(3), "}, {")
	f.P(t(4), "service: \"", method.Parent().Name(), "\",")
	f.P(t(4), "method: \"", method.Name(), "\",")
	f.P(t(3), "}) as Promise<", outputType.Reference(), ">;")
	f.P(t(2), "},")
	return nil
}

func (s serviceGenerator) generateMethodPathValidation(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule httprule.Rule,
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
		f.P(t(4), "throw new Error(", strconv.Quote(errMsg), ");")
		f.P(t(3), "}")
	}
}

func (s serviceGenerator) generateMethodPath(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule httprule.Rule,
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
	f.P(t(3), "const path = `", path, "`; // eslint-disable-line quotes")
}

func (s serviceGenerator) generateMethodBody(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule httprule.Rule,
) {
	switch {
	case rule.Body == "":
		f.P(t(3), "const body = null;")
	case rule.Body == "*":
		f.P(t(3), "const body = JSON.stringify(request);")
	default:
		nullPath := nullPropagationPath(httprule.FieldPath{rule.Body}, input)
		f.P(t(3), "const body = JSON.stringify(request?.", nullPath, " ?? {});")
	}
}

func (s serviceGenerator) generateMethodQuery(
	f *codegen.File,
	input protoreflect.MessageDescriptor,
	rule httprule.Rule,
) {
	f.P(t(3), "const queryParams: string[] = [];")
	// nothing in query
	if rule.Body == "*" {
		return
	}
	// fields covered by path
	pathCovered := make(map[string]struct{})
	for _, segment := range rule.Template.Segments {
		if segment.Kind != httprule.SegmentKindVariable {
			continue
		}
		pathCovered[segment.Variable.FieldPath.String()] = struct{}{}
	}
	walkJSONLeafFields(input, func(path httprule.FieldPath, field protoreflect.FieldDescriptor) {
		if _, ok := pathCovered[path.String()]; ok {
			return
		}
		if rule.Body != "" && path[0] == rule.Body {
			return
		}
		nullPath := nullPropagationPath(path, input)
		jp := jsonPath(path, input)
		f.P(t(3), "if (request.", nullPath, ") {")
		switch {
		case field.IsList():
			f.P(t(4), "request.", jp, ".forEach((x) => {")
			f.P(t(5), "queryParams.push(`", jp, "=${encodeURIComponent(x.toString())}`)")
			f.P(t(4), "})")
		default:
			f.P(t(4), "queryParams.push(`", jp, "=${encodeURIComponent(request.", jp, ".toString())}`)")
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
