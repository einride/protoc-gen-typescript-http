package plugin

import (
	"strings"

	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type commentGenerator struct {
	descriptor protoreflect.Descriptor
}

func (c commentGenerator) generateLeading(f *codegen.File, indent int) {
	path := descriptorSourcePath(c.descriptor)
	loc, ok := descriptorSourceLocation(c.descriptor, path)
	if !ok {
		return
	}
	lines := strings.Split(loc.LeadingComments, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		f.P(t(indent), "// ", strings.TrimSpace(line))
	}
}
