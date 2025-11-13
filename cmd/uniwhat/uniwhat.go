package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FileFormatInfo/fftools/internal"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/runenames"
)

//go:embed README.md
var helpText string

func main() {

	var ascii = pflag.Bool("ascii", false, "Include ASCII characters")
	var codepoint = pflag.Bool("codepoint", true, "Print the U+XXXX codepoint")
	var line = pflag.Bool("line", true, "Print the line number")
	var offset = pflag.Bool("offset", true, "Print the offset")
	var char = pflag.Bool("char", false, "Print the character itself")

	var first = pflag.Bool("first", false, "Only print the first occurrence of each character")

	var help = pflag.Bool("help", false, "Detailed help")
	var version = pflag.Bool("version", false, "Version info")

	pflag.Parse()

	if *version {
		internal.PrintVersion("uniwhat")
		return
	}

	if *help {
		fmt.Printf("%s\n", helpText)
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		fmt.Printf("Usage: uniwhat [options] file ...\n\n")
		pflag.PrintDefaults()
		return
	}

	firstMap := make(map[rune]bool)

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

		reader := bufio.NewReaderSize(file, 1024*1024)

		var pos int = 0
		var lineNum int = 1

		// Loop to read runes one by one
		for {
			r, rsize, err := reader.ReadRune() // ReadRune returns the rune, its size in bytes, and an error
			if err != nil {
				if err == io.EOF { // End of file reached
					break
				}
				log.Fatalf("Error reading rune: %v", err)
			}
			pos += rsize

			if r == '\n' {
				lineNum++
			}

			if !*ascii && ((r >= 0x20 && r <= 0x7E) || r == 0x09 || r == 0x0A || r == 0x0D) {
				continue // Skip ASCII characters if --ascii is not set
			}

			if *first {
				if _, exists := firstMap[r]; exists {
					continue
				}
				firstMap[r] = true
			}

			name := runenames.Name(r)
			if name == "" {
				name = "<unknown>"
			}
			if *offset {
				fmt.Printf("%08x ", pos-rsize)
			}
			if *line {
				fmt.Printf("%6d ", lineNum)
			}
			if *codepoint {
				fmt.Printf("U+%04X ", r)
			}
			if *char {
				fmt.Printf("%c ", r)
			}
			fmt.Printf("%s\n", name)
		}
	}
}
