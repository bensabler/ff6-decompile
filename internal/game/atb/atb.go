// Package atb models Final Fantasy III (USA)'s battle timing layer.
//
// Everything here is Confirmed by EXP-0041 through EXP-0047 and byte-exact
// against the decoded routines. Where the evidence stops, this package stops:
// the increment *formula* and the readiness *threshold* are open questions, so
// increments are an input rather than something computed here. Guessing them
// would put an unverified model behind a confident API.
//
// # The stride trap
//
// The ATB arrays are **stride 2**, not the `$14` of the HP/stat family
// (DISC-0001). Slot n lives at index n*2. Mixing the two strides is the most
// likely way to misread this layer, so SlotIndex exists and the raw offset is
// never written by hand.
//
// # Addresses
//
//	WRAM:+$3AB4  gauges, 10 entries, stride 2, 16-bit   (EXP-0043)
//	WRAM:+$3AC8  per-slot increment, same shape         (EXP-0043)
//	WRAM:+$3AA0  per-slot scheduler flags               (EXP-0043)
//	WRAM:+$3A3E  16-bit battle tick counter             (EXP-0043)
//	WRAM:+$3A8F  Wait flag, sampled at battle entry     (EXP-0042)
//	WRAM:+$3A90  255 - 24*BattleSpeed, sampled at entry (EXP-0042)
//	WRAM:+$2F41  battle submenu flag                    (EXP-0044)
//	WRAM:+$1D4D  Bat.Speed bits 0-2, Bat.Mode bit 3     (EXP-0041)
package atb

import "fmt"

// WRAM addresses, as offsets into the 128 KB work RAM image.
const (
	AddrGauges       = 0x3AB4
	AddrIncrements   = 0x3AC8
	AddrFlags        = 0x3AA0
	AddrTickCounter  = 0x3A3E
	AddrWaitFlag     = 0x3A8F
	AddrSpeedByte    = 0x3A90
	AddrSubmenuFlag  = 0x2F41
	AddrConfigTiming = 0x1D4D
)

// SlotCount is the ten battle slots: 0-3 party, 4-9 enemy (DISC-0001).
const SlotCount = 10

// Stride is the ATB arrays' element stride in bytes. It is 2, not the $14 of
// the HP/stat arrays.
const Stride = 2

// PartySlots is the number of leading slots the party occupies. The battle
// code distinguishes them by index, not by a flag: `CPX #$08 / BCC` at
// ROMCPU:$C209F6 skips the Battle Speed multiply for X < 8, and X is the
// stride-2 index, so slots 0-3.
const PartySlots = 4

// SlotIndex returns the byte index of a slot within a stride-2 ATB array.
func SlotIndex(slot int) int { return slot * Stride }

// IsPartySlot reports whether a slot belongs to the party.
func IsPartySlot(slot int) bool { return slot >= 0 && slot < PartySlots }

// Scheduler flag bits in WRAM:+$3AA0, decoded at ROMCPU:$C2114B and
// ROMCPU:$C211B2.
const (
	// FlagActive is bit 0. ROMCPU:$C2114E does LSR / BCC, so a clear bit
	// skips the slot's entire per-frame update.
	FlagActive uint8 = 1 << 0
	// FlagReady is bit 5, set at ROMCPU:$C211B2 when the gauge fills.
	FlagReady uint8 = 1 << 5
	// FlagPendingAction is bit 6 (EXP-0043).
	FlagPendingAction uint8 = 1 << 6
	// FlagBusy is bit 7, cleared on action completion (EXP-0043).
	FlagBusy uint8 = 1 << 7
)

// Mode is the ACTIVE/WAIT setting.
type Mode uint8

const (
	// Active is Bat.Mode 0: the battle never pauses for a submenu.
	Active Mode = iota
	// Wait is Bat.Mode 1 (WRAM:+$1D4D bit 3 set).
	Wait
)

func (m Mode) String() string {
	if m == Wait {
		return "Wait"
	}
	return "Active"
}

// Config is the timing half of the battle configuration, as stored in
// WRAM:+$1D4D (EXP-0041).
type Config struct {
	// Mode is ACTIVE or WAIT.
	Mode Mode
	// Speed is the stored Battle Speed, 0-5. The Config screen displays it
	// as 1-6; DisplaySpeed converts.
	Speed uint8
}

// DecodeConfig reads the timing settings out of WRAM:+$1D4D.
//
// Bit 3 set means Wait. The Config screen marks the *selected* option with
// tile attribute $20 and the unselected with $28 — the inverse of the
// intuitive reading, and exactly why EXP-0040 misread this byte and lost two
// Whelk attempts to an unmodelled timing system. Read it from memory.
func DecodeConfig(b uint8) Config {
	c := Config{Speed: b & 0x07}
	if b&0x08 != 0 {
		c.Mode = Wait
	}
	return c
}

// DisplaySpeed returns the 1-6 value the Config screen shows.
func (c Config) DisplaySpeed() int { return int(c.Speed) + 1 }

