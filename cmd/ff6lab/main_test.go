package main

import (
	"bytes"
	"testing"
)

func TestParseHexBytes(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    []byte
		wantErr bool
	}{
		{"bridge format", "61 01 00 FF", []byte{0x61, 0x01, 0x00, 0xFF}, false},
		{"newlines and runs of spaces", " 0A\n0b  0C ", []byte{0x0A, 0x0B, 0x0C}, false},
		{"empty", "", []byte{}, false},
		{"bad token", "61 GG", nil, true},
		{"overlong token", "1FF", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseHexBytes(tc.in)
			if tc.wantErr != (err != nil) {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
			if err == nil && !bytes.Equal(got, tc.want) {
				t.Fatalf("got % X, want % X", got, tc.want)
			}
		})
	}
}
