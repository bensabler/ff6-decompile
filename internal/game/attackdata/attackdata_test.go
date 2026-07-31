package attackdata

import "testing"

// Synthetic fixture only — no ROM bytes.
func syntheticRecord() []byte {
	return []byte{0xAA, 0x01, 0x81, 0x80, 0x82, 0x55, 0x3C, 0x66, 0x77, 0x88, 0x82, 0x99, 0xBB, 0xCC}
}

func TestDecodeAndAccessors(t *testing.T) {
	r, err := Decode(syntheticRecord())
	if err != nil {
		t.Fatal(err)
	}
	if r.Element() != 0x01 {
		t.Errorf("Element = %#x, want 0x01", r.Element())
	}
	if r.Flags2() != 0x81 || !r.PhysicalFormula() {
		t.Errorf("Flags2 = %#x, PhysicalFormula = %v", r.Flags2(), r.PhysicalFormula())
	}
	if !r.TargetsMP() {
		t.Error("TargetsMP = false, want true")
	}
	if r.Mode() != 0x82 {
		t.Errorf("Mode = %#x, want 0x82", r.Mode())
	}
	if r.Power() != 0x3C {
		t.Errorf("Power = %#x, want 0x3C (the live Fire Beam power 60)", r.Power())
	}
	if r.AbortBits() != 0x82 {
		t.Errorf("AbortBits = %#x, want 0x82", r.AbortBits())
	}
}

func TestDecodeShortBuffer(t *testing.T) {
	if _, err := Decode(make([]byte, RecordSize-1)); err == nil {
		t.Error("want error for short buffer")
	}
}

func TestRecordAt(t *testing.T) {
	table := make([]byte, RecordSize*3)
	copy(table[RecordSize:], syntheticRecord())
	r, err := RecordAt(table, 1)
	if err != nil {
		t.Fatal(err)
	}
	if r.Power() != 0x3C {
		t.Errorf("Power = %#x, want 0x3C", r.Power())
	}
	if _, err := RecordAt(table, 3); err == nil {
		t.Error("want out-of-range error for index 3")
	}
	if _, err := RecordAt(table, -1); err == nil {
		t.Error("want error for negative index")
	}
}

func FuzzDecode(f *testing.F) {
	f.Add(syntheticRecord())
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, b []byte) {
		r, err := Decode(b)
		if err == nil && len(b) < RecordSize {
			t.Fatal("Decode accepted a short buffer")
		}
		_ = r
	})
}
