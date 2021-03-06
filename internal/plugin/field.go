package plugin

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

func fieldName(field protoreflect.FieldDescriptor) string {
	return field.JSONName()
}

func fieldType(field protoreflect.FieldDescriptor) interface{} {
	switch field.Kind() {
	case protoreflect.StringKind, protoreflect.BytesKind:
		return "string"
	case protoreflect.BoolKind:
		return "boolean"
	case
		protoreflect.Int32Kind,
		protoreflect.Int64Kind,
		protoreflect.Uint32Kind,
		protoreflect.Uint64Kind,
		protoreflect.DoubleKind,
		protoreflect.Fixed32Kind,
		protoreflect.Fixed64Kind,
		protoreflect.Sfixed32Kind,
		protoreflect.Sfixed64Kind,
		protoreflect.FloatKind:
		return "number"
	default:
		return "unknown"
	}
}

func fieldCardinality(field protoreflect.FieldDescriptor) string {
	if field.Cardinality() == protoreflect.Repeated {
		return "[]"
	}
	return ""
}
