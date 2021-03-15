package plugin

import "google.golang.org/protobuf/reflect/protoreflect"

type Type struct {
	IsNamed bool
	Name    string

	IsList     bool
	IsMap      bool
	Underlying *Type
}

func (t Type) Reference() string {
	switch {
	case t.IsMap:
		return "{ [key: string]: " + t.Underlying.Reference() + " }"
	case t.IsList:
		return t.Underlying.Reference() + "[]"
	default:
		return t.Name
	}
}

func typeFromField(pkg protoreflect.FullName, field protoreflect.FieldDescriptor) Type {
	switch {
	case field.IsMap():
		underlying := namedTypeFromField(pkg, field.MapValue())
		return Type{
			IsMap:      true,
			Underlying: &underlying,
		}
	case field.IsList():
		underlying := namedTypeFromField(pkg, field)
		return Type{
			IsList:     true,
			Underlying: &underlying,
		}
	default:
		return namedTypeFromField(pkg, field)
	}
}

func namedTypeFromField(pkg protoreflect.FullName, field protoreflect.FieldDescriptor) Type {
	switch field.Kind() {
	case protoreflect.StringKind, protoreflect.BytesKind:
		return Type{IsNamed: true, Name: "string"}
	case protoreflect.BoolKind:
		return Type{IsNamed: true, Name: "boolean"}
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
		protoreflect.Sint32Kind,
		protoreflect.Sint64Kind,
		protoreflect.FloatKind:
		return Type{IsNamed: true, Name: "number"}
	case protoreflect.MessageKind:
		return typeFromMessage(pkg, field.Message())
	case protoreflect.EnumKind:
		desc := field.Enum()
		if wkt, ok := WellKnownType(field.Enum()); ok {
			return Type{IsNamed: true, Name: wkt.Name()}
		}
		return Type{IsNamed: true, Name: scopedDescriptorTypeName(pkg, desc)}
	default:
		return Type{IsNamed: true, Name: "unknown"}
	}
}

func typeFromMessage(pkg protoreflect.FullName, message protoreflect.MessageDescriptor) Type {
	if wkt, ok := WellKnownType(message); ok {
		return Type{IsNamed: true, Name: wkt.Name()}
	}
	return Type{IsNamed: true, Name: scopedDescriptorTypeName(pkg, message)}
}
