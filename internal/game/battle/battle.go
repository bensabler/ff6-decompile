// Package battle reconstructs the party-slot HP/MP delta engine observed in
// bank 0xC2 of Final Fantasy III (USA): the dispatched routine pair at ROM
// CPU 0xC21323 (HP) and 0xC21350 (MP), reached via JSR ($131F,X), and the
// death handler at 0xC21390.
//
// Provenance (see docs/sessions/SESSION_003.md): the session that produced
// this package was interrupted before documentation, and its ROM dumps were
// not preserved. What survives as raw evidence (mesen/out/events.log) are
// write captures of the damage store (near 0xC21347) firing on enemy hits,
// the heal store (near 0xC21338), and the death zero (near 0xC21396), all
// into the per-slot array at WRAM 0x3BF4 with Y = slot*2 and a dispatch
// return consistent with JSR (abs,X) at 0xC212FF. The remaining detail in
// this file — exact routine entries, the clamp/overflow arithmetic shapes,
// the MP/max arrays, and the status gates — is a reconstruction pending
// re-verification against a fresh ROM dump (open question 1b in
// docs/sessions/08_OPEN_QUESTIONS.md). Evidence records:
// docs/sessions/02_DISCOVERED_FUNCTIONS.md and
// docs/sessions/05_DATA_STRUCTURES.md.
//
// Not modeled here (insufficient evidence, see
// docs/sessions/08_OPEN_QUESTIONS.md):
// the delta-fetch routine at 0xC213A7 (pending-delta arrays 0x33E4/0x33D0
// with 0xFFFF "none" sentinels), the MP routine's exit tail
// (LDA #$0080 / JMP $464C, purpose unknown), and the death handler's
// clearing of 0x3A89.
package battle

// PartySlotCount is the number of party slots in every battle-scoped
// per-slot array (all observed loops process four entries).
const PartySlotCount = 4

// WRAM-relative offsets (SNES CPU bank 0x7E) of the authoritative
// battle-scoped arrays, four 16-bit entries each, one per party slot.
// The whole region is filled with 0xFF at battle teardown and rewritten at
// battle init, so these values only exist during a battle.
const (
	AddrCurrentHP         = 0x3BF4
	AddrCurrentMP         = 0x3C08
	AddrMaxHP             = 0x3C1C
	AddrMaxMP             = 0x3C30
	AddrUnknownStatus3EE4 = 0x3EE4
	AddrUnknown3EF8       = 0x3EF8
	AddrUnknownFlags3C95  = 0x3C95
)

// PartySlots holds the authoritative per-slot battle values operated on by
// the delta engine. Layout follows the original struct-of-arrays form.
type PartySlots struct {
	CurrentHP [PartySlotCount]uint16 // $3BF4
	CurrentMP [PartySlotCount]uint16 // $3C08
	MaxHP     [PartySlotCount]uint16 // $3C1C: heals clamp to this (code-confirmed)
	MaxMP     [PartySlotCount]uint16 // $3C30: MP heals clamp to this (code-confirmed)

	// UnknownStatus3EE4 is a per-slot 16-bit flags field ($3EE4). Bit 1
	// suppresses the death event when HP reaches zero ($C2139C: BIT #$0002).
	// Other bits unknown. It is also copied (masked) into display records.
	UnknownStatus3EE4 [PartySlotCount]uint16

	// DiesAtZeroMP mirrors bit 0 of the per-slot array at $3C95: when set,
	// current MP reaching zero invokes the death handler ($C21380-$C21386).
	// The rest of $3C95 is unknown.
	DiesAtZeroMP [PartySlotCount]bool
}

// ApplyHPDelta reproduces the HP-delta routine at ROM CPU $C21323.
// Positive delta heals: newHP = min(currentHP+delta, maxHP) — a 16-bit
// overflow during the add also clamps to max. Negative delta damages:
// newHP = currentHP-|delta|, and if the result is zero or underflows, HP is
// forced to 0 and the death event fires unless UnknownStatus3EE4 bit 1 is
// set (the original then jumps to $C20E32 with A=$0080).
//
// The returned bool reports that death event. |delta| must fit in 16 bits,
// matching the original's operand width; larger magnitudes are clamped.
func (p *PartySlots) ApplyHPDelta(slot int, delta int32) bool {
	if delta >= 0 {
		p.CurrentHP[slot] = healClamped(p.CurrentHP[slot], clamp16(delta), p.MaxHP[slot])
		return false
	}
	remaining, hitZero := damageFloored(p.CurrentHP[slot], clamp16(-delta))
	p.CurrentHP[slot] = remaining
	if !hitZero {
		return false
	}
	return p.deathEvent(slot)
}

// ApplyMPDelta reproduces the MP-delta routine at ROM CPU $C21350: the same
// clamped arithmetic against CurrentMP/MaxMP. When MP reaches zero (exactly
// or by underflow) and DiesAtZeroMP is set for the slot, the original calls
// the death handler ($C21386: JSR $1390), which zeroes the slot's HP and
// fires the death event under the same status-bit-1 suppression rule; the
// returned bool reports that event.
func (p *PartySlots) ApplyMPDelta(slot int, delta int32) bool {
	if delta >= 0 {
		p.CurrentMP[slot] = healClamped(p.CurrentMP[slot], clamp16(delta), p.MaxMP[slot])
		return false
	}
	remaining, hitZero := damageFloored(p.CurrentMP[slot], clamp16(-delta))
	p.CurrentMP[slot] = remaining
	if !hitZero || !p.DiesAtZeroMP[slot] {
		return false
	}
	p.CurrentHP[slot] = 0
	return p.deathEvent(slot)
}

// deathEvent reproduces the death handler at ROM CPU $C21390 as far as the
// evidence supports: HP is zeroed (callers may already have done so), and
// the event is suppressed when UnknownStatus3EE4 bit 1 is set. The
// original also clears $3A89 and continues into $C20E32 with A=$0080;
// neither is modeled yet.
func (p *PartySlots) deathEvent(slot int) bool {
	p.CurrentHP[slot] = 0
	return p.UnknownStatus3EE4[slot]&0x0002 == 0
}

// healClamped is the shared heal shape: ADC / BCS overflow-clamp /
// CMP max / clamp when >= max.
func healClamped(current, delta, max uint16) uint16 {
	sum := uint32(current) + uint32(delta)
	if sum > 0xFFFF || uint16(sum) >= max {
		return max
	}
	return uint16(sum)
}

// damageFloored is the shared damage shape: SBC, then zero (BEQ) and
// underflow (borrow) both route to the zero floor.
func damageFloored(current, magnitude uint16) (remaining uint16, hitZero bool) {
	if current > magnitude {
		return current - magnitude, false
	}
	return 0, true
}

// clamp16 bounds a non-negative delta magnitude to the 16-bit operand
// width of the original routines.
func clamp16(v int32) uint16 {
	if v > 0xFFFF {
		return 0xFFFF
	}
	return uint16(v)
}
