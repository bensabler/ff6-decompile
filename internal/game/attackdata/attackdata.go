// Package attackdata reads the 14-byte attack/spell records of Final
// Fantasy III (USA) discovered in EXP-0019: the battle engine loads one
// record per action with MVN from ROM CPU 0xC46AC0 + 14*index (ROM file
// offset 0x046AC0, HiROM) into the action-state block at WRAM
// 0x11A0-0x11AD. Field semantics are mapped from the byte-exact damage
// pipeline decodes (EXP-0001..0017): entry byte k lands at WRAM
// 0x11A0+k, so every verified 0x11Ax meaning names a record field.
//
// The package never embeds ROM data; callers supply record bytes read
// from a locally stored ROM (see local_artifacts/).
package attackdata

import "fmt"

// RecordSize is the attack-record stride: the loader multiplies the
// index by 14 (LDA #$0E / JSR $4781) and MVNs 14 bytes (LDA #$000D).
const RecordSize = 14

// TableCPUAddr is the CPU address of record 0, as it appears in the
// loader's MVN (ROMCPU:$C46AC0).
const TableCPUAddr = 0xC46AC0

// TableFileOffset is the ROM-file offset of record 0 (ROMFILE:0x046AC0),
// the HiROM mapping of TableCPUAddr. The two are asserted against
// internal/platform/snesaddr in this package's tests rather than being
// derived, so the constant stays a constant.
const TableFileOffset = 0x046AC0

// Record is one 14-byte attack/spell entry. Accessors exist only for
// fields whose semantics are verified through the damage-pipeline
// decodes; the rest are reachable through the raw bytes.
type Record [RecordSize]byte

// Decode copies one record from src, which must hold at least
// RecordSize bytes.
func Decode(src []byte) (Record, error) {
	var r Record
	if len(src) < RecordSize {
		return r, fmt.Errorf("decode attack record: need %d bytes, got %d", RecordSize, len(src))
	}
	copy(r[:], src[:RecordSize])
	return r, nil
}

// RecordAt slices record i out of a full table blob (the bytes starting
// at TableFileOffset). It bounds-checks both ends.
func RecordAt(table []byte, i int) (Record, error) {
	if i < 0 {
		return Record{}, fmt.Errorf("attack record index %d is negative", i)
	}
	off := i * RecordSize
	if off+RecordSize > len(table) {
		return Record{}, fmt.Errorf("attack record %d out of range (table holds %d records)", i, len(table)/RecordSize)
	}
	return Decode(table[off:])
}

// Element returns entry byte 1 (WRAM 0x11A1): the attack's element
// mask, consumed by the elemental-modifier block (EXP-0010).
func (r Record) Element() uint8 { return r[1] }

// Flags2 returns entry byte 2 (WRAM 0x11A2): bit 0 selects the
// physical base formula (EXP-0017), bit 7 dispatches the MP delta
// path (EXP-0001). Other bits observed in gates but unnamed.
func (r Record) Flags2() uint8 { return r[2] }

// PhysicalFormula reports Flags2 bit 0 (EXP-0017 path select).
func (r Record) PhysicalFormula() bool { return r[2]&0x01 != 0 }

// TargetsMP reports byte 3 bit 7 (WRAM 0x11A3): the delta fetch
// retargets every battle array read from HP to MP by adding 0x14 to
// both index registers (EXP-0011).
func (r Record) TargetsMP() bool { return r[3]&0x80 != 0 }

// Mode returns entry byte 4 (WRAM 0x11A4): bit 7 selects the
// base-amount variant (fraction-of-HP when set, EXP-0011), bit 1
// participates in the damage/heal flip chain (EXP-0010), bit 0 gates
// the party-vs-party halving (EXP-0011).
func (r Record) Mode() uint8 { return r[4] }

// Power returns entry byte 6 (WRAM 0x11A6): the attack power consumed
// by both base formulas (EXP-0014/0015 live-verified).
func (r Record) Power() uint8 { return r[6] }

// AbortBits returns entry byte 10 (WRAM 0x11AA): tested with #$82 in
// the polarity flip chain (EXP-0010); semantics unknown beyond the
// abort effect.
func (r Record) AbortBits() uint8 { return r[10] }
