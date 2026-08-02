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

func TestEncodeRoundTrip(t *testing.T) {
	// Every byte Decode maps must survive Decode -> Encode -> Decode. The
	// one intentional asymmetry is the pair $FE/$FF, which both decode to
	// ' ' and re-encode to $FF.
	for i := 0; i < 256; i++ {
		b := byte(i)
		r, ok := Decode(b)
		if !ok {
			continue
		}
		got, ok := Encode(r)
		if !ok {
			t.Errorf("Decode($%02X) = %q, but Encode(%q) reports no glyph", b, r, r)
			continue
		}
		back, ok := Decode(got)
		if !ok || back != r {
			t.Errorf("$%02X -> %q -> $%02X -> %q, want a stable round trip", b, r, got, back)
		}
		if b != got && r != ' ' {
			t.Errorf("Encode(%q) = $%02X, want $%02X", r, got, b)
		}
	}
}

func TestEncodeRejectsUnverifiedGlyphs(t *testing.T) {
	// 53 bytes in the HUD block are non-blank but unidentified (EXP-0049).
	// Guessing any of them would put invented data into the reconstruction.
	for _, r := range []rune{'!', '.', ',', ':', '/', '@', 'é', '漢', 0} {
		if b, ok := Encode(r); ok {
			t.Errorf("Encode(%q) = $%02X, but that glyph value is not verified", r, b)
		}
	}
}

func TestEncodeFixed(t *testing.T) {
	got, err := EncodeFixed("WEDGE", 6)
	if err != nil {
		t.Fatalf("EncodeFixed: %v", err)
	}
	want := []byte{0x96, 0x84, 0x83, 0x86, 0x84, 0xFF}
	if len(got) != len(want) {
		t.Fatalf("EncodeFixed(\"WEDGE\", 6) = %X, want %X", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("EncodeFixed(\"WEDGE\", 6) = %X, want %X", got, want)
		}
	}
	if s := DecodeFixed(got); s != "WEDGE" {
		t.Errorf("round trip through the name field gave %q", s)
	}

	if _, err := EncodeFixed("TOOLONG", 4); err == nil {
		t.Error("overlong text should be an error, not a truncation")
	}
	if _, err := EncodeFixed("no!", 8); err == nil {
		t.Error("an unverified glyph should be an error, not a substitution")
	}
}

// TestQuestionMarkGlyph pins the EXP-0049 addition.
func TestQuestionMarkGlyph(t *testing.T) {
	if r, ok := Decode(0xBF); !ok || r != '?' {
		t.Errorf("Decode($BF) = %q, %v; want '?', true (EXP-0049)", r, ok)
	}
}

func FuzzEncodeDecode(f *testing.F) {
	for _, s := range []string{"WEDGE", "Were-Rat", "0", "29", "?????", ""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		b, err := EncodeFixed(s, len([]rune(s))+2)
		if err != nil {
			return // unverified glyph or overlong: both are valid answers
		}
		if got := DecodeFixed(b); got != trimTrailingSpace(s) {
			t.Fatalf("EncodeFixed(%q) -> DecodeFixed = %q", s, got)
		}
	})
}

func trimTrailingSpace(s string) string {
	r := []rune(s)
	for len(r) > 0 && r[len(r)-1] == ' ' {
		r = r[:len(r)-1]
	}
	return string(r)
}
