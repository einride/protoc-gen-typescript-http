package plugin

import (
	"encoding/binary"
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

type fileSourcePath = int32

const (
	fileSourcePathMessage fileSourcePath = 4
	fileSourcePathEnum                   = 5
)

type messageSourcePath = int32

const (
	messageSourcePathField   messageSourcePath = 2
	messageSourcePathMessage                   = 3
	messageSourcePathEnum                      = 4
)

type enumSourcePath = int32

const (
	enumSourcePathValue enumSourcePath = 2
)

func descriptorSourcePath(desc protoreflect.Descriptor) protoreflect.SourcePath {
	if _, ok := desc.(protoreflect.FileDescriptor); ok {
		return nil
	}
	var path protoreflect.SourcePath
	switch desc.Parent().(type) {
	case protoreflect.FileDescriptor:
		switch v := desc.(type) {
		case protoreflect.MessageDescriptor:
			path = protoreflect.SourcePath{fileSourcePathMessage, int32(v.Index())}
		case protoreflect.EnumDescriptor:
			path = protoreflect.SourcePath{fileSourcePathEnum, int32(v.Index())}
		}
	case protoreflect.MessageDescriptor:
		switch v := desc.(type) {
		case protoreflect.FieldDescriptor:
			path = protoreflect.SourcePath{messageSourcePathField, int32(v.Index())}
		case protoreflect.MessageDescriptor:
			path = protoreflect.SourcePath{messageSourcePathMessage, int32(v.Index())}
		case protoreflect.EnumDescriptor:
			path = protoreflect.SourcePath{messageSourcePathEnum, int32(v.Index())}
		}
	case protoreflect.EnumDescriptor:
		switch v := desc.(type) {
		case protoreflect.EnumValueDescriptor:
			path = protoreflect.SourcePath{enumSourcePathValue, int32(v.Index())}
		}
	}
	return append(descriptorSourcePath(desc.Parent()), path...)
}

func descriptorSourceLocation(desc protoreflect.Descriptor, path protoreflect.SourcePath) (protoreflect.SourceLocation, bool) {
	locs := desc.ParentFile().SourceLocations()
	key := newPathKey(path)
	for i := 0; i < locs.Len(); i++ {
		loc := locs.Get(i)
		if newPathKey(loc.Path) == key {
			return loc, true
		}
	}
	return protoreflect.SourceLocation{}, false
}

// A pathKey is a representation of a location path suitable for use as a map key.
type pathKey string

// newPathKey converts a location path to a pathKey.
func newPathKey(path protoreflect.SourcePath) pathKey {
	buf := make([]byte, 4*len(path))
	for i, x := range path {
		binary.LittleEndian.PutUint32(buf[i*4:], uint32(x))
	}
	return pathKey(buf)
}
