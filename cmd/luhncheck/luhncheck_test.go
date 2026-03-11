package main

import "testing"

func TestLuhnCheckDigit(t *testing.T) {
	cd, err := luhnCheckDigit("7992739871")
	if err != nil {
		t.Fatalf("luhnCheckDigit error: %v", err)
	}
	if cd != '3' {
		t.Fatalf("luhnCheckDigit = %q, want %q", cd, '3')
	}
}

func TestLuhnValid(t *testing.T) {
	if !luhnValid("79927398713") {
		t.Fatalf("expected known valid number to pass")
	}
	if luhnValid("79927398714") {
		t.Fatalf("expected known invalid number to fail")
	}
}

func TestCorrectedLuhn(t *testing.T) {
	fixed, err := correctedLuhn("79927398714")
	if err != nil {
		t.Fatalf("correctedLuhn error: %v", err)
	}
	if fixed != "79927398713" {
		t.Fatalf("correctedLuhn = %s, want 79927398713", fixed)
	}
}

func TestCorrectedLuhnErrors(t *testing.T) {
	if _, err := correctedLuhn("1"); err == nil {
		t.Fatalf("expected short input error")
	}
	if _, err := correctedLuhn("12a3"); err == nil {
		t.Fatalf("expected non-digit input error")
	}
}
