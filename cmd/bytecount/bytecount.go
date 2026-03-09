package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/pflag"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

func outputPretty(out io.Writer, counts map[byte]int) {

	// LATER: use https://github.com/jeandeaual/go-locale to determine locale
	prettyPrinter := message.NewPrinter(language.English)
	// LATER: option to output markdown or pure ASCII (i.e. not using box-drawing characters)
	table := tablewriter.NewTable(out,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Symbols: tw.NewSymbols(tw.StyleRounded),
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.Off},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignRight}},
		}),
	)

	header := []string{""}
	for i := 0; i <= 0x0f; i += 1 {
		header = append(header, fmt.Sprintf("0x%02X", i))
	}
	header = append(header, "")
	table.Header(header)

	for row := 0; row <= 0xF0; row += 0x10 {
		data := []string{fmt.Sprintf("0x%02X", row)}
		for col := 0; col <= 0x0F; col += 0x01 {
			i := row + col
			data = append(data, prettyPrinter.Sprintf("%d", counts[byte(i)]))
		}
		data = append(data, fmt.Sprintf("0x%02X", row))
		table.Append(data)
	}

	table.Render()

}

func processFile(fileName string) error {
	counts := make(map[byte]int)
	for i := 0; i < 256; i++ {
		counts[byte(i)] = 0
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024)
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		counts[b]++
	}

	//LATER: other output formats: plain, JSON, CSV
	outputPretty(os.Stdout, counts)

	return nil
}

func main() {

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "bytecount version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: bytecount [options] [file...]\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			fmt.Printf("No input files specified and no data piped to stdin.\n\n")
			os.Exit(1)
		}
		args = append(args, "-")
	}

	for _, arg := range args {
		if arg == "-" {
			arg = "/dev/stdin"
		}

		if len(args) > 1 {
			fmt.Printf("Processing file: %s\n", arg)
		}

		processFile(arg)
	}
}
