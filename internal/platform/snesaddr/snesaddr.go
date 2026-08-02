// Package snesaddr converts between the SNES CPU address space and file
// offsets in the ROM image, for the HiROM cartridge layout Final Fantasy III
// (USA) uses.
//
// Every offset in this project used to be a hand-computed constant with the
// mapping arithmetic sitting in a comment. That is fine until two records of
// the same fact disagree — which is exactly how the battle HUD font block came
// to be extracted from the wrong address for three days (DEMO-0001 deviation
// D1). This package exists so the arithmetic is written once, tested, and
// assertable against the constants that already exist.
//
// Address notation follows docs/research/ADDRESS_NOTATION.md: CPU addresses
// are ROMCPU (`ROMCPU:$C46AC0`), file offsets are ROMFILE
// (`ROMFILE:0x046AC0`).
//
// # Confidence
//
// The windows below do not carry equal evidence, and the doc comments say so
// rather than presenting one uniform mapping:
//
//   - Banks $C0-$FF are **Confirmed**: 18 of 18 Mesen bridge ROM captures in
//     mesen/out/ are byte-identical to the ROM file at offset = CPU − $C00000
//     (CORR-0001). Every tracked ROM region in manifests/rom-regions.json uses
//     this window.
//   - Banks $40-$7D and the $8000-$FFFF halves of banks $00-$3F / $80-$BF are
//     **Strong hypothesis**: standard HiROM behaviour, but this project has
//     not observed FF6 addressing through them. They are implemented because
//     refusing them would silently push callers back to hand arithmetic, and
//     they are flagged by Window so a caller can require Confirmed mapping.
package snesaddr

import "fmt"

// ImageSize is the byte length of the supported ROM image. It duplicates
// rom.Size deliberately: internal/platform models SNES hardware and must not
// depend on the project's ROM loader. A test asserts the two agree.
const ImageSize = 3145728

// Window names the mapping window a CPU address falls in, so callers can
// distinguish a Confirmed translation from a plausible one.
type Window int

const (
	// WindowNone means the address does not map to ROM.
	WindowNone Window = iota
	// WindowHiROMUpper is banks $C0-$FF, offset = cpu − $C00000. Confirmed.
	WindowHiROMUpper
	// WindowHiROMLower is banks $40-$7D, offset = cpu − $400000. Strong
	// hypothesis: standard HiROM, unobserved in FF6 evidence.
	WindowHiROMLower
	// WindowHiROMMirror is the $8000-$FFFF half of banks $00-$3F and
	// $80-$BF. Strong hypothesis: standard HiROM, unobserved in FF6
	// evidence.
	WindowHiROMMirror
)

// Confirmed reports whether this project has runtime evidence for the window.
func (w Window) Confirmed() bool { return w == WindowHiROMUpper }

func (w Window) String() string {
	switch w {
	case WindowHiROMUpper:
		return "HiROM banks $C0-$FF"
	case WindowHiROMLower:
		return "HiROM banks $40-$7D"
	case WindowHiROMMirror:
		return "HiROM mirror banks $00-$3F/$80-$BF, $8000-$FFFF"
	default:
		return "unmapped"
	}
}

// An Error reports a CPU address that does not name ROM data, and says which
// region it does name so the caller can tell a typo from a WRAM address.
type Error struct {
	CPU    uint32
	Reason string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ROMCPU:$%06X does not map to ROM: %s", e.CPU, e.Reason)
}

// ROMFile converts a CPU address to a ROMFILE offset, reporting the window it
// resolved through. It returns an *Error for addresses that name WRAM, the
// system area, or data past the end of the image.
func ROMFile(cpu uint32) (int, Window, error) {
	if cpu > 0xFFFFFF {
		return 0, WindowNone, &Error{cpu, "address exceeds the 24-bit CPU space"}
	}
	bank := byte(cpu >> 16)
	low := cpu & 0xFFFF

	var off int
	var win Window
	switch {
	case bank >= 0xC0:
		off, win = int(cpu-0xC00000), WindowHiROMUpper
	case bank >= 0x7E && bank <= 0x7F:
		return 0, WindowNone, &Error{cpu, "banks $7E-$7F are WRAM"}
	case bank >= 0x40 && bank <= 0x7D:
		off, win = int(cpu-0x400000), WindowHiROMLower
	case bank <= 0x3F || (bank >= 0x80 && bank <= 0xBF):
		if low < 0x8000 {
			return 0, WindowNone, &Error{cpu, "the $0000-$7FFF half of this bank is the system area (WRAM mirror, PPU/CPU registers), not ROM"}
		}
		off, win = int(uint32(bank&0x3F)<<16|low), WindowHiROMMirror
	default:
		return 0, WindowNone, &Error{cpu, "no HiROM window covers this bank"}
	}

	if off >= ImageSize {
		return 0, WindowNone, &Error{cpu, fmt.Sprintf("maps to ROMFILE:0x%06X, past the end of the %d-byte image", off, ImageSize)}
	}
	return off, win, nil
}

// MustROMFile is ROMFile for addresses proven to be in range, and panics
// otherwise. Use it only in tests and in package-level initialisation of
// values that are themselves covered by a test.
func MustROMFile(cpu uint32) int {
	off, _, err := ROMFile(cpu)
	if err != nil {
		panic(err)
	}
	return off
}

// ROMCPU converts a ROMFILE offset to its canonical CPU address, which is
// always in the Confirmed $C0-$FF window. The lower and mirror windows alias
// the same bytes, so the mapping is many-to-one and this is the inverse that
// round-trips.
func ROMCPU(off int) (uint32, error) {
	if off < 0 || off >= ImageSize {
		return 0, fmt.Errorf("ROMFILE:0x%06X is outside the %d-byte image", off, ImageSize)
	}
	return uint32(off) + 0xC00000, nil
}
