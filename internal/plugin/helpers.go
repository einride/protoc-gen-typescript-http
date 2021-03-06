package plugin

import (
	"encoding/binary"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func rangeFields(message protoreflect.MessageDescriptor, f func(field protoreflect.FieldDescriptor)) {
	for i := 0; i < message.Fields().Len(); i++ {
		f(message.Fields().Get(i))
	}
}

func t(n int) string {
	return strings.Repeat("\t", n)
}

func descriptorSourcePath(desc protoreflect.Descriptor) protoreflect.SourcePath {
	if _, ok := desc.(protoreflect.FileDescriptor); ok {
		return nil
	}
	var path protoreflect.SourcePath
	switch v := desc.(type) {
	case protoreflect.FieldDescriptor:
		path = protoreflect.SourcePath{2, int32(v.Index())}
	case protoreflect.MessageDescriptor:
		path = protoreflect.SourcePath{4, int32(v.Index())}
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
