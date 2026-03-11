package main

import (
	"strings"
	"testing"
)

func TestLuhnCheckDigit(t *testing.T) {
	cd, err := luhnCheckDigit("7992739871")
	if err != nil {
		t.Fatalf("luhnCheckDigit returned error: %v", err)
	}
	if cd != '3' {
		t.Fatalf("luhnCheckDigit returned %q, want %q", cd, '3')
	}
}

func TestGenerateLuhnPrefixAndLength(t *testing.T) {
	n, err := generateLuhn("411111", 16, newSeededRandomSource(1234))
	if err != nil {
		t.Fatalf("generateLuhn returned error: %v", err)
	}
	if len(n) != 16 {
		t.Fatalf("generated length = %d, want 16", len(n))
	}
	if !strings.HasPrefix(n, "411111") {
		t.Fatalf("generated prefix mismatch: %s", n)
	}
	if !luhnValid(n) {
		t.Fatalf("generated value is not Luhn-valid: %s", n)
	}
}

func TestGenerateLuhnErrors(t *testing.T) {
	if _, err := generateLuhn("12a", 16, newSeededRandomSource(1)); err == nil {
		t.Fatalf("expected non-digit prefix error")
	}
	if _, err := generateLuhn("123456", 6, newSeededRandomSource(1)); err == nil {
		t.Fatalf("expected prefix too long error")
	}
	if _, err := generateLuhn("", 1, newSeededRandomSource(1)); err == nil {
		t.Fatalf("expected minimum length error")
	}
}

func TestPickCardPrefix(t *testing.T) {
	prefix, length, err := pickCardPrefix("A", newSeededRandomSource(9))
	if err != nil {
		t.Fatalf("pickCardPrefix returned error: %v", err)
	}
	if length != 15 {
		t.Fatalf("amex length = %d, want 15", length)
	}
	if prefix != "34" && prefix != "37" {
		t.Fatalf("amex prefix = %s, want 34 or 37", prefix)
	}

	if _, _, err := pickCardPrefix("X", newSeededRandomSource(9)); err == nil {
		t.Fatalf("expected invalid cardtype error")
	}
}

func TestGenerateLuhnDeterministicWithSeed(t *testing.T) {
	rs1 := newSeededRandomSource(42)
	rs2 := newSeededRandomSource(42)
	n1, err := generateLuhn("411111", 16, rs1)
	if err != nil {
		t.Fatalf("generateLuhn first call error: %v", err)
	}
	n2, err := generateLuhn("411111", 16, rs2)
	if err != nil {
		t.Fatalf("generateLuhn second call error: %v", err)
	}
	if n1 != n2 {
		t.Fatalf("expected deterministic output, got %q and %q", n1, n2)
	}
}
