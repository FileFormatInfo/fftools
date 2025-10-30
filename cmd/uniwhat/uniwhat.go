package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FileFormatInfo/fftools/internal"
	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/runenames"
)

func main() {

	var control = pflag.Bool("control", false, "Include control characters")
	var ascii = pflag.Bool("ascii", false, "Include ASCII characters")
	var codepoint = pflag.Bool("codepoint", true, "Print the U+XXXX codepoint")
	var offset = pflag.Bool("offset", true, "Print the offset")
	var char = pflag.Bool("char", false, "Print the character itself")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		internal.PrintVersion("uniwhat")
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

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

		var pos int = 0

		// Loop to read runes one by one
		for {
			r, rsize, err := reader.ReadRune() // ReadRune returns the rune, its size in bytes, and an error
			if err != nil {
				if err == io.EOF { // End of file reached
					break
				}
				log.Fatalf("Error reading rune: %v", err)
			}
			if r < 0x1F && !*control {
				pos += rsize
				continue // Skip control characters if --control is not set
			}
			if r <= 0x7E && !*ascii {
				pos += rsize
				continue // Skip ASCII characters if --ascii is not set
			}
			name := runenames.Name(r)
			if name == "" {
				name = "<unknown>"
			}
			if *offset {
				// Note: Getting the exact byte offset of the rune is complex; this is a placeholder
				fmt.Printf("%08x ", pos)
			}
			if *codepoint {
				fmt.Printf("U+%04X ", r)
			}
			if *char {
				fmt.Printf("%c ", r)
			}
			fmt.Printf("%s\n", name)

			pos += rsize
		}
	}
}
