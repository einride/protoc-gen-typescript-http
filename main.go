package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.einride.tech/protoc-gen-typescript-http/internal/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var flags flag.FlagSet

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func run() error {
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return err
	}
	opts := plugin.Options{
		UseProtoNames:  flags.Bool("use_proto_names", false, "Uses proto field name instead of lowerCamelCase name in JSON field names"),
		UseEnumNumbers: flags.Bool("use_enum_numbers", false, "Emits enum values as numbers."),
		BodyStringify:  flags.Bool("body_stringify", true, "Stringify request body"),
	}
	resp, err := plugin.Generate(req, opts)
	if err != nil {
		return err
	}
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}
