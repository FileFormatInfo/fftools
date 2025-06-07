package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/text/encoding/charmap"
)

var asciiMap = map[int]string{
	0x00: "NUL", 0x01: "SOH", 0x02: "STX", 0x03: "ETX",
	0x04: "EOT", 0x05: "ENQ", 0x06: "ACK", 0x07: "BEL",
	0x08: "BS", 0x09: "HT", 0x0A: "LF", 0x0B: "VT",
	0x0C: "FF", 0x0D: "CR", 0x0E: "SO", 0x0F: "SI",
	0x10: "DLE", 0x11: "DC1", 0x12: "DC2", 0x13: "DC3",
	0x14: "DC4", 0x15: "NAK", 0x16: "SYN", 0x17: "ETB",
	0x18: "CAN", 0x19: "EM", 0x1A: "SUB", 0x1B: "ESC",
	0x1C: "FS", 0x1D: "GS", 0x1E: "RS", 0x1F: "US",
	0x7F: "DEL"}

func main() {

	// LATER: option to output markdown or pure ASCII (i.e. not using box-drawing characters)
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Symbols: tw.NewSymbols(tw.StyleRounded),
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.Off},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
		}),
	)

	// LATER: option to use a different code page
	decoder := charmap.CodePage437.NewDecoder()

	header := []string{""}
	for i := 0; i <= 0xf0; i += 0x10 {
		header = append(header, fmt.Sprintf("0x%02X", i))
	}
	table.Header(header)

	for row := 0; row <= 0x0F; row += 0x01 {
		data := []string{fmt.Sprintf("0x%02X", row)}
		for col := 0; col <= 0xF0; col += 0x10 {
			i := row + col
			if i < 0x20 || i == 0x7F {
				data = append(data, asciiMap[i])
			} else if i > 0x7F {
				utf8, err := decoder.Bytes([]byte{byte(i)})
				if err != nil {
					data = append(data, fmt.Sprintf("0x%02x", i))
				} else {
					data = append(data, string(utf8))
				}
			} else {
				data = append(data, string(rune(i)))
			}
		}
		table.Append(data)
	}

	table.Render()
}
