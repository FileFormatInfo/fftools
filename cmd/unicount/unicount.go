package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/spf13/pflag"
	"golang.org/x/text/unicode/runenames"
)

func main() {

	var control = pflag.Bool("control", false, "Include control characters")
	var ascii = pflag.Bool("ascii", false, "Include ASCII characters")
	var codepoint = pflag.Bool("codepoint", true, "Print the U+XXXX codepoint")
	var char = pflag.Bool("char", false, "Print the character itself")
	var uname = pflag.Bool("name", true, "Print the unicode character name")

	pflag.Parse()

	args := pflag.Args()
	if len(args) == 0 {
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
	for k, _ := range runeCounts {
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
