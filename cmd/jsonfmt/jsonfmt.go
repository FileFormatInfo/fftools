package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	fracturedjson "github.com/FileFormatInfo/go-fractured-json"
	"github.com/spf13/pflag"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

func normalizeLF(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func applyEOL(s string, eol string) (string, error) {
	normalized := normalizeLF(s)
	switch eol {
	case "lf":
		return normalized, nil
	case "cr":
		return strings.ReplaceAll(normalized, "\n", "\r"), nil
	case "crlf":
		return strings.ReplaceAll(normalized, "\n", "\r\n"), nil
	default:
		return "", fmt.Errorf("invalid --eol value %q (expected: lf, cr, crlf)", eol)
	}
}

func decodeJSON(input []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(input))
	dec.UseNumber()

	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("extra JSON content after first value")
	}
	return v, nil
}

func encodeSortedJSON(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			li := strings.ToLower(keys[i])
			lj := strings.ToLower(keys[j])
			if li == lj {
				return keys[i] < keys[j]
			}
			return li < lj
		})

		var buf bytes.Buffer
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			keyJSON, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			buf.Write(keyJSON)
			buf.WriteByte(':')
			valJSON, err := encodeSortedJSON(t[k])
			if err != nil {
				return nil, err
			}
			buf.Write(valJSON)
		}
		buf.WriteByte('}')
		return buf.Bytes(), nil
	case []any:
		var buf bytes.Buffer
		buf.WriteByte('[')
		for i, item := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			itemJSON, err := encodeSortedJSON(item)
			if err != nil {
				return nil, err
			}
			buf.Write(itemJSON)
		}
		buf.WriteByte(']')
		return buf.Bytes(), nil
	default:
		return json.Marshal(t)
	}
}

func maybeSortKeys(input []byte, sortKeys bool) ([]byte, error) {
	if !sortKeys {
		return input, nil
	}
	decoded, err := decodeJSON(input)
	if err != nil {
		return nil, err
	}
	return encodeSortedJSON(decoded)
}

type fracturedFlagValues struct {
	jsonEolStyle               string
	maxTotalLineLength         int
	maxInlineComplexity        int
	maxCompactArrayComplexity  int
	maxTableRowComplexity      int
	maxPropNamePadding         int
	colonBeforePropNamePadding bool
	tableCommaPlacement        string
	minCompactArrayRowItems    int
	alwaysExpandDepth          int
	nestedBracketPadding       bool
	simpleBracketPadding       bool
	colonPadding               bool
	commaPadding               bool
	commentPadding             bool
	numberListAlignment        string
	indentSpaces               int
	useTabToIndent             bool
	prefixString               string
	commentPolicy              string
	preserveBlankLines         bool
	allowTrailingCommas        bool
}

func parseEolStyle(v string) (fracturedjson.EolStyle, error) {
	switch strings.ToLower(v) {
	case "default":
		return fracturedjson.EolDefault, nil
	case "crlf":
		return fracturedjson.EolCRLF, nil
	case "lf":
		return fracturedjson.EolLF, nil
	default:
		return fracturedjson.EolDefault, fmt.Errorf("invalid fractured json eol style %q (expected: default, crlf, lf)", v)
	}
}

func parseTableCommaPlacement(v string) (fracturedjson.TableCommaPlacement, error) {
	switch strings.ToLower(v) {
	case "before-padding":
		return fracturedjson.CommaBeforePadding, nil
	case "after-padding":
		return fracturedjson.CommaAfterPadding, nil
	case "before-padding-except-numbers":
		return fracturedjson.CommaBeforePaddingExceptNumbers, nil
	default:
		return fracturedjson.CommaBeforePaddingExceptNumbers, fmt.Errorf("invalid fractured table comma placement %q (expected: before-padding, after-padding, before-padding-except-numbers)", v)
	}
}

func parseNumberListAlignment(v string) (fracturedjson.NumberListAlignment, error) {
	switch strings.ToLower(v) {
	case "left":
		return fracturedjson.NumberLeft, nil
	case "right":
		return fracturedjson.NumberRight, nil
	case "decimal":
		return fracturedjson.NumberDecimal, nil
	case "normalize":
		return fracturedjson.NumberNormalize, nil
	default:
		return fracturedjson.NumberDecimal, fmt.Errorf("invalid fractured number list alignment %q (expected: left, right, decimal, normalize)", v)
	}
}

func parseCommentPolicy(v string) (fracturedjson.CommentPolicy, error) {
	switch strings.ToLower(v) {
	case "error":
		return fracturedjson.CommentTreatAsError, nil
	case "remove":
		return fracturedjson.CommentRemove, nil
	case "preserve":
		return fracturedjson.CommentPreserve, nil
	default:
		return fracturedjson.CommentTreatAsError, fmt.Errorf("invalid fractured comment policy %q (expected: error, remove, preserve)", v)
	}
}

