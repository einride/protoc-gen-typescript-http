package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.einride.tech/protoc-gen-typescript-http/internal/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func options(parameter string) plugin.Options {
	opts := plugin.Options{BodyStringify: true}
	for _, param := range strings.Split(parameter, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		switch param {
		case "use_proto_names":
			opts.UseProtoNames = value == "true"
		case "use_enum_numbers":
			opts.UseEnumNumbers = value == "true"
		case "body_stringify":
			opts.BodyStringify = value == "true"
		}
	}
	return opts
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
	resp, err := plugin.Generate(req, options(req.GetParameter()))
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
