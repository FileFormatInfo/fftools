package main

import (
	"strings"
	"testing"
)

func TestStripTrailingDot(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "trailing dot", in: "example.com.", want: "example.com"},
		{name: "no trailing dot", in: "example.com", want: "example.com"},
		{name: "root dot", in: ".", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTrailingDot(tt.in)
			if got != tt.want {
				t.Fatalf("stripTrailingDot(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEnsureFQDN(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "adds trailing dot", in: "example.com", want: "example.com."},
		{name: "preserves trailing dot", in: "example.com.", want: "example.com."},
		{name: "root", in: ".", want: "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureFQDN(tt.in)
			if got != tt.want {
				t.Fatalf("ensureFQDN(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPunycodeRoundTrip(t *testing.T) {
	ascii, err := toASCIIHost("bücher.ch")
	if err != nil {
		t.Fatalf("toASCIIHost error: %v", err)
	}
	if ascii != "xn--bcher-kva.ch" {
		t.Fatalf("toASCIIHost returned %q, want %q", ascii, "xn--bcher-kva.ch")
	}

	unicodeHost, err := toUnicodeHost(ascii)
	if err != nil {
		t.Fatalf("toUnicodeHost error: %v", err)
	}
	if unicodeHost != "bücher.ch" {
		t.Fatalf("toUnicodeHost returned %q, want %q", unicodeHost, "bücher.ch")
	}
}

func TestTLDFromHost(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "simple", in: "www.example.com", want: "com"},
		{name: "with trailing dot", in: "www.example.com.", want: "com"},
		{name: "single label", in: "localhost", want: "localhost"},
		{name: "root", in: ".", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tldFromHost(tt.in)
			if got != tt.want {
				t.Fatalf("tldFromHost(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidHostname(t *testing.T) {
	max63 := strings.Repeat("a", 63)
	tooLongLabel := strings.Repeat("a", 64)
	validMaxFQDN := strings.Join([]string{
		strings.Repeat("a", 63),
		strings.Repeat("b", 63),
		strings.Repeat("c", 63),
		strings.Repeat("d", 61),
	}, ".")
	invalidMaxFQDN := strings.Join([]string{
		strings.Repeat("a", 63),
		strings.Repeat("b", 63),
		strings.Repeat("c", 63),
		strings.Repeat("d", 62),
	}, ".")

	tests := []struct {
		name  string
		host  string
		valid bool
	}{
		{name: "basic valid", host: "example.com", valid: true},
		{name: "valid with trailing dot", host: "example.com.", valid: true},
		{name: "unicode valid", host: "bücher.ch", valid: true},
		{name: "max label valid", host: max63 + ".com", valid: true},
		{name: "max fqdn valid", host: validMaxFQDN, valid: true},
		{name: "empty label invalid", host: "example..com", valid: false},
		{name: "leading hyphen invalid", host: "-example.com", valid: false},
		{name: "trailing hyphen invalid", host: "example-.com", valid: false},
		{name: "invalid char underscore", host: "exa_mple.com", valid: false},
		{name: "label too long", host: tooLongLabel + ".com", valid: false},
		{name: "fqdn too long", host: invalidMaxFQDN, valid: false},
		{name: "root invalid", host: ".", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validHostname(tt.host)
			if tt.valid && err != nil {
				t.Fatalf("validHostname(%q) unexpected error: %v", tt.host, err)
			}
			if !tt.valid && err == nil {
				t.Fatalf("validHostname(%q) expected error, got nil", tt.host)
			}
		})
	}
}

func TestValidTLD(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{name: "known icann tld", host: "example.com", want: true},
		{name: "known icann tld uppercase", host: "example.COM", want: true},
		{name: "unknown tld", host: "example.notatld", want: false},
		{name: "empty", host: ".", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validTLD(tt.host)
			if got != tt.want {
				t.Fatalf("validTLD(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestValidPublicSuffix(t *testing.T) {
	tests := []struct {
		name string
		host string
		want bool
	}{
		{name: "known suffix", host: "co.uk", want: true},
		{name: "known suffix with trailing dot", host: "co.uk.", want: true},
		{name: "regular host not suffix", host: "www.example.com", want: false},
		{name: "unknown suffix", host: "notasuffix", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validPublicSuffix(tt.host)
			if got != tt.want {
				t.Fatalf("validPublicSuffix(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}
