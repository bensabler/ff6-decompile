package mesenstate

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"strings"
	"testing"
)

// block encodes one name/length/data record the way Mesen writes them.
func block(name string, data []byte) []byte {
	var b bytes.Buffer
	b.WriteString(name)
	b.WriteByte(0)
	var n [4]byte
	binary.LittleEndian.PutUint32(n[:], uint32(len(data)))
	b.Write(n[:])
	b.Write(data)
	return b.Bytes()
}

// savestate wraps sections into an MSS image: signature, a short header,
// then one zlib stream per section.
func savestate(t *testing.T, sections ...[]byte) []byte {
	t.Helper()
	var out bytes.Buffer
	out.WriteString(signature)
	out.Write([]byte{0x01, 0x02, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00})
	for _, s := range sections {
		zw := zlib.NewWriter(&out)
		if _, err := zw.Write(s); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return out.Bytes()
}

func fill(n int, seed byte) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed + byte(i%251)
	}
	return b
}

func realSection() []byte {
	var b bytes.Buffer
	b.Write(block("cpu.a", []byte{0x00, 0x00}))
	b.Write(block("cpu.k", []byte{0xC0}))
	b.Write(block(WorkRAMBlock, fill(WorkRAMSize, 1)))
	b.Write(block(SaveRAMBlock, fill(SaveRAMSize, 2)))
	return b.Bytes()
}

func TestParseRoundTrip(t *testing.T) {
	st, err := Parse(savestate(t, realSection()))
	if err != nil {
		t.Fatal(err)
	}
	wram, err := st.WorkRAM()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(wram, fill(WorkRAMSize, 1)) {
		t.Error("work RAM does not round-trip")
	}
	sram, err := st.SaveRAM()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sram, fill(SaveRAMSize, 2)) {
		t.Error("save RAM does not round-trip")
	}
	if got := st.Names(); len(got) != 4 || got[0] != "cart.saveRam" {
		t.Errorf("Names() = %v", got)
	}
	if k, ok := st.Block("cpu.k"); !ok || len(k) != 1 || k[0] != 0xC0 {
		t.Errorf("cpu.k = %v, %v", k, ok)
	}
}

// A savestate leads with a stream that is not the block section (Mesen
// stores a preview image there). Parse must skip it, not fail.
func TestParseSkipsNonBlockStream(t *testing.T) {
	decoy := bytes.Repeat([]byte{0xAA, 0xBB, 0xCC}, 4096)
	st, err := Parse(savestate(t, decoy, realSection()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.WorkRAM(); err != nil {
		t.Fatalf("work RAM not found past the decoy stream: %v", err)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "not a savestate",
			data: []byte("PNG\x00whatever"),
			want: "missing \"MSS\" signature",
		},
		{
			name: "empty input",
			data: nil,
			want: "missing \"MSS\" signature",
		},
		{
			name: "no block section",
			data: savestate(t, bytes.Repeat([]byte{0xAA}, 512)),
			want: "no named-block section",
		},
		{
			name: "length field truncated",
			data: savestate(t, append(block("cpu.a", []byte{1, 2}), "cpu.k\x00\x01"...)),
			want: "no named-block section",
		},
		{
			name: "block runs past end of section",
			data: savestate(t, block("cpu.a", make([]byte, 8))[:10]),
			want: "no named-block section",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.data)
			if err == nil {
				t.Fatal("want an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err, tt.want)
			}
		})
	}
}

// parseBlocks carries the specific diagnostics; Parse collapses them
// into "no named-block section" because it cannot tell a malformed
// section from a stream that was never a section.
func TestParseBlocksDiagnostics(t *testing.T) {
	tests := []struct {
		name    string
		section []byte
		want    string
	}{
		{
			name:    "unterminated name",
			section: []byte("cpu.a"),
			want:    "block name missing or empty",
		},
		{
			name:    "empty name",
			section: []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0xFF},
			want:    "block name missing or empty",
		},
		{
			name:    "non-printable name",
			section: append([]byte{0x01, 0x02, 0x00}, 0x00, 0x00, 0x00, 0x00),
			want:    "not printable ASCII",
		},
		{
			name:    "truncated length",
			section: []byte("cpu.a\x00\x02\x00"),
			want:    "truncated length field",
		},
		{
			name:    "size past end",
			section: []byte("cpu.a\x00\xFF\xFF\xFF\xFF"),
			want:    "past end of section",
		},
		{
			name:    "names without a dot",
			section: block("abc", []byte{1}),
			want:    "no named blocks",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseBlocks(tt.section)
			if err == nil {
				t.Fatal("want an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err, tt.want)
			}
		})
	}
}

func TestRegionSizeMismatch(t *testing.T) {
	st, err := Parse(savestate(t, block("memoryManager.workRam", fill(1024, 0))))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.WorkRAM(); err == nil || !strings.Contains(err.Error(), "want 131072") {
		t.Fatalf("want a size-mismatch error, got %v", err)
	}
	if _, err := st.SaveRAM(); err == nil || !strings.Contains(err.Error(), "no \"cart.saveRam\" block") {
		t.Fatalf("want a missing-block error, got %v", err)
	}
}
