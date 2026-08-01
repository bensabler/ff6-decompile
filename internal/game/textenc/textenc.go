// Package textenc decodes the fixed-width menu/name text encoding of the
// verified FF3us ROM.
//
// The glyph values were derived on-screen in EXP-0026 (menu tilemap
// cross-checked against ROM byte searches and the battle HUD tilemap) and
// exercised across all 54 spell names and 27 esper names in EXP-0027:
//
//	'A'..'Z' = $80.. ; 'a'..'z' = $9A.. ; '0'..'9' = $B4.. ;
//	'-' = $C4 ; narrow space = $FE ; blank/padding = $FF.
//
// Of the digit run only '2'=$B6 and '3'=$B7 were individually observed
// on-screen; the $B4 base is the arithmetic consequence and is carried at
// the same confidence as the letter runs. Every other byte is UNMAPPED and
// decodes to an explicit \xNN escape — unknown glyphs are surfaced, never
// guessed.
package textenc

import "fmt"

// Decode maps one encoded byte to its glyph. ok is false for unmapped
// bytes.
func Decode(b byte) (r rune, ok bool) {
	switch {
	case b >= 0x80 && b <= 0x99:
		return rune('A' + b - 0x80), true
	case b >= 0x9A && b <= 0xB3:
		return rune('a' + b - 0x9A), true
	case b >= 0xB4 && b <= 0xBD:
		return rune('0' + b - 0xB4), true
	case b == 0xC4:
		return '-', true
	case b == 0xFE, b == 0xFF:
		return ' ', true
	default:
		return 0, false
	}
}

// DecodeFixed decodes a fixed-width name field. Trailing padding/space is
// trimmed; unmapped bytes are rendered as \xNN escapes so they are visible
// in output rather than silently dropped.
func DecodeFixed(src []byte) string {
	buf := make([]rune, 0, len(src)*4)
	for _, b := range src {
		if r, ok := Decode(b); ok {
			buf = append(buf, r)
			continue
		}
		buf = append(buf, []rune(fmt.Sprintf(`\x%02X`, b))...)
	}
	// trim trailing spaces (both $FE and $FF decode to ' ')
	end := len(buf)
	for end > 0 && buf[end-1] == ' ' {
		end--
	}
	return string(buf[:end])
}

// Clean reports whether every byte of src is a mapped glyph (padding
// included) — i.e. DecodeFixed produced no escapes.
func Clean(src []byte) bool {
	for _, b := range src {
		if _, ok := Decode(b); !ok {
			return false
		}
	}
	return true
}
