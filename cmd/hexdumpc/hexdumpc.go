package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
)

func hexdump(fileName string, offset, length int64, w io.Writer) error {

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Seek to the offset
	if offset > 0 {
		if _, err := file.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	reader := bufio.NewReaderSize(file, 1024*1024)

	var count int64 = 0
	for {
		// Read a chunk of data
		buf := make([]byte, 16)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		printLine(w, offset+count, buf[:n])

		count += int64(n)
		if length > 0 && count >= length {
			break
		}
	}
	return nil
}

func printLine(w io.Writer, offset int64, buf []byte) {
	// Print the offset
	fmt.Fprintf(w, "%08x  ", offset)

	// Print the hex values
	for i := 0; i < len(buf); i++ {
		if i%8 == 0 && i != 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%02x ", buf[i])
	}
	// Pad the line to 48 characters
	for i := len(buf); i < 16; i++ {
		if i%8 == 0 && i != 0 {
			fmt.Fprintf(w, " ")
		}
		fmt.Fprint(w, "   ")
	}

	// Print the ASCII values
	fmt.Fprint(w, " |")
	for _, b := range buf {
		if b >= 32 && b <= 126 {
			fmt.Fprintf(w, "%c", b)
		} else {
			fmt.Fprint(w, ".")
		}
	}
	// Pad the line to 16 characters
	for i := len(buf); i < 16; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintln(w, "|")
}

func main() {

	var head = pflag.Int64("head", 0, "number of bytes to read at the beginning of the file")
	var offset = pflag.Int64("offset", 0, "number of bytes to skip at the beginning of the file")
	var length = pflag.Int64("length", 0, "number of bytes to read from the file")
	//LATER: support for tail

	pflag.Parse()

	if *head > 0 {
		*offset = 0
		*length = *head
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, arg := range args {
		if arg == "-" {
			arg = "/dev/stdin"
		}

		if err := hexdump(arg, *offset, *length, bufio.NewWriter(os.Stdout)); err != nil {
			panic(err)
		}
	}
}
