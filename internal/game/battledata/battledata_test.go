package battledata

import "testing"

// syntheticFormation builds a record with recognisable field values. No ROM
// bytes.
func syntheticFormation() []byte {
	return []byte{
		0x34, 0x12, // leading word $1234
		0x05, 0xFF, 0x07, 0xFF, 0xFF, 0xFF, // ids: 5, empty, 7, empty...
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, // trailing
	}
}

func TestDecodeFormation(t *testing.T) {
	f, err := DecodeFormation(syntheticFormation())
	if err != nil {
		t.Fatal(err)
	}
	if got := f.LeadingWord(); got != 0x1234 {
		t.Errorf("LeadingWord = $%04X, want $1234 (little-endian)", got)
	}
	if got := f.OccupiedIDs(); len(got) != 2 || got[0] != 5 || got[1] != 7 {
		t.Errorf("OccupiedIDs = %v, want [5 7]", got)
	}
	ids := f.MonsterIDs()
	if ids[1] != EmptySlot {
		t.Errorf("MonsterIDs[1] = $%02X, want $FF; empty slots must survive the raw accessor", ids[1])
	}
	if got := f.TrailingBytes(); len(got) != 7 || got[0] != 0xA0 || got[6] != 0xA6 {
		t.Errorf("TrailingBytes = % 02X, want A0..A6", got)
	}
}

func TestDecodeFormationRejectsShortInput(t *testing.T) {
	for _, n := range []int{0, 1, FormationSize - 1} {
		if _, err := DecodeFormation(make([]byte, n)); err == nil {
			t.Errorf("DecodeFormation with %d bytes should fail", n)
		}
	}
}

