package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validJSON = `{
	"id": 1,
	"username": "alice",
	"email": "alice@example.com",
	"is_admin": false,
	"profile": {
		"bio": "Software Engineer",
		"avatar_url": "https://example.com/avatar.png"
	},
	"roles": ["user", "developer"]
}`

const invalidJSON = `{"invalid": json}`

func TestParseCLIFlags_AllFlags(t *testing.T) {
	args := []string{"cmd", "-input", "input.json", "-output", "output.go", "-lang", "python", "-root", "User"}
	stderr := &bytes.Buffer{}

	config := parseCLIFlags(args, stderr)

	if config.InputFile != "input.json" {
		t.Errorf("expected InputFile=input.json, got %s", config.InputFile)
	}
	if config.OutputFile != "output.go" {
		t.Errorf("expected OutputFile=output.go, got %s", config.OutputFile)
	}
	if config.Language != "python" {
		t.Errorf("expected Language=python, got %s", config.Language)
	}
	if config.RootName != "User" {
		t.Errorf("expected RootName=User, got %s", config.RootName)
	}
}

func TestParseCLIFlags_Defaults(t *testing.T) {
	args := []string{"cmd"}
	stderr := &bytes.Buffer{}

	config := parseCLIFlags(args, stderr)
	if config.Language != "go" {
		t.Errorf("expected default Language=go, got %s", config.Language)
	}
	if config.RootName != "Root" {
		t.Errorf("expected default RootName=Root, got %s", config.RootName)
	}
}

func TestParseCLIFlags_Shorthand(t *testing.T) {
	args := []string{"cmd", "-i", "in.json", "-o", "out.py", "-l", "typescript", "-r", "MyClass"}
	stderr := &bytes.Buffer{}

	config := parseCLIFlags(args, stderr)

	if config.InputFile != "in.json" {
		t.Errorf("expected InputFile=in.json, got %s", config.InputFile)
	}
	if config.OutputFile != "out.py" {
		t.Errorf("expected OutputFile=out.py, got %s", config.OutputFile)
	}
	if config.Language != "typescript" {
		t.Errorf("expected Language=typescript, got %s", config.Language)
	}
	if config.RootName != "MyClass" {
		t.Errorf("expected RootName=MyClass, got %s", config.RootName)
	}
}

func TestRunCLI_ReadFromStdin(t *testing.T) {
	args := []string{"cmd", "-l", "go"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "type Root struct") {
		t.Error("expected Go struct in output")
	}
	if !strings.Contains(output, "type Profile struct") {
		t.Error("expected nested Profile struct in output")
	}
}

func TestRunCLI_ReadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.json")
	err := os.WriteFile(inputFile, []byte(validJSON), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	args := []string{"cmd", "-i", inputFile, "-l", "go"}
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err = runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "type Root struct") {
		t.Error("expected Go struct in output")
	}
}

func TestRunCLI_FileNotFound(t *testing.T) {
	args := []string{"cmd", "-i", "/nonexistent/file.json", "-l", "go"}
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}

	if !strings.Contains(err.Error(), "error reading file") {
		t.Errorf("expected 'error reading file' in error, got: %v", err)
	}
}

func TestRunCLI_InvalidJSON(t *testing.T) {
	args := []string{"cmd", "-l", "go"}
	stdin := strings.NewReader(invalidJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "error parsing JSON") {
		t.Errorf("expected 'error parsing JSON' in error, got: %v", err)
	}
}

func TestRunCLI_WriteToStdout(t *testing.T) {
	args := []string{"cmd", "-l", "python"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "class Root:") {
		t.Error("expected Python class in stdout")
	}
}

func TestRunCLI_WriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.go")

	args := []string{"cmd", "-o", outputFile, "-l", "go"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("output file was not created")
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "type Root struct") {
		t.Error("expected Go struct in output file")
	}

	if !strings.Contains(stderr.String(), "Code generated successfully") {
		t.Error("expected success message in stderr")
	}
}

func TestRunCLI_GoGeneration(t *testing.T) {
	args := []string{"cmd", "-l", "go", "-r", "User"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "type User struct") {
		t.Error("expected User struct with custom root name")
	}
	if !strings.Contains(output, "ID") && !strings.Contains(output, "Id") {
		t.Error("expected ID field in struct")
	}
	if !strings.Contains(output, "Username") {
		t.Error("expected Username field in struct")
	}
}

func TestRunCLI_PythonGeneration(t *testing.T) {
	args := []string{"cmd", "-l", "python", "-r", "User"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "class User:") {
		t.Error("expected User class")
	}
	if !strings.Contains(output, "class Profile:") {
		t.Error("expected Profile class")
	}
}

func TestRunCLI_TypeScriptGeneration(t *testing.T) {
	args := []string{"cmd", "-l", "typescript", "-r", "User"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "class User") {
		t.Error("expected User class")
	}
	if !strings.Contains(output, "class Profile") {
		t.Error("expected Profile class")
	}
}

func TestRunCLI_JavaGeneration(t *testing.T) {
	args := []string{"cmd", "-l", "java", "-r", "User"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("runCLI failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "class User") {
		t.Error("expected User class")
	}
	if !strings.Contains(output, "class Profile") {
		t.Error("expected Profile class")
	}
}

func TestRunCLI_UnsupportedLanguage(t *testing.T) {
	args := []string{"cmd", "-l", "ruby"}
	stdin := strings.NewReader(validJSON)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err := runCLI(args, stdin, stdout, stderr)
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("expected 'not supported' in error message, got: %v", err)
	}
}

func TestRunCLI_CustomRootName(t *testing.T) {
	testCases := []struct {
		lang     string
		expected string
	}{
		{"go", "type CustomRoot struct"},
		{"python", "class CustomRoot:"},
		{"typescript", "class CustomRoot"},
		{"java", "class CustomRoot"},
	}

	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			args := []string{"cmd", "-l", tc.lang, "-r", "CustomRoot"}
			stdin := strings.NewReader(validJSON)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			err := runCLI(args, stdin, stdout, stderr)
			if err != nil {
				t.Fatalf("runCLI failed for %s: %v", tc.lang, err)
			}

			output := stdout.String()
			if !strings.Contains(output, tc.expected) {
				t.Errorf("expected '%s' in %s output", tc.expected, tc.lang)
			}
		})
	}
}
