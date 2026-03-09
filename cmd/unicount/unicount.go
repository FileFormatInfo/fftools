package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/runenames"
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

	var control = pflag.Bool("control", false, "Include control characters")
	var ascii = pflag.Bool("ascii", false, "Include ASCII characters")
	var codepoint = pflag.Bool("codepoint", true, "Print the U+XXXX codepoint")
	var char = pflag.Bool("char", false, "Print the character itself")
	var uname = pflag.Bool("name", true, "Print the unicode character name")

	var help = pflag.Bool("help", false, "Detailed help")
	var version = pflag.Bool("version", false, "Version info")

	pflag.Parse()

	if *version {
		fmt.Printf("unicount version %s (built on %s from %s by %s)\n", VERSION, LASTMOD, COMMIT, BUILDER)
		return
	}

	if *help {
		fmt.Println("unicount - count unicode characters")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprintf(os.Stderr, "ERROR: no files specified and stdin is not piped\n\n")
			fmt.Printf("Usage: unicount [options] file ...\n\n")
			pflag.PrintDefaults()
			os.Exit(1)
		}
		args = []string{"-"}
	}

	runeCounts := make(map[rune]int)

	for _, arg := range args {
		if arg == "-" {
			arg = "/dev/stdin"
		}

		fmt.Printf("Processing file: %s\n", arg)

		file, err := os.Open(arg)
		if err != nil {
			log.Fatalf("Error opening file '%s': %v", arg, err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		for {
			r, _, err := reader.ReadRune() // ReadRune returns the rune, its size in bytes, and an error
			if err != nil {
				if err == io.EOF { // End of file reached
					break
				}
				log.Fatalf("Error reading rune: %v", err)
			}
			if r < 0x1F && !*control {
				continue // Skip control characters if --control is not set
			}
			if r <= 0x7E && !*ascii {
				continue // Skip ASCII characters if --ascii is not set
			}
			runeCounts[r]++
		}
	}

	keys := make([]rune, 0, len(runeCounts))
	for k := range runeCounts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, r := range keys {
		name := runenames.Name(r)
		if name == "" {
			name = "<unknown>"
		}
		if *codepoint {
			fmt.Printf("U+%04X ", r)
		}
		fmt.Printf("%8d", runeCounts[r])

		if *char {
			fmt.Printf(" %c", r)
		}
		if *uname {
			fmt.Printf(" %s", name)
		}
		fmt.Println()
	}
}
