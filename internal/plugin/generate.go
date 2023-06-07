package plugin

import (
	"fmt"
	"path"
	"strings"

	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Options controls the generated code.
type Options struct {
	// UseProtoNames controls the casing of generated field names.
	// If set to true, fields will use proto names (typically snake_case).
	// If omitted or set to false, fields will use JSON names (typically camelCase).
	UseProtoNames bool
	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool
	// DisableBodyStringify By default, this should use JSON.stringify to not be a breaking change.
	DisableBodyStringify bool
}

func Generate(request *pluginpb.CodeGeneratorRequest, opts Options) (*pluginpb.CodeGeneratorResponse, error) {
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
	packaged := make(map[protoreflect.FullName][]protoreflect.FileDescriptor)
	for _, f := range request.FileToGenerate {
		file, err := registry.FindFileByPath(f)
		if err != nil {
			return nil, fmt.Errorf("find file %s: %w", f, err)
		}
		packaged[file.Package()] = append(packaged[file.Package()], file)
	}

	var res pluginpb.CodeGeneratorResponse
	for pkg, files := range packaged {
		var index codegen.File
		indexPathElems := append(strings.Split(string(pkg), "."), "index.ts")
		if err := (packageGenerator{pkg: pkg, files: files, opts: opts}).Generate(&index); err != nil {
			return nil, fmt.Errorf("generate package '%s': %w", pkg, err)
		}
		index.P()
		index.P("// @@protoc_insertion_point(typescript-http-eof)")
		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(path.Join(indexPathElems...)),
			Content: proto.String(string(index.Content())),
		})
	}
	res.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))
	return &res, nil
}
