package plugin

import (
	"go.einride.tech/protoc-gen-typescript-http/internal/httprule"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type jsonLeafWalkFunc func(path httprule.FieldPath, field protoreflect.FieldDescriptor)

func walkJSONLeafFields(message protoreflect.MessageDescriptor, f jsonLeafWalkFunc) {
	var w jsonWalker
	w.walkMessage(nil, message, f)
}

type jsonWalker struct {
	seen map[protoreflect.FullName]struct{}
}

func (w *jsonWalker) enter(name protoreflect.FullName) bool {
	if _, ok := w.seen[name]; ok {
		return false
	}
	if w.seen == nil {
		w.seen = make(map[protoreflect.FullName]struct{})
	}
	w.seen[name] = struct{}{}
	return true
}

func (w *jsonWalker) walkMessage(path httprule.FieldPath, message protoreflect.MessageDescriptor, f jsonLeafWalkFunc) {
	if w.enter(message.FullName()) {
		for i := 0; i < message.Fields().Len(); i++ {
			field := message.Fields().Get(i)
			p := append(httprule.FieldPath{}, path...)
			p = append(p, string(field.Name()))
			switch {
			case !field.IsMap() && field.Kind() == protoreflect.MessageKind:
				if IsWellKnownType(field.Message()) {
					f(p, field)
				} else {
					w.walkMessage(p, field.Message(), f)
				}
			default:
				f(p, field)
			}
		}
	}
}
