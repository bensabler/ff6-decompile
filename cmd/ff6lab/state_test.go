package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeState builds a minimal Mesen savestate carrying a work-RAM image,
// so the state commands can be exercised without a real .mss (savestates
// are ROM-derived and never committed).
func writeState(t *testing.T, dir, name string, wram []byte) string {
	t.Helper()
	var section bytes.Buffer
	writeBlock := func(blk string, data []byte) {
		section.WriteString(blk)
		section.WriteByte(0)
		var n [4]byte
		binary.LittleEndian.PutUint32(n[:], uint32(len(data)))
		section.Write(n[:])
		section.Write(data)
	}
	writeBlock("cpu.k", []byte{0xC0})
	writeBlock("memoryManager.workRam", wram)

	var out bytes.Buffer
	out.WriteString("MSS")
	out.Write([]byte{0x01, 0x02, 0x02, 0x00})
	zw := zlib.NewWriter(&out)
	if _, err := zw.Write(section.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, out.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestStateParseHex(t *testing.T) {
	cases := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"1609", 0x1609, false},
		{"$1609", 0x1609, false},
		{"+$1609", 0x1609, false},
		{"0x1609", 0x1609, false},
		{"1d4e", 0x1D4E, false},
		{"", 0, true},
		{"zz", 0, true},
	}
	for _, tt := range cases {
		got, err := parseHex(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseHex(%q): want error", tt.in)
			}
			continue
		}
		if err != nil || got != tt.want {
			t.Errorf("parseHex(%q) = %d, %v; want %d", tt.in, got, err, tt.want)
		}
	}
}

func TestTakeOutputFlag(t *testing.T) {
	rest, dest := takeOutputFlag([]string{"-o", "x.bin"})
	if dest != "x.bin" || len(rest) != 0 {
		t.Errorf("got rest=%v dest=%q", rest, dest)
	}
	rest, dest = takeOutputFlag([]string{"extra"})
	if dest != "" || len(rest) != 1 {
		t.Errorf("got rest=%v dest=%q", rest, dest)
	}
}

func TestStateReadAndDiff(t *testing.T) {
	dir := t.TempDir()
	a := make([]byte, 128*1024)
	copy(a[0x1609:], []byte{0x4C, 0x00, 0x4D, 0x00})
	b := make([]byte, 128*1024)
	copy(b[0x1609:], []byte{0x4C, 0x00, 0x4D, 0x00})
	b[0x1D4E] = 0x40 // the single-bit config change

	pa := writeState(t, dir, "a.mss", a)
	pb := writeState(t, dir, "b.mss", b)

	var out bytes.Buffer
	if err := stateRead(pa, "wram", "+$1609", "4", &out); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.Contains(got, "WRAM:+$1609  4c 00 4d 00") {
		t.Errorf("read output = %q", got)
	}

	out.Reset()
	if err := stateDiff(pa, pb, "wram", &out); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "WRAM:+$1D4E-WRAM:+$1D4E  A=00  B=40") {
		t.Errorf("diff output = %q", got)
	}
	if !strings.Contains(got, "1 differing bytes in 1 runs") {
		t.Errorf("diff summary = %q", got)
	}

	out.Reset()
	if err := stateDiff(pa, pa, "wram", &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "identical") {
		t.Errorf("self-diff = %q", out.String())
	}
}

func TestStateErrors(t *testing.T) {
	dir := t.TempDir()
	p := writeState(t, dir, "a.mss", make([]byte, 128*1024))
	var out bytes.Buffer

	if err := stateRead(p, "vram", "0", "1", &out); err == nil ||
		!strings.Contains(err.Error(), "unknown region") {
		t.Errorf("unknown region: got %v", err)
	}
	if err := stateRead(p, "wram", "1FFFF", "8", &out); err == nil ||
		!strings.Contains(err.Error(), "past the end") {
		t.Errorf("out of range: got %v", err)
	}
	if err := stateRead(p, "wram", "0", "0", &out); err == nil ||
		!strings.Contains(err.Error(), "positive decimal") {
		t.Errorf("bad length: got %v", err)
	}
	if err := stateRead(p, "sram", "0", "1", &out); err == nil ||
		!strings.Contains(err.Error(), "cart.saveRam") {
		t.Errorf("missing sram block: got %v", err)
	}
	if err := stateDump(p, "wram", "", &out); err == nil ||
		!strings.Contains(err.Error(), "-o") {
		t.Errorf("dump without -o: got %v", err)
	}
	if err := runState([]string{"bogus", p}, &out); err == nil ||
		!strings.Contains(err.Error(), "usage:") {
		t.Errorf("bad subcommand: got %v", err)
	}
}

func TestStateDumpWritesFile(t *testing.T) {
	dir := t.TempDir()
	wram := make([]byte, 128*1024)
	wram[0x1D4E] = 0x40
	p := writeState(t, dir, "a.mss", wram)
	dest := filepath.Join(dir, "wram.bin")

	var out bytes.Buffer
	if err := stateDump(p, "wram", dest, &out); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 128*1024 || got[0x1D4E] != 0x40 {
		t.Errorf("dumped image is wrong: len=%d byte=%02x", len(got), got[0x1D4E])
	}
}
