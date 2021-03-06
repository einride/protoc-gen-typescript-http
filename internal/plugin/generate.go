package plugin

import (
	"fmt"
	"path"
	"strings"

	"github.com/einride/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	generate := make(map[string]struct{})
	registry, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: request.ProtoFile,
	})
	if err != nil {
		return nil, fmt.Errorf("create proto registry: %w", err)
	}
	for _, f := range request.FileToGenerate {
		generate[f] = struct{}{}
	}
	packaged := make(map[string][]protoreflect.FileDescriptor)
	for _, f := range request.FileToGenerate {
		file, err := registry.FindFileByPath(f)
		if err != nil {
			return nil, fmt.Errorf("find file %s: %w", f, err)
		}
		packaged[string(file.Package())] = append(packaged[string(file.Package())], file)
	}

	var res pluginpb.CodeGeneratorResponse
	for pkg, files := range packaged {
		var index codegen.File
		indexPathElems := append(strings.Split(pkg, "."), "index.ts")
		generatePackage(&index, files)
		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(path.Join(indexPathElems...)),
			Content: proto.String(string(index.Content())),
		})
	}
	return &res, nil
}

func generatePackage(f *codegen.File, files []protoreflect.FileDescriptor) {
	for _, file := range files {
		messages := file.Messages()
		for i := 0; i < messages.Len(); i++ {
			message := messages.Get(i)
			messageGenerator{message: message}.Generate(f)
		}
	}
}