// SpeedByte returns the value battle entry stores at WRAM:+$3A90:
// 255 - 24*Speed (EXP-0042, ROMCPU:$C22472).
//
// It scales **enemy** gauges only. ROMCPU:$C209F6's `CPX #$08 / BCC` skips
// the multiply for party slots, and EXP-0043 measured party increments
// byte-identical at Battle Speed 3 and 6 while enemy increments went
// 240 -> 156.
func (c Config) SpeedByte() uint8 { return 255 - 24*c.Speed }

// WaitFlag returns the value battle entry stores at WRAM:+$3A8F.
func (c Config) WaitFlag() uint8 {
	if c.Mode == Wait {
		return 1
	}
	return 0
}

// Entry is the timing state sampled once at battle entry (ROMCPU:$C22472).
//
// Neither field is re-read during the battle. That is the program's staging
// rule: ACTIVE/WAIT and Battle Speed must be set *before* entry, or injected
// at these two cells (EXP-0042).
type Entry struct {
	WaitFlag  uint8 // WRAM:+$3A8F
	SpeedByte uint8 // WRAM:+$3A90
}

// SampleEntry performs the battle-entry decomposition.
func SampleEntry(c Config) Entry {
	return Entry{WaitFlag: c.WaitFlag(), SpeedByte: c.SpeedByte()}
}

// Paused reports whether the per-frame battle update is skipped this frame.
//
// ROMCPU:$C21124 is `LDA $2F41 / AND $3A8F / BNE`, so the update is skipped
// only when both the submenu flag and the Wait flag are non-zero. Everything
// downstream freezes together: gauges, the tick counter, and the second
// accumulator (EXP-0044).
//
// The pause is narrower than the folk model. EXP-0044 established that the
// main battle command window does **not** pause under WAIT — only the ability
// list and target selection do — and EXP-0046/0047 established that the
// action-execution path is periodic and ungated, so work already queued still
// completes while the gate is shut.
func Paused(submenuFlag, waitFlag uint8) bool {
	return submenuFlag&waitFlag != 0
}

// State is the mutable per-frame ATB layer.
type State struct {
	// Gauges is WRAM:+$3AB4.
	Gauges [SlotCount]uint16
	// Increments is WRAM:+$3AC8. It is an input: the formula that produces
	// it is only partly decoded (ROMCPU:$C209F0 reads per-slot Speed and
	// adds $14 before an enemy-only multiply), so this package does not
	// compute it.
	Increments [SlotCount]uint16
	// Flags is WRAM:+$3AA0.
	Flags [SlotCount]uint8
	// Tick is WRAM:+$3A3E, incremented once per non-paused frame.
	Tick uint16
	// Entry is the timing state sampled at battle entry.
	Entry Entry
}

// AdvanceSlot advances one slot's gauge by increment>>1 and reports whether
// it filled.
//
// Byte-exact against ROMCPU:$C21193:
//
//	REP #$20 / LDA $3AC8,X / LSR / CLC / ADC $3AB4,X / STA $3AB4,X
//
// The add is 16-bit and wraps; the carry out is what signals a full gauge.
// EXP-0043 measured a slot advancing by exactly $4E = $9C>>1 per increment,
// and separately caught one wrapping $00B6 -> $0004.
//
// The readiness *threshold* compare that follows (`CMP $322C,X` at
// ROMCPU:$C211A5) is not modelled: `$322C` is undecoded, so only the carry
// path is Confirmed.
func (s *State) AdvanceSlot(slot int) (filled bool, err error) {
	if slot < 0 || slot >= SlotCount {
		return false, fmt.Errorf("atb: slot %d out of range 0-%d", slot, SlotCount-1)
	}
	if s.Flags[slot]&FlagActive == 0 {
		// ROMCPU:$C2114E: LSR / BCC skips the slot entirely.
		return false, nil
	}
	sum := uint32(s.Gauges[slot]) + uint32(s.Increments[slot]>>1)
	s.Gauges[slot] = uint16(sum)
	if sum > 0xFFFF {
		s.Flags[slot] |= FlagReady
		return true, nil
	}
	return false, nil
}

// Frame advances the whole layer one frame, and reports which slots filled.
//
// When the gate is shut nothing moves — not the gauges, not the tick counter.
// ROMCPU:$C21127's BNE jumps past the entire update.
func (s *State) Frame(submenuFlag uint8) (filled []int, err error) {
	if Paused(submenuFlag, s.Entry.WaitFlag) {
		return nil, nil
	}
	s.Tick++
	// ROMCPU:$C21141 loads X = $12 (18 = slot 9 * 2) and counts down by two
	// to zero, so slots are visited high to low.
	for slot := SlotCount - 1; slot >= 0; slot-- {
		f, err := s.AdvanceSlot(slot)
		if err != nil {
			return filled, err
		}
		if f {
			filled = append(filled, slot)
		}
	}
	return filled, nil
}

// Ready reports whether a slot's gauge has filled.
func (s *State) Ready(slot int) bool {
	if slot < 0 || slot >= SlotCount {
		return false
	}
	return s.Flags[slot]&FlagReady != 0
}
