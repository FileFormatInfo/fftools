package main

import (
	"bytes"
	"testing"
)

func TestResolveCharmapByName(t *testing.T) {
	if _, ok := resolveCharmapByName("cp437"); !ok {
		t.Fatalf("expected cp437 to resolve")
	}
	if _, ok := resolveCharmapByName("windows-1252"); !ok {
		t.Fatalf("expected windows-1252 to resolve")
	}
	if _, ok := resolveCharmapByName("iso_8859-1"); !ok {
		t.Fatalf("expected iso_8859-1 alias to resolve")
	}
}

func TestDecodeInputNone(t *testing.T) {
	input := []byte("héllo")
	out, err := decodeInput(input, "none")
	if err != nil {
		t.Fatalf("decodeInput none error: %v", err)
	}
	if !bytes.Equal(out, input) {
		t.Fatalf("decodeInput none changed bytes")
	}
}

func TestDecodeInputExplicitCharmap(t *testing.T) {
	out, err := decodeInput([]byte{0x82}, "cp437")
	if err != nil {
		t.Fatalf("decodeInput cp437 error: %v", err)
	}
	if string(out) != "é" {
		t.Fatalf("decodeInput cp437 = %q, want %q", string(out), "é")
	}
}

func TestDecodeInputAutoUtf8(t *testing.T) {
	input := []byte("hello world")
	out, err := decodeInput(input, "auto")
	if err != nil {
		t.Fatalf("decodeInput auto error: %v", err)
	}
	if string(out) != string(input) {
		t.Fatalf("decodeInput auto changed utf8/ascii input")
	}
}

func TestDecodeInputUnknownCharmap(t *testing.T) {
	if _, err := decodeInput([]byte("abc"), "not-a-charmap"); err == nil {
		t.Fatalf("expected unknown charmap error")
	}
}
