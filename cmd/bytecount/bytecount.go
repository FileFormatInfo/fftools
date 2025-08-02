package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
	for i := 0; i <= 0xf0; i += 0x10 {
		header = append(header, fmt.Sprintf("0x%02X", i))
	}
	table.Header(header)

	for row := 0; row <= 0x0F; row += 0x01 {
		data := []string{fmt.Sprintf("0x%02X", row)}
		for col := 0; col <= 0xF0; col += 0x10 {
			i := row + col
			data = append(data, prettyPrinter.Sprintf("%d", counts[byte(i)]))
		}
		table.Append(data)
	}

	table.Render()

}

func main() {
	counts := make(map[byte]int)
	for i := 0; i < 256; i++ {
		counts[byte(i)] = 0
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to read input: %v\n", err)
			os.Exit(1)
		}
		counts[b]++
	}

	//LATER: other output formats: plain, JSON, CSV
	outputPretty(os.Stdout, counts)
}
