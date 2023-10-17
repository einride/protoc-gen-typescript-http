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

type Options struct {
	// UseProtoNames controls the casing of generated field names.
	// If set to true, fields will use proto names (typically snake_case).
	// If omitted or set to false, fields will use JSON names (typically camelCase).
	UseProtoNames bool
	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool
	// The method names of service methods naming case.
	// Only work when `UseEnumNumbers=true`
	// opt:
	// camelcase: convert name to lower camel case like `camelCase`
	// pascalcase: convert name to pascalcase like `PascalCase`
	// default is pascalcase
	EnumFieldNaming string
	// Generate comments as multiline comments.
	UseMultiLineComment bool
	// force add `undefined` to message field.
	// default true
	ForceMessageFieldUndefinable bool
	// If set to true, body will be JSON.stringify before send
	// default true
	UseBodyStringify bool
	// The method names of service methods naming case.
	// opt:
	// camelcase: convert name to lower camel case like `camelCase`
	// pascalcase: convert name to pascalcase like `PascalCase`
	// default is pascalcase
	ServiceMethodNaming string
	// If set to true, field int64 and uint64 will convert to string
	ForceLongAsString bool
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
		if err := (packageGenerator{opts: opts, pkg: pkg, files: files}).Generate(&index); err != nil {
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
