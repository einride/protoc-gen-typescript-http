package main

import (
	"context"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/tools/sgconvco"
	"go.einride.tech/sage/tools/sggit"
	"go.einride.tech/sage/tools/sggo"
	"go.einride.tech/sage/tools/sggolangcilint"
	"go.einride.tech/sage/tools/sgmdformat"
	"go.einride.tech/sage/tools/sgyamlfmt"
)

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
		sg.Makefile{
			Namespace:     Proto{},
			DefaultTarget: Proto.All,
			Path:          sg.FromGitRoot("examples", "proto", "Makefile"),
		},
	)
}

func All(ctx context.Context) error {
	sg.Deps(ctx, ConvcoCheck, GoLint, GoTest, FormatMarkdown, FormatYAML)
	sg.SerialDeps(ctx, Proto.All, TypescriptLint)
	sg.SerialDeps(ctx, GoModTidy, GitVerifyNoDiff)
	return nil
}

func FormatYAML(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting YAML files...")
	return sgyamlfmt.Run(ctx)
}

func GoModTidy(ctx context.Context) error {
	sg.Logger(ctx).Println("tidying Go module files...")
	return sg.Command(ctx, "go", "mod", "tidy", "-v").Run()
}

func GoTest(ctx context.Context) error {
	sg.Logger(ctx).Println("running Go tests...")
	return sggo.TestCommand(ctx).Run()
}

func GoLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting Go files...")
	return sggolangcilint.Run(ctx)
}

func GoLintFix(ctx context.Context) error {
	sg.Logger(ctx).Println("fixing Go files...")
	return sggolangcilint.Fix(ctx)
}

func FormatMarkdown(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Markdown files...")
	return sgmdformat.Command(ctx).Run()
}

func ConvcoCheck(ctx context.Context) error {
	sg.Logger(ctx).Println("checking git commits...")
	return sgconvco.Command(ctx, "check", "origin/master..HEAD").Run()
}

func GitVerifyNoDiff(ctx context.Context) error {
	sg.Logger(ctx).Println("verifying that git has no diff...")
	return sggit.VerifyNoDiff(ctx)
}

func TypescriptLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting typescript files...")
	return eslintCommand(
		ctx,
		"--config",
		sg.FromGitRoot(".eslintrc.js"),
		"--quiet",
		"examples/proto/gen/typescript/**/*.ts",
	).Run()
}
