// Package battledata decodes Final Fantasy III (USA)'s formation and monster
// record formats.
//
// Both layouts are Confirmed, and both are only *partly* mapped. This package
// exposes what is verified through named accessors and keeps every record's
// raw bytes reachable, so an undecoded field is visible rather than absent.
// Nothing here guesses a field's meaning.
//
// # Addresses
//
//	ROMCPU:$CF6200 + 15*id   formation records   (EXP-0030, ROM-0024)
//	ROMCPU:$CF5900 +  4*id   formation flags     (EXP-0030, ROM-0025)
//	ROMCPU:$CF0000 + 32*id   monster records     (EXP-0028/0029, ROM-0022)
//
// # Verified extent
//
// The archived tables are larger than the verified spans. Formations are
// Confirmed through record 44 and monsters through record 78; the extractor
// archives 576 and 384 respectively because the true table lengths are
// unknown. Decoding a record beyond the verified span is legitimate — it is
// how the span grows — but callers should not treat it as Confirmed. Verified
// reports the distinction.
package battledata

import (
	"encoding/binary"
	"fmt"
)

// Formation record geometry (EXP-0030).
const (
	FormationSize = 15
	// FormationCPUAddr is record 0's CPU address.
	FormationCPUAddr = 0xCF6200
	// FormationMonsterSlots is the count of monster-id bytes, at +$02..+$07.
	FormationMonsterSlots = 6
	// EmptySlot marks an unused monster slot.
	EmptySlot byte = 0xFF
	// FormationsVerified is the highest formation record any experiment has
	// verified (EXP-0030, record 44 = the mines encounter).
	FormationsVerified = 44
)

// Monster record geometry (EXP-0028/0029).
const (
	MonsterSize = 32
	// MonsterCPUAddr is record 0's CPU address.
	MonsterCPUAddr = 0xCF0000
	// MonstersVerified is the highest monster record any experiment has
	// cross-checked (EXP-0028/0030, records 77 and 78).
	MonstersVerified = 78
)

// Formation is one 15-byte formation record.
//
// Mapped: bytes +$02..+$07 are monster ids, `$FF` meaning empty. Unknown:
// the leading word at +$00/+$01 and the trailing bytes +$08..+$0E
// (CEN-MONSTER-0002).
type Formation [FormationSize]byte

// DecodeFormation copies one record from src.
func DecodeFormation(src []byte) (Formation, error) {
	var f Formation
	if len(src) < FormationSize {
		return f, fmt.Errorf("decode formation: need %d bytes, got %d", FormationSize, len(src))
	}
	copy(f[:], src[:FormationSize])
	return f, nil
}

// FormationAt slices record id out of a full table blob.
func FormationAt(table []byte, id int) (Formation, error) {
	var f Formation
	if id < 0 {
		return f, fmt.Errorf("formation id %d is negative", id)
	}
	off := id * FormationSize
	if off+FormationSize > len(table) {
		return f, fmt.Errorf("formation %d needs bytes %d-%d, table has %d", id, off, off+FormationSize-1, len(table))
	}
	copy(f[:], table[off:off+FormationSize])
	return f, nil
}

// LeadingWord returns the 16-bit little-endian value at +$00/+$01.
//
// Its meaning is **Unknown**. It is exposed because it is the standing
// candidate for the monster-id extension the Whelk formation needs: formation
// 432's low id bytes read as records 0 and 52, which cannot be right — record
// 0 is the 40 HP opening guard, while EXP-0040 measured Whelk's two entities
// at 50000 and 1600 HP. Something outside the id bytes must be selecting a
// different record, and this word is where it would live.
func (f Formation) LeadingWord() uint16 { return binary.LittleEndian.Uint16(f[0:2]) }

// MonsterIDs returns the six raw id bytes at +$02..+$07, including `$FF`
// markers.
func (f Formation) MonsterIDs() [FormationMonsterSlots]byte {
	var out [FormationMonsterSlots]byte
	copy(out[:], f[2:2+FormationMonsterSlots])
	return out
}

