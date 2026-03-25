package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

type Config struct {
	InputPath   string
	OutputPath  string
	RootName    string
	PackageName string
}

func Parse(args []string, stderr io.Writer) (Config, error) {
	var config Config

	flagSet := flag.NewFlagSet("structify-go", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	flagSet.StringVar(&config.InputPath, "input", "", "Path to a JSON file. If omitted, JSON is read from stdin.")
	flagSet.StringVar(&config.RootName, "name", "", "Required root struct name for the generated model.")
	flagSet.StringVar(&config.OutputPath, "output", "", "Path to save generated Go code. If omitted, code is printed to stdout.")
	flagSet.StringVar(&config.PackageName, "package", "main", "Package name for the generated Go code.")

	flagSet.Usage = func() {
		fmt.Fprintln(stderr, "Convert JSON into Go struct definitions.")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  structify-go -input sample.json -name UserResponse")
		fmt.Fprintln(stderr, "  cat sample.json | structify-go -name UserResponse")
		fmt.Fprintln(stderr, "  structify-go -input sample.json -name UserResponse -output model.go")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Flags:")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		return Config{}, err
	}

	if flagSet.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected positional arguments: %s", strings.Join(flagSet.Args(), " "))
	}

	if strings.TrimSpace(config.RootName) == "" {
		return Config{}, errors.New("missing required -name flag")
	}

	return config, nil
}
