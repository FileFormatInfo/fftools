package main

import (
	"fmt"
	"os"

	anyascii "github.com/anyascii/go"
	"golang.org/x/text/encoding/charmap"
)

func main() {

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