// OccupiedIDs returns the non-empty monster ids, in slot order.
//
// The values are the **low bytes only**. Whether a record can address monsters
// above 255 is undecoded, so a caller resolving these against the monster
// table is resolving an id this project has verified only for formations
// within FormationsVerified.
func (f Formation) OccupiedIDs() []byte {
	var out []byte
	for _, id := range f.MonsterIDs() {
		if id != EmptySlot {
			out = append(out, id)
		}
	}
	return out
}

// TrailingBytes returns +$08..+$0E, whose semantics are Unknown.
func (f Formation) TrailingBytes() []byte { return f[8:FormationSize] }

// Verified reports whether an experiment has checked this formation id
// against live battle state.
func FormationVerified(id int) bool { return id >= 0 && id <= FormationsVerified }

// Monster is one 32-byte monster stat record.
//
// Mapped: +$01 battle power, +$08 HP, +$0A MP. Unknown: +$06, +$0D, +$12,
// +$14..+$19, +$1B..+$1D, +$1F, and the reward fields (CEN-MONSTER-0001).
type Monster [MonsterSize]byte

// DecodeMonster copies one record from src.
func DecodeMonster(src []byte) (Monster, error) {
	var m Monster
	if len(src) < MonsterSize {
		return m, fmt.Errorf("decode monster: need %d bytes, got %d", MonsterSize, len(src))
	}
	copy(m[:], src[:MonsterSize])
	return m, nil
}

// MonsterAt slices record id out of a full table blob.
func MonsterAt(table []byte, id int) (Monster, error) {
	var m Monster
	if id < 0 {
		return m, fmt.Errorf("monster id %d is negative", id)
	}
	off := id * MonsterSize
	if off+MonsterSize > len(table) {
		return m, fmt.Errorf("monster %d needs bytes %d-%d, table has %d", id, off, off+MonsterSize-1, len(table))
	}
	copy(m[:], table[off:off+MonsterSize])
	return m, nil
}

// Power returns the battle power at +$01.
//
// Confirmed by the EXP-0030 cross-check: formation 44's monsters {19, 77}
// carry powers {13, 19}, the exact pair EXP-0018 measured live.
func (m Monster) Power() uint8 { return m[0x01] }

// HP returns the 16-bit value at +$08, which battle init stages as both
// current and max HP (ROMCPU:$C22CCE is the max-HP store).
func (m Monster) HP() uint16 { return binary.LittleEndian.Uint16(m[0x08:0x0A]) }

// MP returns the 16-bit value at +$0A.
func (m Monster) MP() uint16 { return binary.LittleEndian.Uint16(m[0x0A:0x0C]) }

// MonsterVerified reports whether an experiment has cross-checked this record.
func MonsterVerified(id int) bool { return id >= 0 && id <= MonstersVerified }

// Table wraps a decoded table blob with its record count.
type Table struct {
	raw    []byte
	stride int
}

// NewTable wraps a blob, rejecting one that is not a whole number of records.
func NewTable(raw []byte, stride int) (*Table, error) {
	if stride <= 0 {
		return nil, fmt.Errorf("record stride %d", stride)
	}
	if len(raw)%stride != 0 {
		return nil, fmt.Errorf("table of %d bytes is not a whole number of %d-byte records", len(raw), stride)
	}
	return &Table{raw: raw, stride: stride}, nil
}

// Count returns the number of records.
func (t *Table) Count() int { return len(t.raw) / t.stride }

// Raw returns the record bytes for id.
func (t *Table) Raw(id int) ([]byte, error) {
	if id < 0 || id >= t.Count() {
		return nil, fmt.Errorf("record %d out of range 0-%d", id, t.Count()-1)
	}
	off := id * t.stride
	out := make([]byte, t.stride)
	copy(out, t.raw[off:off+t.stride])
	return out, nil
}
