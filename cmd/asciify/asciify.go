package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	anyascii "github.com/anyascii/go"
	"github.com/mattn/go-isatty"
	"github.com/spf13/pflag"
	"github.com/wlynxg/chardet"
	"golang.org/x/text/encoding/charmap"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

var charmapAliases = map[string]*charmap.Charmap{
	"cp037":             charmap.CodePage037,
	"codepage037":       charmap.CodePage037,
	"cp437":             charmap.CodePage437,
	"codepage437":       charmap.CodePage437,
	"cp850":             charmap.CodePage850,
	"codepage850":       charmap.CodePage850,
	"cp852":             charmap.CodePage852,
	"codepage852":       charmap.CodePage852,
	"cp855":             charmap.CodePage855,
	"codepage855":       charmap.CodePage855,
	"cp858":             charmap.CodePage858,
	"codepage858":       charmap.CodePage858,
	"cp860":             charmap.CodePage860,
	"codepage860":       charmap.CodePage860,
	"cp862":             charmap.CodePage862,
	"codepage862":       charmap.CodePage862,
	"cp863":             charmap.CodePage863,
	"codepage863":       charmap.CodePage863,
	"cp865":             charmap.CodePage865,
	"codepage865":       charmap.CodePage865,
	"cp866":             charmap.CodePage866,
	"codepage866":       charmap.CodePage866,
	"cp1047":            charmap.CodePage1047,
	"codepage1047":      charmap.CodePage1047,
	"iso88591":          charmap.ISO8859_1,
	"iso88592":          charmap.ISO8859_2,
	"iso88593":          charmap.ISO8859_3,
	"iso88594":          charmap.ISO8859_4,
	"iso88595":          charmap.ISO8859_5,
	"iso88596":          charmap.ISO8859_6,
	"iso88597":          charmap.ISO8859_7,
	"iso88598":          charmap.ISO8859_8,
	"iso885910":         charmap.ISO8859_10,
	"iso885913":         charmap.ISO8859_13,
	"iso885914":         charmap.ISO8859_14,
	"iso885915":         charmap.ISO8859_15,
	"iso885916":         charmap.ISO8859_16,
	"koi8r":             charmap.KOI8R,
	"koi8u":             charmap.KOI8U,
	"macintosh":         charmap.Macintosh,
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	"windows1250":       charmap.Windows1250,
	"windows1251":       charmap.Windows1251,
	"windows1252":       charmap.Windows1252,
	"windows1253":       charmap.Windows1253,
	"windows1254":       charmap.Windows1254,
	"windows1255":       charmap.Windows1255,
	"windows1256":       charmap.Windows1256,
	"windows1257":       charmap.Windows1257,
	"windows1258":       charmap.Windows1258,
}

func normalizeEncodingName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, " ", "")
	return name
}

func resolveCharmapByName(name string) (*charmap.Charmap, bool) {
	cm, ok := charmapAliases[normalizeEncodingName(name)]
	return cm, ok
}

func decodeInput(input []byte, charmapName string) ([]byte, error) {
	mode := normalizeEncodingName(charmapName)

	if mode == "" {
		mode = "none"
	}

	if mode == "none" || mode == "utf8" || mode == "utf" || mode == "ascii" || mode == "usascii" {
		return input, nil
	}

	if mode == "auto" {
		detected := chardet.Detect(input)
		detectedMode := normalizeEncodingName(detected.Encoding)
		if detectedMode == "" {
			return nil, fmt.Errorf("unable to detect charmap")
		}
		if detectedMode == "utf8" || detectedMode == "ascii" || detectedMode == "usascii" {
			return input, nil
		}
		cm, ok := resolveCharmapByName(detectedMode)
		if !ok {
			return nil, fmt.Errorf("detected encoding %q is not a supported charmap", detected.Encoding)
		}
		return cm.NewDecoder().Bytes(input)
	}

	cm, ok := resolveCharmapByName(mode)
	if !ok {
		return nil, fmt.Errorf("unknown charmap %q (use none, auto, cp437, windows-1252, iso-8859-1, etc.)", charmapName)
	}
	return cm.NewDecoder().Bytes(input)
}

func main() {

	var charmapName = pflag.String("charmap", "utf8", "Character map for input decoding (e.g. none, auto, utf8, cp437, windows-1252, iso-8859-1)")
	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "asciify version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: asciify [options] [file...]\n\n")
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

		input, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: unable to read file:", err)
			os.Exit(1)
		}

		utf8, err := decodeInput(input, *charmapName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: unable to decode input:", err)
			os.Exit(1)
		}

		output := anyascii.Transliterate(string(utf8))
		fmt.Print(output)
	}
}
