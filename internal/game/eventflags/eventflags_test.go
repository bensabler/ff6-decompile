package eventflags

import "testing"

// Golden vectors from CONTRA-0002 (static decode of ROMCPU:$C0BAED and
// the mask tables) and the EXP-0037 live write inventory.
func TestRefDerivation(t *testing.T) {
	cases := []struct {
		name     string
		ref      Ref
		byteAddr uint16
		bit      uint8
		setMask  byte
	}{
		// The CONTRA-0002 byte: $1EA5 holds flags $28-$2F of ArrayB.
		{"1EA5 bit0", Ref{ArrayB, 0x28}, 0x1EA5, 0, 0x01},
		{"1EA5 bit2", Ref{ArrayB, 0x2A}, 0x1EA5, 2, 0x04},
		{"1EA5 bit3", Ref{ArrayB, 0x2B}, 0x1EA5, 3, 0x08},
		// EXP-0037 live-observed story flags in the other arrays.
		{"1EA1 bit3", Ref{ArrayB, 0x0B}, 0x1EA1, 3, 0x08},
		{"1EDC bit0", Ref{ArrayC, 0xE0}, 0x1EDC, 0, 0x01},
		{"1EDF bit7", Ref{ArrayC, 0xFF}, 0x1EDF, 7, 0x80},
		// Array edges.
		{"first flag", Ref{ArrayA, 0x00}, 0x1E80, 0, 0x01},
		{"last flag", Ref{ArrayC, 0xFF}, 0x1EDF, 7, 0x80},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ref.ByteAddr(); got != c.byteAddr {
				t.Errorf("ByteAddr = $%04X, want $%04X", got, c.byteAddr)
			}
			if got := c.ref.Bit(); got != c.bit {
				t.Errorf("Bit = %d, want %d", got, c.bit)
			}
			if got := c.ref.SetMask(); got != c.setMask {
				t.Errorf("SetMask = $%02X, want $%02X", got, c.setMask)
			}
			if got, want := c.ref.ClearMask(), ^c.setMask; got != want {
				t.Errorf("ClearMask = $%02X, want $%02X", got, want)
			}
		})
	}
}

// The mask tables must stay byte-identical to the ROM tables they
// mirror (ROMCPU:$C0BAFC / $C0BB04, read in CONTRA-0002).
func TestMaskTablesMatchROMDecode(t *testing.T) {
	wantSet := [8]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
	wantClear := [8]byte{0xFE, 0xFD, 0xFB, 0xF7, 0xEF, 0xDF, 0xBF, 0x7F}
	if SetMasks != wantSet {
		t.Errorf("SetMasks = %#v, want %#v", SetMasks, wantSet)
	}
	if ClearMasks != wantClear {
		t.Errorf("ClearMasks = %#v, want %#v", ClearMasks, wantClear)
	}
	for i := range SetMasks {
		if SetMasks[i] != ^ClearMasks[i] {
			t.Errorf("mask %d: set $%02X is not the complement of clear $%02X", i, SetMasks[i], ClearMasks[i])
		}
	}
}

func TestFlagAtRoundTrip(t *testing.T) {
	for _, base := range []ArrayBase{ArrayA, ArrayB, ArrayC} {
		for f := 0; f < 256; f++ {
			ref := Ref{Array: base, Flag: uint8(f)}
			got, err := FlagAt(ref.ByteAddr(), ref.Bit())
			if err != nil {
				t.Fatalf("FlagAt($%04X,%d): %v", ref.ByteAddr(), ref.Bit(), err)
			}
			if got != ref {
				t.Fatalf("FlagAt($%04X,%d) = %+v, want %+v", ref.ByteAddr(), ref.Bit(), got, ref)
			}
		}
	}
}

func TestFlagAtRejectsOutside(t *testing.T) {
	for _, addr := range []uint16{0x1E7F, 0x1EE0, 0x0000, 0xFFFF} {
		if _, err := FlagAt(addr, 0); err == nil {
			t.Errorf("FlagAt($%04X,0) accepted an address outside the arrays", addr)
		}
	}
	if _, err := FlagAt(0x1E80, 8); err == nil {
		t.Error("FlagAt accepted bit 8")
	}
}

// The milestone-05 value of $1EA5 is $0D = flags $28, $2A, $2B set
// (CONTRA-0002; reproduced live by EXP-0037 at frames 34298, 39090,
// 50699).
func TestMilestone05Value(t *testing.T) {
	array := make([]byte, 0x20)
	for _, f := range []uint8{0x28, 0x2A, 0x2B} {
		r := Ref{ArrayB, f}
		array[r.ByteOffset()] |= r.SetMask()
	}
	if array[5] != 0x0D {
		t.Errorf("$1EA5 = $%02X after setting $28/$2A/$2B, want $0D", array[5])
	}
	for _, c := range []struct {
		flag uint8
		want bool
	}{{0x28, true}, {0x29, false}, {0x2A, true}, {0x2B, true}, {0x2C, false}} {
		set, err := (Ref{ArrayB, c.flag}).IsSet(array)
		if err != nil {
			t.Fatal(err)
		}
		if set != c.want {
			t.Errorf("flag $%02X set = %v, want %v", c.flag, set, c.want)
		}
	}
}

func TestID(t *testing.T) {
	if got := (Ref{ArrayB, 0x2B}).ID(); got != "EVF-1EA0-$2B" {
		t.Errorf("ID = %q", got)
	}
	if got := (Ref{ArrayA, 0x00}).ID(); got != "EVF-1E80-$00" {
		t.Errorf("ID = %q", got)
	}
}

func TestIsSetBounds(t *testing.T) {
	if _, err := (Ref{ArrayA, 0xFF}).IsSet(make([]byte, 4)); err == nil {
		t.Error("IsSet accepted a short snapshot")
	}
}
