package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/sgtool"
)

const (
	name               = "eslint"
	packageJSONContent = `{
  "dependencies": {
    "@einride/eslint-plugin": "4.2.0",
    "eslint": "8.5.0"
  }
}`
)

func eslintCommand(ctx context.Context, args ...string) *exec.Cmd {
	sg.Deps(ctx, prepareEslintCommand)
	// eslint plugins should be resolved from the tool dir
	defaultArgs := []string{
		"--resolve-plugins-relative-to",
		sg.FromToolsDir(name),
	}
	cmd := sg.Command(ctx, sg.FromBinDir(name), append(defaultArgs, args...)...)
	return cmd
}

func prepareEslintCommand(ctx context.Context) error {
	toolDir := sg.FromToolsDir(name)
	binary := filepath.Join(toolDir, "node_modules", ".bin", name)
	packageJSON := filepath.Join(toolDir, "package.json")
	if err := os.MkdirAll(toolDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(packageJSON, []byte(packageJSONContent), 0o600); err != nil {
		return err
	}
	sg.Logger(ctx).Println("installing packages...")
	if err := sg.Command(
		ctx,
		"npm",
		"--silent",
		"install",
		"--prefix",
		toolDir,
		"--no-save",
		"--no-audit",
		"--ignore-script",
	).Run(); err != nil {
		return err
	}
	if _, err := sgtool.CreateSymlink(binary); err != nil {
		return err
	}
	return nil
}
