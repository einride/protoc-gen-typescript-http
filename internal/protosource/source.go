package protosource

import (
	"encoding/binary"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Path returns the source path to the descriptor in its parent file.
func Path(desc protoreflect.Descriptor) protoreflect.SourcePath {
	if _, ok := desc.(protoreflect.FileDescriptor); ok {
		return nil
	}
	// TODO: add support for all "vectors" unless upstream proto has added this.
	// https://github.com/protocolbuffers/protobuf-go/blob/v1.25.0/reflect/protoreflect/source.go#L16
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
		if v, ok := desc.(protoreflect.EnumValueDescriptor); ok {
			path = protoreflect.SourcePath{enumSourcePathValue, int32(v.Index())}
		}
	}
	return append(Path(desc.Parent()), path...)
}

// Location returns the source location, if any, at the given source path.
func Location(
	desc protoreflect.Descriptor,
	path protoreflect.SourcePath,
) (protoreflect.SourceLocation, bool) {
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

type fileSourcePath = int32

const (
	fileSourcePathMessage fileSourcePath = 4
	fileSourcePathEnum    fileSourcePath = 5
)

type messageSourcePath = int32

const (
	messageSourcePathField   messageSourcePath = 2
	messageSourcePathMessage messageSourcePath = 3
	messageSourcePathEnum    messageSourcePath = 4
)

type enumSourcePath = int32

const (
	enumSourcePathValue enumSourcePath = 2
)

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
