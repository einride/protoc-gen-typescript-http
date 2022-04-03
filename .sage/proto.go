package main

import (
	"context"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sgbuf"
)

type Proto sg.Namespace

func (Proto) All(ctx context.Context) error {
	sg.SerialDeps(ctx, Proto.Build)
	sg.Deps(ctx, Proto.BufLint, Proto.BufFormat, Proto.BufGenerate)
	return nil
}

func (Proto) Build(ctx context.Context) error {
	sg.Logger(ctx).Println("installing binary...")
	return sg.Command(
		ctx,
		"go",
		"build",
		"-o",
		sg.FromBinDir("protoc-gen-typescript-http"),
		".",
	).Run()
}

func (Proto) BufLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting proto files...")
	cmd := sgbuf.Command(ctx, "lint")
	cmd.Dir = sg.FromGitRoot("examples", "proto")
	return cmd.Run()
}

func (Proto) BufFormat(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting proto files...")
	cmd := sgbuf.Command(ctx, "format", "--write")
	cmd.Dir = sg.FromGitRoot("examples", "proto")
	return cmd.Run()
}

func (Proto) BufGenerate(ctx context.Context) error {
	sg.Logger(ctx).Println("generating from proto files...")
	cmd := sgbuf.Command(
		ctx,
		"generate",
		"--template",
		"buf.gen.yaml",
		"--path",
		"einride",
	)
	cmd.Dir = sg.FromGitRoot("examples", "proto")
	return cmd.Run()
}