func buildFracturedOptions(flags fracturedFlagValues) (fracturedjson.Options, error) {
	options := fracturedjson.RecommendedOptions()

	eolStyle, err := parseEolStyle(flags.jsonEolStyle)
	if err != nil {
		return options, err
	}
	commaPlacement, err := parseTableCommaPlacement(flags.tableCommaPlacement)
	if err != nil {
		return options, err
	}
	numberAlignment, err := parseNumberListAlignment(flags.numberListAlignment)
	if err != nil {
		return options, err
	}
	commentPolicy, err := parseCommentPolicy(flags.commentPolicy)
	if err != nil {
		return options, err
	}

	options.JsonEolStyle = eolStyle
	options.MaxTotalLineLength = flags.maxTotalLineLength
	options.MaxInlineComplexity = flags.maxInlineComplexity
	options.MaxCompactArrayComplexity = flags.maxCompactArrayComplexity
	options.MaxTableRowComplexity = flags.maxTableRowComplexity
	options.MaxPropNamePadding = flags.maxPropNamePadding
	options.ColonBeforePropNamePadding = flags.colonBeforePropNamePadding
	options.TableCommaPlacement = commaPlacement
	options.MinCompactArrayRowItems = flags.minCompactArrayRowItems
	options.AlwaysExpandDepth = flags.alwaysExpandDepth
	options.NestedBracketPadding = flags.nestedBracketPadding
	options.SimpleBracketPadding = flags.simpleBracketPadding
	options.ColonPadding = flags.colonPadding
	options.CommaPadding = flags.commaPadding
	options.CommentPadding = flags.commentPadding
	options.NumberListAlignment = numberAlignment
	options.IndentSpaces = flags.indentSpaces
	options.UseTabToIndent = flags.useTabToIndent
	options.PrefixString = strings.ReplaceAll(flags.prefixString, `\t`, "\t")
	options.CommentPolicy = commentPolicy
	options.PreserveBlankLines = flags.preserveBlankLines
	options.AllowTrailingCommas = flags.allowTrailingCommas

	return options, nil
}

