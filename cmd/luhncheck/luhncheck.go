package main

import (
	_ "embed"
	"fmt"
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

func allDigits(s string) bool {
	if s == "" {
		return false
	}
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

func correctedLuhn(number string) (string, error) {
	if len(number) < 2 {
		return "", fmt.Errorf("number must have at least 2 digits")
	}
	if !allDigits(number) {
		return "", fmt.Errorf("number must contain only digits")
	}
	cd, err := luhnCheckDigit(number[:len(number)-1])
	if err != nil {
		return "", err
	}
	return number[:len(number)-1] + string(cd), nil
}

func main() {
	var noErrorLevel = pflag.Bool("no-error-level", false, "Do not return non-zero exit code when Luhn check fails")
	var verbose = pflag.Bool("verbose", false, "If check fails, print corrected number with fixed last digit")
	var quiet = pflag.Bool("quiet", false, "Do not print PASS/FAIL")

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "luhncheck version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: luhncheck [options] <number>\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: missing number argument\n")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "WARNING: ignoring extra arguments (count=%d)\n", len(args)-1)
	}

	number := strings.TrimSpace(args[0])
	if len(number) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR: number must have at least 2 digits\n")
		os.Exit(1)
	}
	if !allDigits(number) {
		fmt.Fprintf(os.Stderr, "ERROR: number must contain only digits\n")
		os.Exit(1)
	}

	valid := luhnValid(number)
	if !*quiet {
		if valid {
			fmt.Fprint(os.Stdout, "PASS")
		} else {
			fmt.Fprint(os.Stdout, "FAIL")
		}
	}

	if !valid && *verbose {
		fixed, err := correctedLuhn(number)
		if err == nil {
			if !*quiet {
				fmt.Fprintln(os.Stdout)
			}
			fmt.Fprintf(os.Stdout, "Fixed: %s", fixed)
		}
	}

	if valid || *noErrorLevel {
		return
	}
	os.Exit(1)
}
