package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

const TrueString = "true"

func NewOptions(parameter string) plugin.Options {
	opts := plugin.Options{
		UseMultiLineComment:          true,
		ForceMessageFieldUndefinable: true,
		UseBodyStringify:             true,
		ServiceMethodNaming:          "none",
		EnumFieldNaming:              "pascalcase",
	}

	for _, param := range strings.Split(parameter, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}

		switch param {
		case "use_enum_numbers":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}

			opts.UseEnumNumbers = enable
		case "use_proto_names":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}
			opts.UseProtoNames = enable
		case "use_multiline_comment":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}
			opts.UseMultiLineComment = enable
		case "force_message_field_undefined":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}
			opts.ForceMessageFieldUndefinable = enable
		case "force_long_as_string":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}
			opts.ForceLongAsString = enable
		case "use_body_stringify":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				enable = false
			}
			opts.UseBodyStringify = enable
		case "service_method_naming":
			opts.ServiceMethodNaming = value
		case "enum_field_naming":
			opts.EnumFieldNaming = value
		}
	}

	return opts
}

func run() error {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return err
	}
	opts := NewOptions(req.GetParameter())
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
