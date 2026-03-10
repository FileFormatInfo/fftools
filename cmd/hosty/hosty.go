package main

import (
	_ "embed"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

var (
	BUILDER = "unknown"
	COMMIT  = "(local)"
	LASTMOD = "(local)"
	VERSION = "internal"
)

//go:embed README.md
var helpText string

var labelRE = regexp.MustCompile(`^[a-z0-9-]+$`)

func stripTrailingDot(host string) string {
	if host == "." {
		return ""
	}
	return strings.TrimSuffix(host, ".")
}

func ensureFQDN(host string) string {
	h := stripTrailingDot(host)
	if h == "" {
		return "."
	}
	return h + "."
}

func toASCIIHost(host string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("hostname is empty")
	}
	return idna.Lookup.ToASCII(host)
}

func toUnicodeHost(host string) (string, error) {
	if host == "" {
		return "", fmt.Errorf("hostname is empty")
	}
	return idna.Lookup.ToUnicode(host)
}

func tldFromHost(host string) string {
	h := stripTrailingDot(host)
	if h == "" {
		return ""
	}
	parts := strings.Split(h, ".")
	return parts[len(parts)-1]
}

func validHostname(host string) error {
	h := stripTrailingDot(host)
	if h == "" {
		return fmt.Errorf("hostname is empty")
	}

	ascii, err := toASCIIHost(h)
	if err != nil {
		return fmt.Errorf("invalid hostname: %w", err)
	}

	ascii = strings.ToLower(ascii)
	if len(ascii) > 253 {
		return fmt.Errorf("hostname exceeds 253 octets")
	}

	labels := strings.Split(ascii, ".")
	for _, label := range labels {
		if label == "" {
			return fmt.Errorf("hostname has an empty label")
		}
		if len(label) > 63 {
			return fmt.Errorf("label %q exceeds 63 octets", label)
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return fmt.Errorf("label %q starts or ends with '-'", label)
		}
		if !labelRE.MatchString(label) {
			return fmt.Errorf("label %q contains invalid characters", label)
		}
	}

	return nil
}

func validTLD(host string) bool {
	tld := strings.ToLower(tldFromHost(host))
	if tld == "" {
		return false
	}
	suffix, icann := publicsuffix.PublicSuffix(tld)
	return icann && suffix == tld
}

func validPublicSuffix(host string) bool {
	h := strings.ToLower(stripTrailingDot(host))
	if h == "" {
		return false
	}
	suffix, icann := publicsuffix.PublicSuffix(h)
	return icann && suffix == h
}

func printValidationResult(ok bool, onInvalid string, original string) {
	if ok {
		os.Exit(0)
	}

	switch onInvalid {
	case "blank":
		fmt.Print("")
		os.Exit(0)
	case "original":
		fmt.Print(original)
		os.Exit(0)
	default:
		os.Exit(1)
	}
}

func main() {
	var toPunycode = pflag.Bool("to-punycode", false, "Convert hostname to punycode")
	var fromPunycode = pflag.Bool("from-punycode", false, "Convert punycode hostname to Unicode")
	var getPublicSuffix = pflag.Bool("public-suffix", false, "Output public suffix")
	var getETLD1 = pflag.Bool("etld1", false, "Output effective TLD+1 (public suffix plus one label)")
	var getTLD = pflag.Bool("tld", false, "Output top-level domain label")
	var getBare = pflag.Bool("bare", false, "Output hostname without trailing dot")
	var getFQDN = pflag.Bool("fqdn", false, "Output fully qualified hostname with trailing dot")

	var checkHostname = pflag.Bool("check-host", false, "Validate hostname syntax and lengths")
	var checkResolve = pflag.Bool("check-resolve", false, "Validate hostname by DNS lookup")
	var checkTLD = pflag.Bool("check-tld", false, "Validate the TLD against ICANN public suffix data")
	var checkSuffix = pflag.Bool("check-suffix", false, "Validate that input is a known public suffix")
	var onInvalid = pflag.String("on-invalid", "exit", "Validation failure behavior: exit, blank, original")

	var help = pflag.BoolP("help", "h", false, "Show help message")
	var version = pflag.Bool("version", false, "Print version information")

	pflag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "hosty version %s (built by %s on %s, commit %s)\n", VERSION, BUILDER, LASTMOD, COMMIT)
		return
	}

	if *help {
		fmt.Printf("Usage: hosty [options] <hostname>\n\n")
		pflag.PrintDefaults()
		fmt.Printf("%s\n", helpText)
		return
	}

	if *onInvalid != "exit" && *onInvalid != "blank" && *onInvalid != "original" {
		fmt.Fprintf(os.Stderr, "ERROR: --on-invalid must be one of: exit, blank, original\n")
		os.Exit(1)
	}

	actions := 0
	for _, selected := range []bool{*toPunycode, *fromPunycode, *getPublicSuffix, *getETLD1, *getTLD, *getBare, *getFQDN, *checkHostname, *checkResolve, *checkTLD, *checkSuffix} {
		if selected {
			actions++
		}
	}

	if actions == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: no action selected (try --help)\n")
		os.Exit(1)
	}
	if actions > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: select exactly one action\n")
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: missing hostname argument\n")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "WARNING: ignoring extra arguments (count=%d)\n", len(args)-1)
	}

	original := strings.TrimSpace(args[0])
	host := original

	if *toPunycode {
		ascii, err := toASCIIHost(stripTrailingDot(host))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		if strings.HasSuffix(host, ".") {
			fmt.Print(ascii + ".")
			return
		}
		fmt.Print(ascii)
		return
	}

	if *fromPunycode {
		unicodeHost, err := toUnicodeHost(stripTrailingDot(host))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		if strings.HasSuffix(host, ".") {
			fmt.Print(unicodeHost + ".")
			return
		}
		fmt.Print(unicodeHost)
		return
	}

	if *getPublicSuffix {
		h := strings.ToLower(stripTrailingDot(host))
		if h == "" {
			fmt.Fprintf(os.Stderr, "ERROR: hostname is empty\n")
			os.Exit(1)
		}
		suffix, _ := publicsuffix.PublicSuffix(h)
		fmt.Print(suffix)
		return
	}

	if *getETLD1 {
		h := strings.ToLower(stripTrailingDot(host))
		etld1, err := publicsuffix.EffectiveTLDPlusOne(h)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to derive eTLD+1: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(etld1)
		return
	}

	if *getTLD {
		fmt.Print(strings.ToLower(tldFromHost(host)))
		return
	}

	if *getBare {
		fmt.Print(stripTrailingDot(host))
		return
	}

	if *getFQDN {
		fmt.Print(ensureFQDN(host))
		return
	}

	if *checkHostname {
		printValidationResult(validHostname(host) == nil, *onInvalid, original)
	}

	if *checkResolve {
		h := stripTrailingDot(host)
		if h == "" {
			printValidationResult(false, *onInvalid, original)
		}
		_, err := net.LookupHost(h)
		printValidationResult(err == nil, *onInvalid, original)
	}

	if *checkTLD {
		printValidationResult(validTLD(host), *onInvalid, original)
	}

	if *checkSuffix {
		printValidationResult(validPublicSuffix(host), *onInvalid, original)
	}
}
