package main

import (
	crand "crypto/rand"
	_ "embed"
	"fmt"
	mrand "math/rand"
	"os"
	"strings"

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

type cardPreset struct {
	length   int
	prefixes []string
}

var cardPresets = map[string]cardPreset{
	"V": {length: 16, prefixes: []string{"4"}},
	"M": {length: 16, prefixes: []string{"51", "52", "53", "54", "55"}},
	"D": {length: 16, prefixes: []string{"6011"}},
	"A": {length: 15, prefixes: []string{"34", "37"}},
}

type randomSource interface {
	randomDigit() (byte, error)
	randomIndex(max int) (int, error)
}

type cryptoRandomSource struct{}

func (s *cryptoRandomSource) randomDigit() (byte, error) {
	buf := []byte{0}
	if _, err := crand.Read(buf); err != nil {
		return 0, err
	}
	return (buf[0] % 10) + '0', nil
}

func (s *cryptoRandomSource) randomIndex(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be > 0")
	}
	buf := []byte{0}
	if _, err := crand.Read(buf); err != nil {
		return 0, err
	}
	return int(buf[0]) % max, nil
}

type seededRandomSource struct {
	rng *mrand.Rand
}

func newSeededRandomSource(seed int64) *seededRandomSource {
	return &seededRandomSource{rng: mrand.New(mrand.NewSource(seed))}
}

func (s *seededRandomSource) randomDigit() (byte, error) {
	return byte(s.rng.Intn(10)) + '0', nil
}

func (s *seededRandomSource) randomIndex(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be > 0")
	}
	return s.rng.Intn(max), nil
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func luhnCheckDigit(withoutCheckDigit string) (byte, error) {
	if withoutCheckDigit == "" {
		return 0, fmt.Errorf("input is empty")
	}
	if !allDigits(withoutCheckDigit) {
		return 0, fmt.Errorf("input contains non-digit characters")
	}

	sum := 0
	double := true
	for i := len(withoutCheckDigit) - 1; i >= 0; i-- {
		d := int(withoutCheckDigit[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}

	check := (10 - (sum % 10)) % 10
	return byte('0' + check), nil
}

func luhnValid(number string) bool {
	if len(number) < 2 || !allDigits(number) {
		return false
	}
	cd, err := luhnCheckDigit(number[:len(number)-1])
	if err != nil {
		return false
	}
	return number[len(number)-1] == cd
}

func generateLuhn(prefix string, length int, rs randomSource) (string, error) {
	if length < 2 {
		return "", fmt.Errorf("length must be at least 2")
	}
	if prefix == "" {
		prefix = "0"
	}
	if !allDigits(prefix) {
		return "", fmt.Errorf("prefix must contain only digits")
	}
	if len(prefix) >= length {
		return "", fmt.Errorf("prefix length (%d) must be less than length (%d)", len(prefix), length)
	}

	var b strings.Builder
	b.Grow(length)
	b.WriteString(prefix)

	randDigits := length - len(prefix) - 1
	for i := 0; i < randDigits; i++ {
		d, err := rs.randomDigit()
		if err != nil {
			return "", err
		}
		b.WriteByte(d)
	}

	body := b.String()
	cd, err := luhnCheckDigit(body)
	if err != nil {
		return "", err
	}
	b.WriteByte(cd)

	out := b.String()
	if !luhnValid(out) {
		return "", fmt.Errorf("internal error: generated value failed Luhn validation")
	}
	return out, nil
}

func pickCardPrefix(cardType string, rs randomSource) (string, int, error) {
	preset, ok := cardPresets[cardType]
	if !ok {
		return "", 0, fmt.Errorf("invalid --cardtype %q (expected: V, M, D, A)", cardType)
	}
	idx, err := rs.randomIndex(len(preset.prefixes))
	if err != nil {
		return "", 0, err
	}
	return preset.prefixes[idx], preset.length, nil
}

func main() {
	var length = pflag.Int("length", 16, "Number of digits")
	var prefix = pflag.String("prefix", "", "Starting digits")
	var cardType = pflag.String("cardtype", "", "Card type preset: V (Visa), M (Mastercard), D (Discover), A (American Express)")
	var seed = pflag.Int64("seed", 0, "Seed for deterministic pseudo-random generation")
	var trailingNewline = pflag.Bool("trailing-newline", false, "Emit a trailing newline")

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "luhngen version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: luhngen [options]\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	if len(pflag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "WARNING: ignoring extra arguments (count=%d)\n", len(pflag.Args()))
	}

	effectivePrefix := strings.TrimSpace(*prefix)
	effectiveLength := *length

	var rs randomSource = &cryptoRandomSource{}
	if pflag.CommandLine.Lookup("seed").Changed {
		rs = newSeededRandomSource(*seed)
	}

	if *cardType != "" {
		ct := strings.ToUpper(strings.TrimSpace(*cardType))
		presetPrefix, presetLength, err := pickCardPrefix(ct, rs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		if !pflag.CommandLine.Lookup("prefix").Changed {
			effectivePrefix = presetPrefix
		}
		if !pflag.CommandLine.Lookup("length").Changed {
			effectiveLength = presetLength
		}
	}

	generated, err := generateLuhn(effectivePrefix, effectiveLength, rs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	if *trailingNewline {
		fmt.Fprintln(os.Stdout, generated)
		return
	}
	fmt.Fprint(os.Stdout, generated)
}
