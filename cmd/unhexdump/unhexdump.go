package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var linePattern = regexp.MustCompile(`^.*[|]([ 0-9a-fA-F]+)[|].*$`)

func parseLine(line string, dst *os.File) error {

	var hexString string

	matches := linePattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil // No hex data found, skip this line
	}
	hexString = matches[1]

	hexBytes := strings.Split(hexString, " ")
	for _, hb := range hexBytes {
		if hb == "" {
			continue
		}
		b, err := strconv.ParseUint(hb, 16, 8)
		if err != nil {
			return err
		}
		if _, err := dst.Write([]byte{byte(b)}); err != nil {
			return err
		}
	}
	return nil
}

func unhexdump(srcfile string, dstfile string) error {
	// Open the source file
	src, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(dstfile)
	if err != nil {
		return err
	}
	defer dst.Close()

	reader := bufio.NewReader(src)

	for {
		// Read a line from the source file
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Parse the line and write the binary data to the destination file
		if err := parseLine(line, dst); err != nil {
			return err
		}
	}

	return nil
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source file> <destination file>\n", os.Args[0])
		os.Exit(1)
	}
	srcfile := os.Args[1]
	dstfile := os.Args[2]

	if srcfile == "-" {
		srcfile = "/dev/stdin"
	}
	if dstfile == "-" {
		dstfile = "/dev/stdout"
	}

	err := unhexdump(srcfile, dstfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
