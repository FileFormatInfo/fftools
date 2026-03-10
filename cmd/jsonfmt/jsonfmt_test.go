package main

import (
	"strings"
	"testing"
)

func TestResolveModes(t *testing.T) {
	canonical, line, fractured, err := resolveModes(false, false, false)
	if err != nil {
		t.Fatalf("resolveModes default error: %v", err)
	}
	if canonical || line || !fractured {
		t.Fatalf("resolveModes default = canonical:%v line:%v fractured:%v", canonical, line, fractured)
	}

	if _, _, _, err := resolveModes(true, false, true); err == nil {
		t.Fatalf("resolveModes expected conflict error")
	}
}

func TestApplyEOL(t *testing.T) {
	in := "a\r\nb\rc\n"

	lf, err := applyEOL(in, "lf")
	if err != nil {
		t.Fatalf("applyEOL lf error: %v", err)
	}
	if lf != "a\nb\nc\n" {
		t.Fatalf("applyEOL lf = %q", lf)
	}

	cr, err := applyEOL(in, "cr")
	if err != nil {
		t.Fatalf("applyEOL cr error: %v", err)
	}
	if cr != "a\rb\rc\r" {
		t.Fatalf("applyEOL cr = %q", cr)
	}

	crlf, err := applyEOL(in, "crlf")
	if err != nil {
		t.Fatalf("applyEOL crlf error: %v", err)
	}
	if crlf != "a\r\nb\r\nc\r\n" {
		t.Fatalf("applyEOL crlf = %q", crlf)
	}

	if _, err := applyEOL(in, "bad"); err == nil {
		t.Fatalf("applyEOL expected error for invalid mode")
	}
}

func TestFormatJSONLine(t *testing.T) {
	out, err := formatJSON([]byte("{\n  \"b\": 2,\n  \"a\": 1\n}\n"), false, true, false)
	if err != nil {
		t.Fatalf("formatJSON line error: %v", err)
	}
	if out != "{\"b\":2,\"a\":1}" {
		t.Fatalf("line format output = %q", out)
	}
}

func TestFormatJSONCanonicalSortsKeys(t *testing.T) {
	out, err := formatJSON([]byte("{\"b\":2,\"a\":1}"), true, false, false)
	if err != nil {
		t.Fatalf("formatJSON canonical error: %v", err)
	}
	if !strings.Contains(out, "\n  \"a\": 1,") || !strings.Contains(out, "\n  \"b\": 2") {
		t.Fatalf("canonical output did not contain expected key/value lines: %q", out)
	}
	if strings.Index(out, "\"a\"") > strings.Index(out, "\"b\"") {
		t.Fatalf("canonical output did not sort keys: %q", out)
	}
}

func TestFormatJSONExpanded(t *testing.T) {
	out, err := formatJSON([]byte("{\"k\":\"v\"}"), false, false, false)
	if err != nil {
		t.Fatalf("formatJSON expanded error: %v", err)
	}
	if out != "{\n  \"k\": \"v\"\n}" {
		t.Fatalf("expanded output = %q", out)
	}
}

func TestFormatJSONFractured(t *testing.T) {
	out, err := formatJSON([]byte("{\"a\":1,\"b\":2}"), false, false, true)
	if err != nil {
		t.Fatalf("formatJSON fractured error: %v", err)
	}
	if !strings.Contains(out, "\"a\"") || !strings.Contains(out, "\"b\"") {
		t.Fatalf("fractured output missing keys: %q", out)
	}
}
