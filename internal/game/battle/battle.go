// Package battle reconstructs the battle-slot HP/MP delta engine observed
// in bank 0xC2 of Final Fantasy III (USA): the dispatched routine pair at
// ROM CPU 0xC21323 (HP) and 0xC21350 (MP), reached via JSR ($131F,X), and
// the death handler at 0xC21390. The engine is slot-uniform: one 10-entry
// array family serves party entries 0-3 and enemy entries 4-9 (EXP-0005).
//
// Provenance: the session that produced this package (Session 003,
// docs/sessions/SESSION_003.md) was interrupted before documentation and
// its original ROM dumps were lost, so the arithmetic here briefly stood
// as an unverified reconstruction. EXP-0001
// (docs/experiments/EXP-0001-c2-delta-engine-dump.md) re-dumped
// ROMCPU $C212F0-$C2141F and verified every instruction this package
// models byte-exact: dispatch JSR ($131F,X) at 0xC21300 selecting
// 0xC21323 (HP) or 0xC21350 (MP) via bit 7 of WRAM 0x11A2; heal store
// 0xC21338 clamped to 0x3C1C,Y; damage store 0xC21347; death handler
// 0xC21390 with zero store 0xC21396 and the 0x3EE4,Y bit-1 suppression.
// Semantic labels (HP vs MP, "max", "death") remain strong hypotheses
// from battle context. Full annotated listing:
// docs/sessions/02_DISCOVERED_FUNCTIONS.md; structures:
// docs/sessions/05_DATA_STRUCTURES.md.
//
// The queue side of the pending-delta arrays (0x33E4/0x33D0, 0xFFFF
// "none" sentinels) is modeled by AccumulatePending (EXP-0006, byte-exact
// vs ROM CPU 0xC20C76). Not modeled here (insufficient evidence, see
// docs/sessions/08_OPEN_QUESTIONS.md): the delta-fetch routine at
// 0xC213A7 (its gate reads are semantically unresolved), the MP routine's
// exit tail (LDA #$0080 / JMP $464C, purpose unknown), and the death
// handler's clearing of 0x3A89.
package battle

// BattleSlotCount is the number of entries in every battle-scoped
// per-slot array: entries 0-3 are the party slots, 4-9 the enemy slots.
// Confirmed by EXP-0004/0005: the four main arrays sit exactly 0x14
// bytes apart (10 x 16-bit entries), enemy HP was observed live at
// entries 4-5, and the same delta-engine stores operate party and enemy
// entries alike. Entries 6-9 are stride-bounded but not yet observed in
// use (no encounter with >2 enemies instrumented).
const BattleSlotCount = 10

// PartySlotCount is the number of party entries (0-3) — the slice of the
// arrays that the display copier (ROM CPU 0xC25D26) exposes to the HUD.
const PartySlotCount = 4

// WRAM-relative offsets (SNES CPU bank 0x7E) of the authoritative
// battle-scoped arrays, ten 16-bit entries each (0x14-byte stride between
// the family's base addresses). The whole region is filled with 0xFF at
// battle teardown and rewritten at battle init, so these values only
// exist during a battle.
const (
	AddrCurrentHP         = 0x3BF4
	AddrCurrentMP         = 0x3C08
	AddrMaxHP             = 0x3C1C
	AddrMaxMP             = 0x3C30
	AddrUnknownStatus3EE4 = 0x3EE4
	AddrUnknown3EF8       = 0x3EF8
	AddrUnknownFlags3C95  = 0x3C95
)

