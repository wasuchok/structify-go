package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"structify-go/internal/cli"
	"structify-go/internal/generator"
	"structify-go/internal/parser"
	schematypes "structify-go/internal/types"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdin *os.File, stdout, stderr io.Writer) error {
	config, err := cli.Parse(args, stderr)
	if err != nil {
		return err
	}

	input, err := readInput(config.InputPath, stdin)
	if err != nil {
		return err
	}

	value, err := parser.ParseBytes(input)
	if err != nil {
		return err
	}

	inferrer := schematypes.NewInferrer()
	root, err := inferrer.InferRoot(config.RootName, value)
	if err != nil {
		return err
	}

	code, err := generator.Generate(root, generator.Options{
		PackageName: config.PackageName,
	})
	if err != nil {
		return err
	}

	if config.OutputPath != "" {
		if err := os.WriteFile(config.OutputPath, code, 0o644); err != nil {
			return fmt.Errorf("write output file %q: %w", config.OutputPath, err)
		}
		return nil
	}

	if _, err := stdout.Write(code); err != nil {
		return fmt.Errorf("write stdout: %w", err)
	}

	return nil
}

func readInput(inputPath string, stdin *os.File) ([]byte, error) {
	if inputPath != "" {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, fmt.Errorf("read input file %q: %w", inputPath, err)
		}

		if len(bytes.TrimSpace(data)) == 0 {
			return nil, fmt.Errorf("input file %q is empty", inputPath)
		}

		return data, nil
	}

	info, err := stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect stdin: %w", err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		return nil, errors.New("no input provided: pass -input <file> or pipe JSON into stdin")
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("stdin did not contain any JSON content")
	}

	return data, nil
}
