package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jguerreno/JSON-Converter/internal/generator"
)

type CLIConfig struct {
	InputFile  string
	OutputFile string
	Language   string
	RootName   string
}

func parseCLIFlags(args []string, stderr io.Writer) *CLIConfig {
	config := &CLIConfig{}

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&config.InputFile, "input", "", "Input JSON file (required, or use stdin)")
	fs.StringVar(&config.InputFile, "i", "", "Input JSON file (shorthand)")
	fs.StringVar(&config.OutputFile, "output", "", "Output file (optional, default: stdout)")
	fs.StringVar(&config.OutputFile, "o", "", "Output file (shorthand)")
	fs.StringVar(&config.Language, "lang", "go", "Target language: go, python, typescript, java")
	fs.StringVar(&config.Language, "l", "go", "Target language (shorthand)")
	fs.StringVar(&config.RootName, "root", "Root", "Root struct/class name")
	fs.StringVar(&config.RootName, "r", "Root", "Root name (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "JSON Code Generator - Convert JSON to Go/Python/TypeScript/Java\n\n")
		fmt.Fprintf(stderr, "Usage:\n")
		fmt.Fprintf(stderr, "  %s [options]\n\n", args[0])
		fmt.Fprintf(stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nExamples:\n")
		fmt.Fprintf(stderr, "  %s -i input.json -l go\n", args[0])
		fmt.Fprintf(stderr, "  %s -i input.json -l python -o output.py\n", args[0])
		fmt.Fprintf(stderr, "  cat input.json | %s -l typescript\n", args[0])
		fmt.Fprintf(stderr, "  %s -i input.json -l go -p models -r User\n", args[0])
	}

	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
	return config
}

func runCLI(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	config := parseCLIFlags(args, stderr)

	var jsonData []byte
	var err error

	if config.InputFile != "" {
		jsonData, err = os.ReadFile(config.InputFile)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
	} else {
		jsonData, err = io.ReadAll(stdin)
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
	}

	service := generator.NewGeneratorService()
	output, err := service.GenerateFromJSON(jsonData, config.RootName, config.Language)
	if err != nil {
		return err
	}
	if config.OutputFile != "" {
		if err := os.WriteFile(config.OutputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
		fmt.Fprintf(stderr, "Code generated successfully: %s\n", config.OutputFile)
	} else {
		fmt.Fprint(stdout, output)
	}

	return nil
}

func main() {
	if err := runCLI(os.Args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