// BattleSlots holds the authoritative per-slot battle values operated on
// by the delta engine — party entries 0-3, enemy entries 4-9. Layout
// follows the original struct-of-arrays form.
type BattleSlots struct {
	CurrentHP [BattleSlotCount]uint16 // $3BF4
	CurrentMP [BattleSlotCount]uint16 // $3C08
	MaxHP     [BattleSlotCount]uint16 // $3C1C: heals clamp to this (code-confirmed)
	MaxMP     [BattleSlotCount]uint16 // $3C30: MP heals clamp to this (code-confirmed)

	// UnknownStatus3EE4 is a per-slot 16-bit flags field ($3EE4; the
	// family keeps the 0x14 stride: $3EE4+0x14 = $3EF8). Bit 1
	// suppresses the death event when HP reaches zero ($C2139C:
	// BIT #$0002); the death handler reads it for enemy entries too.
	// Other bits unknown. It is also copied (masked) into display records.
	UnknownStatus3EE4 [BattleSlotCount]uint16

	// DiesAtZeroMP mirrors bit 0 of the per-slot array at $3C95: when set,
	// current MP reaching zero invokes the death handler ($C21380-$C21386).
	// The rest of $3C95 is unknown; its width is assumed to match the
	// family (10) — only party reads are observed so far.
	DiesAtZeroMP [BattleSlotCount]bool
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
func (p *BattleSlots) ApplyHPDelta(slot int, delta int32) bool {
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
func (p *BattleSlots) ApplyMPDelta(slot int, delta int32) bool {
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
func (p *BattleSlots) deathEvent(slot int) bool {
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

// PendingDeltaNone is the "no pending delta" sentinel observed throughout
// the pending arrays ($33D0/$33E4 families): battle init and the sweepers
// fill entries with 0xFFFF, and both the accumulator and the fetch treat
// it as zero.
const PendingDeltaNone = 0xFFFF

// PendingDeltaCap is the accumulation ceiling enforced by the original
// (CMP #$2710 / LDA #$270F at ROM CPU 0xC20C90-0xC20C97): 9999, the
// classic displayed-damage cap, applied at queue time.
const PendingDeltaCap = 0x270F

// AccumulatePending reproduces the pending-delta accumulator at ROM CPU
// 0xC20C76 (byte-exact per EXP-0006): the slot's current pending value
// (0xFFFF sentinel meaning "none") plus a new amount, clamped to 9999.
// A 16-bit overflow during the add also clamps, matching the BCS path.
// The caller decides which pending array (damage or heal) the result
// belongs to; the original retargets +$33D0 to +$33E4 by adding 0x14 to
// Y before this arithmetic.
func AccumulatePending(current, amount uint16) uint16 {
	if current == PendingDeltaNone {
		current = 0
	}
	sum := uint32(current) + uint32(amount)
	if sum > 0xFFFF || sum >= 0x2710 {
		return PendingDeltaCap
	}
	return uint16(sum)
}

// ElementResponse holds one slot's per-element response masks, named by
// their verified effect on the pending amount (ROM CPU 0xC20BE2-0xC20C1D,
// EXP-0010). Candidate semantic labels, unproven: FlipToHeal=absorb,
// Nullify=immune, Halve=resist, Double=weak. In WRAM the masks live in
// two more 0x14-stride 10-entry arrays: FlipToHeal/Nullify are the
// low/high bytes of the 16-bit entry at 0x3BCC+slot*2, Halve/Double the
// high/low bytes at 0x3BE0+slot*2.
type ElementResponse struct {
	FlipToHeal uint8 // $3BCC,Y: matching element flips damage<->heal
	Nullify    uint8 // $3BCD,Y: matching element zeroes the amount
	Halve      uint8 // $3BE1,Y: matching element halves the amount
	Double     uint8 // $3BE0,Y: matching element doubles the amount
}

// ApplyElementResponse reproduces the elemental-modifier block at ROM CPU
// 0xC20BD3-0xC20C1D (byte-exact per EXP-0010): given the running amount,
// the current heal flag (bit 0 of DP $F2), the attack's element byte
// ($11A1), and the battle-wide nullified-element byte ($3EC8), apply the
// target's response masks with first-match-wins priority:
//
//	no elements        -> unchanged
//	all nullified      -> amount 0 (~nullified & elems == 0)
//	FlipToHeal match   -> heal = !heal
//	else Nullify match -> amount 0
//	else Halve match   -> amount >>= 1
//	else Double match  -> amount <<= 1, only when amount < 0x8000
//
// Note the original tests each mask against the full element byte, not
// the nullify-masked value; this function preserves that exactly.
func ApplyElementResponse(amount uint16, heal bool, elems, nullified uint8, r ElementResponse) (uint16, bool) {
	if elems == 0 {
		return amount, heal
	}
	if ^nullified&elems == 0 {
		return 0, heal
	}
	switch {
	case r.FlipToHeal&elems != 0:
		heal = !heal
	case r.Nullify&elems != 0:
		amount = 0
	case r.Halve&elems != 0:
		amount >>= 1
	case r.Double&elems != 0:
		if amount < 0x8000 {
			amount <<= 1
		}
	}
	return amount, heal
}
