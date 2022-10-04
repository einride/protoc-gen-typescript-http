package plugin

import (
	serviceconfigv1 "go.buf.build/protocolbuffers/go/einride/grpc-service-config/einride/serviceconfig/v1"
	"go.einride.tech/protoc-gen-typescript-http/internal/codegen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceConfigGenerator struct {
	pkg   protoreflect.FullName
	files []protoreflect.FileDescriptor
}

func (sc serviceConfigGenerator) Generate(f *codegen.File) error {
	seen := false
	for _, file := range sc.files {
		defaultServiceConfig := proto.GetExtension(
			file.Options(),
			serviceconfigv1.E_DefaultServiceConfig,
		).(*serviceconfigv1.ServiceConfig)
		if defaultServiceConfig == nil {
			continue
		}
		if seen {
			continue
		}
		seen = true

		json, err := protojson.Marshal(defaultServiceConfig)
		if err != nil {
			return err
		}

		f.P("export const defaultServiceConfiguration = ", string(json))
	}

	return nil
}
