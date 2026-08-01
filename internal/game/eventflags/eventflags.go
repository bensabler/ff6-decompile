// Package eventflags models the FF6 event-flag bit arrays and their
// flag-number addressing, as located by CONTRA-0002 and inventoried by
// EXP-0037.
//
// Three parallel 32-byte bit arrays live at WRAM:+$1E80, +$1EA0 and
// +$1EC0 (256 flags each). A shared ROM decoder at ROMCPU:$C0BAED maps
// a flag number to array indices (byte = flag/8 via LSR x3, bit =
// flag&7), and set/clear handlers apply single-bit masks from the ROM
// tables at ROMCPU:$C0BAFC (set) and $C0BB04 (clear). Static decode
// shows the same handler family extends to at least four further bases
// ($1EE0, $1F00, $1F20, $1F40); those arrays are registered in the
// census but not modeled here until runtime evidence covers them.
package eventflags

import "fmt"

// ArrayBase is the WRAM offset of one 32-byte event-flag bit array.
type ArrayBase uint16

// The three runtime-verified arrays (CONTRA-0002; EXP-0037 write
// inventory). Set handlers: $C0B593/$C0B5A7/$C0B5BB. Clear handlers:
// $C0B5CF/$C0B5E3/$C0B5F7.
const (
	ArrayA ArrayBase = 0x1E80
	ArrayB ArrayBase = 0x1EA0
	ArrayC ArrayBase = 0x1EC0
)

// arrayLen is the byte length of one array ($20), covering 256 flags.
const arrayLen = 0x20

// SetMasks mirrors the ROM set-mask table at ROMCPU:$C0BAFC.
var SetMasks = [8]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}

// ClearMasks mirrors the ROM clear-mask table at ROMCPU:$C0BB04
// (exact complements of SetMasks).
var ClearMasks = [8]byte{0xFE, 0xFD, 0xFB, 0xF7, 0xEF, 0xDF, 0xBF, 0x7F}

// Ref locates one flag: its array, byte offset within the array, and
// bit index, exactly as the ROM decoder derives them.
type Ref struct {
	Array ArrayBase
	Flag  uint8
}

// ByteOffset returns the byte index inside the array (flag/8), the Y
// register produced by ROMCPU:$C0BAED.
func (r Ref) ByteOffset() uint8 { return r.Flag >> 3 }

// Bit returns the bit index (flag&7), the X register produced by
// ROMCPU:$C0BAED.
func (r Ref) Bit() uint8 { return r.Flag & 7 }

// ByteAddr returns the WRAM offset of the byte holding the flag.
func (r Ref) ByteAddr() uint16 { return uint16(r.Array) + uint16(r.ByteOffset()) }

// SetMask returns the ORA mask the set handlers apply ($C0BAFC table).
func (r Ref) SetMask() byte { return SetMasks[r.Bit()] }

// ClearMask returns the AND mask the clear handlers apply ($C0BB04 table).
func (r Ref) ClearMask() byte { return ClearMasks[r.Bit()] }

// ID renders the project identifier used by EXP-0037 records,
// e.g. "EVF-1EA0-$2B".
func (r Ref) ID() string { return fmt.Sprintf("EVF-%04X-$%02X", uint16(r.Array), r.Flag) }

// FlagAt is the inverse mapping: which flag does a given WRAM byte
// address and bit index belong to? It bounds-checks the address against
// the three modeled arrays.
func FlagAt(byteAddr uint16, bit uint8) (Ref, error) {
	if bit > 7 {
		return Ref{}, fmt.Errorf("eventflags: bit %d out of range", bit)
	}
	for _, base := range []ArrayBase{ArrayA, ArrayB, ArrayC} {
		if byteAddr >= uint16(base) && byteAddr < uint16(base)+arrayLen {
			off := uint8(byteAddr - uint16(base))
			return Ref{Array: base, Flag: off<<3 | bit}, nil
		}
	}
	return Ref{}, fmt.Errorf("eventflags: WRAM:+$%04X is not inside a modeled flag array", byteAddr)
}

// IsSet reports whether the flag is set in a 32-byte array snapshot.
func (r Ref) IsSet(array []byte) (bool, error) {
	off := int(r.ByteOffset())
	if off >= len(array) {
		return false, fmt.Errorf("eventflags: snapshot too short (%d bytes) for flag $%02X", len(array), r.Flag)
	}
	return array[off]&r.SetMask() != 0, nil
}
