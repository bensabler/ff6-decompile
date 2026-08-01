package textenc

import "testing"

func TestDecodeRuns(t *testing.T) {
	tests := []struct {
		name string
		b    byte
		want rune
		ok   bool
	}{
		{"A", 0x80, 'A', true},
		{"Z", 0x99, 'Z', true},
		{"a", 0x9A, 'a', true},
		{"z", 0xB3, 'z', true},
		{"0", 0xB4, '0', true},
		{"2 (observed on-screen, EXP-0026)", 0xB6, '2', true},
		{"3 (observed on-screen, EXP-0026)", 0xB7, '3', true},
		{"9", 0xBD, '9', true},
		{"hyphen", 0xC4, '-', true},
		{"narrow space", 0xFE, ' ', true},
		{"blank padding", 0xFF, ' ', true},
		{"unmapped control", 0x00, 0, false},
		{"unmapped icon byte", 0xE8, 0, false},
		{"gap between digits and hyphen", 0xBE, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Decode(tt.b)
			if ok != tt.ok || (ok && got != tt.want) {
				t.Fatalf("Decode(%#x) = %q,%v; want %q,%v", tt.b, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestDecodeFixed(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want string
	}{
		// Spell id 45, name bytes from ROMFILE:0x26F6A4 (EXP-0027).
		{"Cure with $FF padding", []byte{0x82, 0xAE, 0xAB, 0x9E, 0xFF, 0xFF}, "Cure"},
		// Spell id 5: the $FE narrow space is a SPACE, not a period.
		// EXP-0027's local script rendered $FE as '.' for visibility;
		// EXP-0026 established the glyph itself as a narrow space.
		{"narrow space is a space", []byte{0x85, 0xA2, 0xAB, 0x9E, 0xFE, 0xB6}, "Fire 2"},
		{"interior narrow space", []byte{0x96, 0xFE, 0x96, 0xA2, 0xA7, 0x9D}, "W Wind"},
		{"unmapped byte is escaped, not dropped", []byte{0x85, 0x01, 0x9E}, `F\x01e`},
		{"all padding trims to empty", []byte{0xFF, 0xFF, 0xFF}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeFixed(tt.src); got != tt.want {
				t.Fatalf("DecodeFixed(% x) = %q, want %q", tt.src, got, tt.want)
			}
		})
	}
}

func TestClean(t *testing.T) {
	if !Clean([]byte{0x82, 0xAE, 0xFF}) {
		t.Fatal("Clean() = false for fully mapped bytes")
	}
	if Clean([]byte{0x82, 0x01}) {
		t.Fatal("Clean() = true despite an unmapped byte")
	}
}