func TestFormationAt(t *testing.T) {
	table := make([]byte, FormationSize*3)
	table[FormationSize*2+2] = 0x2A // record 2, first monster id
	f, err := FormationAt(table, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got := f.OccupiedIDs(); len(got) != 6 || got[0] != 0x2A {
		t.Errorf("record 2 first id = %v, want 0x2A leading", got)
	}
	for _, id := range []int{-1, 3, 999} {
		if _, err := FormationAt(table, id); err == nil {
			t.Errorf("FormationAt(id=%d) on a 3-record table should fail", id)
		}
	}
}

// TestAllSlotsEmptyYieldsNoIDs covers the degenerate record, so a caller
// iterating OccupiedIDs never sees a phantom monster 255.
func TestAllSlotsEmptyYieldsNoIDs(t *testing.T) {
	raw := make([]byte, FormationSize)
	for i := 2; i < 2+FormationMonsterSlots; i++ {
		raw[i] = EmptySlot
	}
	f, err := DecodeFormation(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got := f.OccupiedIDs(); len(got) != 0 {
		t.Errorf("OccupiedIDs = %v, want empty", got)
	}
}

func TestDecodeMonster(t *testing.T) {
	raw := make([]byte, MonsterSize)
	raw[0x01] = 13
	raw[0x08], raw[0x09] = 0x28, 0x00 // HP 40
	raw[0x0A], raw[0x0B] = 0x0F, 0x00 // MP 15
	m, err := DecodeMonster(raw)
	if err != nil {
		t.Fatal(err)
	}
	if m.Power() != 13 || m.HP() != 40 || m.MP() != 15 {
		t.Errorf("power=%d HP=%d MP=%d, want 13/40/15", m.Power(), m.HP(), m.MP())
	}
	// Both stat fields are 16-bit little-endian; a byte-wide read would pass
	// the case above and fail this one.
	raw[0x09] = 0x01 // HP 0x0128 = 296
	m, _ = DecodeMonster(raw)
	if m.HP() != 296 {
		t.Errorf("HP = %d, want 296 (little-endian 16-bit)", m.HP())
	}
}

func TestDecodeMonsterRejectsShortInput(t *testing.T) {
	for _, n := range []int{0, MonsterSize - 1} {
		if _, err := DecodeMonster(make([]byte, n)); err == nil {
			t.Errorf("DecodeMonster with %d bytes should fail", n)
		}
	}
}

func TestMonsterAt(t *testing.T) {
	table := make([]byte, MonsterSize*2)
	table[MonsterSize+0x01] = 99
	m, err := MonsterAt(table, 1)
	if err != nil {
		t.Fatal(err)
	}
	if m.Power() != 99 {
		t.Errorf("record 1 power = %d, want 99", m.Power())
	}
	for _, id := range []int{-1, 2} {
		if _, err := MonsterAt(table, id); err == nil {
			t.Errorf("MonsterAt(id=%d) on a 2-record table should fail", id)
		}
	}
}

// TestVerifiedExtentMatchesTheRecords pins the boundary between what an
// experiment has checked and what is merely archived. The archive holds 576
// formations and 384 monsters because the true lengths are unknown; only the
// spans below carry live verification.
func TestVerifiedExtentMatchesTheRecords(t *testing.T) {
	if !FormationVerified(44) {
		t.Error("formation 44 is the mines encounter EXP-0030 verified")
	}
	if FormationVerified(45) || FormationVerified(432) {
		t.Error("formations past 44 are archived but unverified; Whelk (432) especially")
	}
	if !MonsterVerified(78) || MonsterVerified(79) {
		t.Error("monsters are cross-checked through record 78 (EXP-0028/0030)")
	}
	if FormationVerified(-1) || MonsterVerified(-1) {
		t.Error("negative ids are not verified")
	}
}

func TestNewTable(t *testing.T) {
	if _, err := NewTable(make([]byte, FormationSize*4), FormationSize); err != nil {
		t.Errorf("a whole number of records should be accepted: %v", err)
	}
	// A blob that is not a whole number of records means the stride or the
	// span is wrong, which is exactly the class of error that put the HUD
	// font on the wrong ROM range.
	if _, err := NewTable(make([]byte, FormationSize*4+1), FormationSize); err == nil {
		t.Error("a partial trailing record should be rejected")
	}
	if _, err := NewTable(nil, 0); err == nil {
		t.Error("a zero stride should be rejected")
	}
}

func TestTableRawIsACopy(t *testing.T) {
	tb, err := NewTable(make([]byte, MonsterSize*2), MonsterSize)
	if err != nil {
		t.Fatal(err)
	}
	if tb.Count() != 2 {
		t.Fatalf("Count = %d, want 2", tb.Count())
	}
	raw, err := tb.Raw(0)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] = 0xFF
	again, _ := tb.Raw(0)
	if again[0] != 0 {
		t.Error("Raw must return a copy; mutating it corrupted the table")
	}
	if _, err := tb.Raw(2); err == nil {
		t.Error("Raw past the end should fail")
	}
}

func FuzzDecodeFormation(f *testing.F) {
	f.Add(syntheticFormation())
	f.Add(make([]byte, 3))
	f.Fuzz(func(t *testing.T, b []byte) {
		rec, err := DecodeFormation(b)
		if err != nil {
			return
		}
		// Accessors must not panic or over-report for any byte pattern.
		if got := len(rec.OccupiedIDs()); got > FormationMonsterSlots {
			t.Fatalf("OccupiedIDs returned %d ids, max is %d", got, FormationMonsterSlots)
		}
		_ = rec.LeadingWord()
		if got := len(rec.TrailingBytes()); got != FormationSize-8 {
			t.Fatalf("TrailingBytes length %d", got)
		}
	})
}

func FuzzDecodeMonster(f *testing.F) {
	f.Add(make([]byte, MonsterSize))
	f.Fuzz(func(t *testing.T, b []byte) {
		m, err := DecodeMonster(b)
		if err != nil {
			return
		}
		_, _, _ = m.Power(), m.HP(), m.MP()
	})
}
