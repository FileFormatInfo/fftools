package main

import (
	_ "embed"
	"fmt"
	"os"

	anyascii "github.com/anyascii/go"
	"github.com/spf13/pflag"
	"golang.org/x/text/encoding/charmap"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

func main() {

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "asciify version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: asciify [options] [file...]\n\n")
		fmt.Printf("%s\n", helpText)
		return
	}

	input, err := os.ReadFile("/dev/stdin")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		os.Exit(1)
	}

	// LATER: option to use a different code page
	decoder := charmap.CodePage437.NewDecoder()

	// LATER: option to skip decoding (in case input is already UTF-8)
	utf8, err := decoder.Bytes(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: unable to decode input: ", err)
		os.Exit(1)
	}

	output := anyascii.Transliterate(string(utf8))
	fmt.Print(output)
}
