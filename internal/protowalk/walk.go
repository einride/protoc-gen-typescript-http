package protowalk

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type WalkFunc func(desc protoreflect.Descriptor) bool

func WalkFiles(files []protoreflect.FileDescriptor, f WalkFunc) {
	var w walker
	w.walkFiles(files, f)
}

type walker struct {
	seen map[string]struct{}
}

func (w *walker) enter(name string) bool {
	if _, ok := w.seen[name]; ok {
		return false
	}
	if w.seen == nil {
		w.seen = make(map[string]struct{})
	}
	w.seen[name] = struct{}{}
	return true
}

func (w *walker) walkFiles(files []protoreflect.FileDescriptor, f WalkFunc) {
	for _, file := range files {
		w.walkFile(file, f)
	}
}

func (w *walker) walkFile(file protoreflect.FileDescriptor, f WalkFunc) {
	if w.enter(file.Path()) {
		if !f(file) {
			return
		}
		w.walkEnums(file.Enums(), f)
		w.walkMessages(file.Messages(), f)
		w.walkServices(file.Services(), f)
	}
}

func (w *walker) walkEnums(enums protoreflect.EnumDescriptors, f WalkFunc) {
	for i := 0; i < enums.Len(); i++ {
		w.walkEnum(enums.Get(i), f)
	}
}

func (w *walker) walkEnum(enum protoreflect.EnumDescriptor, f WalkFunc) {
	if w.enter(string(enum.FullName())) {
		f(enum)
	}
}

func (w *walker) walkMessages(messages protoreflect.MessageDescriptors, f WalkFunc) {
	for i := 0; i < messages.Len(); i++ {
		w.walkMessage(messages.Get(i), f)
	}
}

func (w *walker) walkMessage(message protoreflect.MessageDescriptor, f WalkFunc) {
	if w.enter(string(message.FullName())) {
		if !f(message) {
			return
		}
		w.walkFields(message.Fields(), f)
		w.walkMessages(message.Messages(), f)
		w.walkEnums(message.Enums(), f)
	}
}

func (w *walker) walkFields(fields protoreflect.FieldDescriptors, f WalkFunc) {
	for i := 0; i < fields.Len(); i++ {
		w.walkField(fields.Get(i), f)
	}
}

func (w *walker) walkField(field protoreflect.FieldDescriptor, f WalkFunc) {
	if w.enter(string(field.FullName())) {
		if !f(field) {
			return
		}
		if field.IsMap() {
			w.walkField(field.MapKey(), f)
			w.walkField(field.MapValue(), f)
		}
		if field.Message() != nil {
			w.walkMessage(field.Message(), f)
		}
		if field.Enum() != nil {
			w.walkEnum(field.Enum(), f)
		}
	}
}

func (w *walker) walkServices(services protoreflect.ServiceDescriptors, f WalkFunc) {
	for i := 0; i < services.Len(); i++ {
		w.walkService(services.Get(i), f)
	}
}

func (w *walker) walkService(service protoreflect.ServiceDescriptor, f WalkFunc) {
	if w.enter(string(service.FullName())) {
		if !f(service) {
			return
		}
		w.walkMethods(service.Methods(), f)
	}
}

func (w *walker) walkMethods(methods protoreflect.MethodDescriptors, f WalkFunc) {
	for i := 0; i < methods.Len(); i++ {
		w.walkMethod(methods.Get(i), f)
	}
}

func (w *walker) walkMethod(method protoreflect.MethodDescriptor, f WalkFunc) {
	if w.enter(string(method.FullName())) {
		if !f(method) {
			return
		}
		w.walkMessage(method.Input(), f)
		w.walkMessage(method.Output(), f)
	}
}
