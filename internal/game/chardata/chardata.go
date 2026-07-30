// Package chardata reconstructs the character-field copy routine at ROM CPU
// address 0xC10DF3 of Final Fantasy III (USA).
//
// The routine was fully disassembled and verified live in Session 002: it
// runs once per frame during battle (called from 0xC10200 inside the
// per-frame battle update sequence at 0xC101FB), copies six 16-bit fields
// per party slot from a struct-of-arrays source region (WRAM 0x2E78-0x2EA7)
// into four 0x20-byte records based at WRAM 0x2EB5, optionally masking the
// last two fields, and derives a 4-bit per-slot mask that the original
// stores to WRAM 0x61AD.
//
// Evidence and confidence levels are documented in
// docs/sessions/02_DISCOVERED_FUNCTIONS.md and
// docs/sessions/05_DATA_STRUCTURES.md; names
// follow the project naming rules, so fields without demonstrated meaning
// are called Unknown, and fields with strong but unconfirmed meaning carry
// a Candidate marker.
package chardata

// PartySlotCount is the number of records the copy loop processes:
// X advances by 2 until CPX #$0008, giving four iterations.
const PartySlotCount = 4

// WRAM-relative offsets where each reconstructed region was observed
// (SNES CPU bank 0x7E). Recorded for traceability back to the debugger
// evidence; the Go code does not emulate the SNES address space.
const (
	// AddrCharacterFieldsSource is the base of the struct-of-arrays
	// source region read by the copy loop (0x2E78-0x2EA7).
	AddrCharacterFieldsSource = 0x2E78

	// AddrPartySlotRecordBase is the base of the first destination
	// record. Its first field is the confirmed current HP of the first
	// displayed party slot.
	AddrPartySlotRecordBase = 0x2EB5

	// PartySlotRecordStride is the spacing between destination records:
	// the loop advances Y by 0x20 per iteration.
	PartySlotRecordStride = 0x20

	// AddrUnknownFlag628D and AddrUnknownFlagE9EF are the two flag bytes
	// the routine's prologue tests to decide whether to mask the last
	// two copied fields.
	AddrUnknownFlag628D = 0x628D
	AddrUnknownFlagE9EF = 0xE9EF

	// AddrSlotMask61AD is where the original routine stores the 4-bit
	// per-slot mask computed by its second loop.
	AddrSlotMask61AD = 0x61AD
)

// CharacterFieldsSource is the source region the copy loop reads: six
// parallel arrays (struct-of-arrays), one array per field, one entry per
// displayed party slot. Slot order matches the on-screen party list
// (verified against displayed HP values in Session 002).
type CharacterFieldsSource struct {
	CurrentHP      [PartySlotCount]uint16 // $2E78,X -> record +0x00; matches displayed HP (slots 0-2 verified)
	CandidateMaxHP [PartySlotCount]uint16 // $2E80,X -> record +0x02; a full heal set CurrentHP to exactly this value
	Unknown04      [PartySlotCount]uint16 // $2E88,X -> record +0x04; observed 24,0,0,0 (current MP candidate)
	Unknown06      [PartySlotCount]uint16 // $2E90,X -> record +0x06; observed 24,0,0,0 (max MP candidate)
	Unknown08      [PartySlotCount]uint16 // $2E98,X -> record +0x08, masked with 0x0038 in masked mode
	Unknown0A      [PartySlotCount]uint16 // $2EA0,X -> record +0x0A, zeroed in masked mode; bit 13 drives SlotMask
}

// PartySlotRecord is one 0x20-byte destination record; record n is based at
// WRAM 0x2EB5 + n*0x20. Only offsets +0x00 through +0x0B are written by the
// copy routine; bytes +0x0C onward hold live data written by other,
// still-unidentified code, so they are preserved explicitly. Whether 0x2EB5
// is the true start of the record is an open question.
type PartySlotRecord struct {
	CurrentHP      uint16     // +0x00: confirmed (records 0-2 matched displayed HP)
	CandidateMaxHP uint16     // +0x02: strong hypothesis (heal target value, HP gauge ratio)
	Unknown04      uint16     // +0x04
	Unknown06      uint16     // +0x06
	Unknown08      uint16     // +0x08
	Unknown0A      uint16     // +0x0A
	Unknown0C      [0x14]byte // +0x0C-0x1F: written by other code, never by this routine
}

// CopyMode mirrors the two WRAM flag bytes the routine's prologue tests.
// Masked mode (Flag628D clear and FlagE9EF set) ANDs field +0x08 with
// 0x0038 and zeroes field +0x0A; every other combination copies both
// fields unmodified. The gameplay meaning of both flags is unknown; both
// were zero throughout the observed battle.
type CopyMode struct {
	UnknownFlag628D bool // WRAM $628D nonzero: never mask
	UnknownFlagE9EF bool // WRAM $E9EF nonzero (and $628D zero): mask
}

// CopyCharacterFields reproduces the routine at ROM CPU 0xC10DF3
// (behavioral reconstruction, byte-exact against the disassembly): the
// six-field copy loop at 0xC10E11 followed by the mask-derivation loop at
// 0xC10E4D.
//
// The returned slot mask is what the original stores to WRAM 0x61AD: bit n
// is set when bit 13 of source field Unknown0A for slot n is clear. It is
// derived from the unmasked source values, so it is independent of mode.
func CopyCharacterFields(src *CharacterFieldsSource, dst *[PartySlotCount]PartySlotRecord, mode CopyMode) uint8 {
	mask08 := uint16(0xFFFF)
	mask0A := uint16(0xFFFF)
	if !mode.UnknownFlag628D && mode.UnknownFlagE9EF {
		mask08 = 0x0038
		mask0A = 0x0000
	}

	for slot := range dst {
		dst[slot].CurrentHP = src.CurrentHP[slot]
		dst[slot].CandidateMaxHP = src.CandidateMaxHP[slot]
		dst[slot].Unknown04 = src.Unknown04[slot]
		dst[slot].Unknown06 = src.Unknown06[slot]
		dst[slot].Unknown08 = src.Unknown08[slot] & mask08
		dst[slot].Unknown0A = src.Unknown0A[slot] & mask0A
	}

	var slotMask uint8
	for slot := range src.Unknown0A {
		if src.Unknown0A[slot]&0x2000 == 0 {
			slotMask |= 1 << slot
		}
	}
	return slotMask
}
