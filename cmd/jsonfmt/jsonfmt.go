package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	fracturedjson "github.com/FileFormatInfo/go-fractured-json"
	"github.com/spf13/pflag"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

func normalizeLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func applyEOL(s string, eol string) (string, error) {
	normalized := normalizeLF(s)
	switch eol {
	case "lf":
		return normalized, nil
	case "cr":
		return strings.ReplaceAll(normalized, "\n", "\r"), nil
	case "crlf":
		return strings.ReplaceAll(normalized, "\n", "\r\n"), nil
	default:
		return "", fmt.Errorf("invalid --eol value %q (expected: lf, cr, crlf)", eol)
	}
}

func decodeJSON(input []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("extra JSON content after first value")
	}
	return v, nil
}

func formatJSON(input []byte, canonical bool, line bool, fractured bool) (string, error) {
	if fractured {
		return fracturedjson.Reformat(string(input))
	}

	if line {
		var out bytes.Buffer
		if err := json.Compact(&out, input); err != nil {
			return "", err
		}
		return out.String(), nil
	}

	if canonical {
		decoded, err := decodeJSON(input)
		if err != nil {
			return "", err
		}
		out, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	var out bytes.Buffer
	if err := json.Indent(&out, input, "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

func readInput(arg string) ([]byte, error) {
	if arg == "-" {
		return os.ReadFile("/dev/stdin")
	}
	return os.ReadFile(arg)
}

func resolveModes(canonical bool, line bool, fractured bool) (bool, bool, bool, error) {
	modeCount := 0
	for _, v := range []bool{canonical, line, fractured} {
		if v {
			modeCount++
		}
	}
	if modeCount > 1 {
		return false, false, false, fmt.Errorf("use at most one of --canonical, --line, --fractured")
	}
	if modeCount == 0 {
		fractured = true
	}
	return canonical, line, fractured, nil
}

func main() {
	var canonical = pflag.Bool("canonical", false, "Output canonical JSON (same as jq . --sort-keys)")
	var line = pflag.Bool("line", false, "Output JSON on a single line")
	var fractured = pflag.Bool("fractured", false, "Use fractured JSON formatting")
	var trailingNewline = pflag.Bool("trailing-newline", false, "Emit a trailing newline")
	var eol = pflag.String("eol", "lf", "End-of-line style: lf, cr, or crlf")

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "jsonfmt version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: jsonfmt [options] [file|-]\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	canonicalMode, lineMode, fracturedMode, err := resolveModes(*canonical, *line, *fractured)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "WARNING: ignoring extra arguments (count=%d)\n", len(args)-1)
	}
	arg := args[0]

	input, err := readInput(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to read input: %v\n", err)
		os.Exit(1)
	}

	formatted, err := formatJSON(input, canonicalMode, lineMode, fracturedMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to format JSON: %v\n", err)
		os.Exit(1)
	}

	formatted = strings.TrimRight(normalizeLF(formatted), "\n")
	if *trailingNewline {
		formatted += "\n"
	}
	formatted, err = applyEOL(formatted, strings.ToLower(*eol))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, formatted)
}
