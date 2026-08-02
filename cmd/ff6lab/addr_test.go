package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseHexish(t *testing.T) {
	tests := []struct {
		in   string
		want uint64
		bad  bool
	}{
		{in: "$C46AC0", want: 0xC46AC0},
		{in: "0x046AC0", want: 0x046AC0},
		{in: "046AC0", want: 0x046AC0},
		{in: "C4:6AC0", want: 0xC46AC0},
		{in: "ROMCPU:$C46AC0", want: 0xC46AC0},
		{in: "ROMFILE:0x046AC0", want: 0x046AC0},
		{in: "  $C46AC0  ", want: 0xC46AC0},
		{in: "not-hex", bad: true},
		{in: "", bad: true},
	}
	for _, tt := range tests {
		got, err := parseHexish(tt.in)
		if tt.bad {
			if err == nil {
				t.Errorf("parseHexish(%q) = 0x%X, want an error", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseHexish(%q): %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseHexish(%q) = 0x%X, want 0x%X", tt.in, got, tt.want)
		}
	}
}

func TestAddrCmd(t *testing.T) {
	var buf bytes.Buffer
	if err := addrCmd([]string{"cpu", "$C47FB0"}, &buf); err != nil {
		t.Fatalf("addr cpu: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "ROMFILE:0x047FB0") || !strings.Contains(got, "Confirmed") {
		t.Errorf("addr cpu $C47FB0 output:\n%s", got)
	}

	buf.Reset()
	if err := addrCmd([]string{"file", "0x047FB0"}, &buf); err != nil {
		t.Fatalf("addr file: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "ROMCPU:$C47FB0") {
		t.Errorf("addr file 0x047FB0 output:\n%s", got)
	}

	// A mirror-window address must report its Confirmed alias rather than
	// leaving the reader to compute it.
	buf.Reset()
	if err := addrCmd([]string{"cpu", "$00FFEE"}, &buf); err != nil {
		t.Fatalf("addr cpu mirror: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "canonical:  ROMCPU:$C0FFEE") {
		t.Errorf("mirror address should name its canonical form:\n%s", got)
	}

	for _, args := range [][]string{
		{},
		{"cpu"},
		{"cpu", "$7E2EB5"},   // WRAM
		{"cpu", "$002140"},   // system area
		{"file", "0x300000"}, // past the image
		{"sideways", "$C0"},
	} {
		buf.Reset()
		if err := addrCmd(args, &buf); err == nil {
			t.Errorf("addr %v succeeded, want an error", args)
		}
	}
}