func formatJSON(input []byte, canonical bool, line bool, fractured bool, sortKeys bool, fracturedOptions fracturedjson.Options) (string, error) {
	input, err := maybeSortKeys(input, sortKeys)
	if err != nil {
		return "", err
	}

	if fractured {
		f := fracturedjson.NewFormatter()
		f.Options = fracturedOptions
		return f.Reformat(string(input), 0)
	}

	if line {
		var out bytes.Buffer
		if err := json.Compact(&out, input); err != nil {
			return "", err
		}
		return out.String(), nil
	}

	if canonical {
		decoded, err := decodeJSON(input)
		if err != nil {
			return "", err
		}
		out, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	var out bytes.Buffer
	if err := json.Indent(&out, input, "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

func readInput(arg string) ([]byte, error) {
	if arg == "-" {
		return os.ReadFile("/dev/stdin")
	}
	return os.ReadFile(arg)
}

func resolveModes(canonical bool, line bool, fractured bool) (bool, bool, bool, error) {
	modeCount := 0
	for _, v := range []bool{canonical, line, fractured} {
		if v {
			modeCount++
		}
	}
	if modeCount > 1 {
		return false, false, false, fmt.Errorf("use at most one of --canonical, --line, --fractured")
	}
	if modeCount == 0 {
		fractured = true
	}
	return canonical, line, fractured, nil
}

func main() {
	var canonical = pflag.Bool("canonical", false, "Output canonical JSON (same as jq . --sort-keys)")
	var line = pflag.Bool("line", false, "Output JSON on a single line")
	var fractured = pflag.Bool("fractured", false, "Use fractured JSON formatting")
	var sortKeys = pflag.Bool("sort-keys", false, "Sort object keys case-insensitively")
	var trailingNewline = pflag.Bool("trailing-newline", false, "Emit a trailing newline")
	var eol = pflag.String("eol", "lf", "End-of-line style: lf, cr, or crlf")

	defaults := fracturedjson.RecommendedOptions()
	var fracturedJsonEolStyle = pflag.String("fractured-json-eol-style", "default", "Fractured JSON EOL style: default, crlf, lf")
	var fracturedMaxTotalLineLength = pflag.Int("fractured-max-total-line-length", defaults.MaxTotalLineLength, "Fractured max total line length")
	var fracturedMaxInlineComplexity = pflag.Int("fractured-max-inline-complexity", defaults.MaxInlineComplexity, "Fractured max inline complexity")
	var fracturedMaxCompactArrayComplexity = pflag.Int("fractured-max-compact-array-complexity", defaults.MaxCompactArrayComplexity, "Fractured max compact array complexity")
	var fracturedMaxTableRowComplexity = pflag.Int("fractured-max-table-row-complexity", defaults.MaxTableRowComplexity, "Fractured max table row complexity")
	var fracturedMaxPropNamePadding = pflag.Int("fractured-max-prop-name-padding", defaults.MaxPropNamePadding, "Fractured max property name padding")
	var fracturedColonBeforePropNamePadding = pflag.Bool("fractured-colon-before-prop-name-padding", defaults.ColonBeforePropNamePadding, "Fractured: place colon before prop-name padding")
	var fracturedTableCommaPlacement = pflag.String("fractured-table-comma-placement", "before-padding-except-numbers", "Fractured table comma placement: before-padding, after-padding, before-padding-except-numbers")
	var fracturedMinCompactArrayRowItems = pflag.Int("fractured-min-compact-array-row-items", defaults.MinCompactArrayRowItems, "Fractured min compact array row items")
	var fracturedAlwaysExpandDepth = pflag.Int("fractured-always-expand-depth", defaults.AlwaysExpandDepth, "Fractured always-expand depth")
	var fracturedNestedBracketPadding = pflag.Bool("fractured-nested-bracket-padding", defaults.NestedBracketPadding, "Fractured nested bracket padding")
	var fracturedSimpleBracketPadding = pflag.Bool("fractured-simple-bracket-padding", defaults.SimpleBracketPadding, "Fractured simple bracket padding")
	var fracturedColonPadding = pflag.Bool("fractured-colon-padding", defaults.ColonPadding, "Fractured colon padding")
	var fracturedCommaPadding = pflag.Bool("fractured-comma-padding", defaults.CommaPadding, "Fractured comma padding")
	var fracturedCommentPadding = pflag.Bool("fractured-comment-padding", defaults.CommentPadding, "Fractured comment padding")
	var fracturedNumberListAlignment = pflag.String("fractured-number-list-alignment", "decimal", "Fractured number list alignment: left, right, decimal, normalize")
	var fracturedIndentSpaces = pflag.Int("fractured-indent-spaces", defaults.IndentSpaces, "Fractured indent spaces")
	var fracturedUseTabToIndent = pflag.Bool("fractured-use-tab-to-indent", defaults.UseTabToIndent, "Fractured use tabs to indent")
	var fracturedPrefixString = pflag.String("fractured-prefix-string", defaults.PrefixString, "Fractured prefix string (use \\t for tabs)")
	var fracturedCommentPolicy = pflag.String("fractured-comment-policy", "error", "Fractured comment policy: error, remove, preserve")
	var fracturedPreserveBlankLines = pflag.Bool("fractured-preserve-blank-lines", defaults.PreserveBlankLines, "Fractured preserve blank lines")
	var fracturedAllowTrailingCommas = pflag.Bool("fractured-allow-trailing-commas", defaults.AllowTrailingCommas, "Fractured allow trailing commas")

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "jsonfmt version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: jsonfmt [options] [file|-]\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	canonicalMode, lineMode, fracturedMode, err := resolveModes(*canonical, *line, *fractured)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	fracturedOptions, err := buildFracturedOptions(fracturedFlagValues{
		jsonEolStyle:               *fracturedJsonEolStyle,
		maxTotalLineLength:         *fracturedMaxTotalLineLength,
		maxInlineComplexity:        *fracturedMaxInlineComplexity,
		maxCompactArrayComplexity:  *fracturedMaxCompactArrayComplexity,
		maxTableRowComplexity:      *fracturedMaxTableRowComplexity,
		maxPropNamePadding:         *fracturedMaxPropNamePadding,
		colonBeforePropNamePadding: *fracturedColonBeforePropNamePadding,
		tableCommaPlacement:        *fracturedTableCommaPlacement,
		minCompactArrayRowItems:    *fracturedMinCompactArrayRowItems,
		alwaysExpandDepth:          *fracturedAlwaysExpandDepth,
		nestedBracketPadding:       *fracturedNestedBracketPadding,
		simpleBracketPadding:       *fracturedSimpleBracketPadding,
		colonPadding:               *fracturedColonPadding,
		commaPadding:               *fracturedCommaPadding,
		commentPadding:             *fracturedCommentPadding,
		numberListAlignment:        *fracturedNumberListAlignment,
		indentSpaces:               *fracturedIndentSpaces,
		useTabToIndent:             *fracturedUseTabToIndent,
		prefixString:               *fracturedPrefixString,
		commentPolicy:              *fracturedCommentPolicy,
		preserveBlankLines:         *fracturedPreserveBlankLines,
		allowTrailingCommas:        *fracturedAllowTrailingCommas,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "WARNING: ignoring extra arguments (count=%d)\n", len(args)-1)
	}
	arg := args[0]

	input, err := readInput(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to read input: %v\n", err)
		os.Exit(1)
	}

	formatted, err := formatJSON(input, canonicalMode, lineMode, fracturedMode, *sortKeys, fracturedOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to format JSON: %v\n", err)
		os.Exit(1)
	}

	formatted = strings.TrimRight(normalizeLF(formatted), "\n")
	if *trailingNewline {
		formatted += "\n"
	}
	formatted, err = applyEOL(formatted, strings.ToLower(*eol))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, formatted)
}
