package snesdma

import (
	"reflect"
	"testing"
)

// regs builds a channel register block, $43n0 through $43nA.
func regs(ctrl, dest, srcLo, srcHi, srcBank, sizeLo, sizeHi byte) []byte {
	return []byte{ctrl, dest, srcLo, srcHi, srcBank, sizeLo, sizeHi, 0, 0, 0, 0}
}

func TestDecodeControlBits(t *testing.T) {
	tests := []struct {
		name  string
		ctrl  byte
		mode  uint8
		fixed bool
		dec   bool
		dir   Direction
		hdmaI bool
	}{
		{"mode 0, A->B, stepping", 0x00, 0, false, false, AToB, false},
		{"mode 1, the VRAM upload pattern", 0x01, 1, false, false, AToB, false},
		{"fixed source", 0x08, 0, true, false, AToB, false},
		{"decrementing source", 0x10, 0, false, true, AToB, false},
		{"HDMA indirect", 0x40, 0, false, false, AToB, true},
		{"B->A readback", 0x80, 0, false, false, BToA, false},
		{"everything at once", 0xDF, 7, true, true, BToA, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Decode(0, regs(tt.ctrl, 0x18, 0, 0, 0, 0, 0))
			if err != nil {
				t.Fatal(err)
			}
			if c.TransferMode != tt.mode || c.Fixed != tt.fixed || c.Decrement != tt.dec ||
				c.Direction != tt.dir || c.HDMAIndirect != tt.hdmaI {
				t.Errorf("ctrl $%02X decoded to mode=%d fixed=%v dec=%v dir=%v hdmaI=%v",
					tt.ctrl, c.TransferMode, c.Fixed, c.Decrement, c.Direction, c.HDMAIndirect)
			}
		})
	}
}

func TestDecodeAddressesAndSize(t *testing.T) {
	c, err := Decode(3, regs(0x01, 0x18, 0x80, 0xC1, 0x7E, 0x00, 0x20))
	if err != nil {
		t.Fatal(err)
	}
	if c.Index != 3 {
		t.Errorf("Index = %d, want 3", c.Index)
	}
	if got := c.Source(); got != 0x7EC180 {
		t.Errorf("Source = $%06X, want $7EC180", got)
	}
	if c.Dest != 0x18 || !c.TargetsVRAM() {
		t.Errorf("Dest = $%02X, TargetsVRAM = %v", c.Dest, c.TargetsVRAM())
	}
	if got := c.ByteCount(); got != 0x2000 {
		t.Errorf("ByteCount = %d, want %d", got, 0x2000)
	}
}

// The zero-size rule. A 64 KB VRAM upload has Size 0, and reading Size
// directly records it as an empty transfer — which is precisely how the mines
// savestate's channel 0 reads at a glance.
func TestZeroSizeMeans64K(t *testing.T) {
	c, _ := Decode(0, regs(0x01, 0x18, 0, 0, 0x7E, 0x00, 0x00))
	if got := c.ByteCount(); got != 65536 {
		t.Errorf("ByteCount for Size 0 = %d, want 65536", got)
	}
	if c.Size != 0 {
		t.Errorf("Size should still report the raw 0, got %d", c.Size)
	}
}

func TestTransferPattern(t *testing.T) {
	tests := []struct {
		mode uint8
		want []uint8
	}{
		{0, []uint8{0}},
		{1, []uint8{0, 1}},
		{2, []uint8{0, 0}},
		{3, []uint8{0, 0, 1, 1}},
		{4, []uint8{0, 1, 2, 3}},
		{5, []uint8{0, 1, 0, 1}},
		{6, []uint8{0, 0}},       // duplicates 2
		{7, []uint8{0, 0, 1, 1}}, // duplicates 3
	}
	for _, tt := range tests {
		if got := TransferPattern(tt.mode); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TransferPattern(%d) = %v, want %v", tt.mode, got, tt.want)
		}
	}
	// Bits above 2 must be masked, not passed through.
	if !reflect.DeepEqual(TransferPattern(0xF9), TransferPattern(1)) {
		t.Error("TransferPattern must mask the control byte to modes 0-7")
	}
}

func TestSourceSpan(t *testing.T) {
	tests := []struct {
		name             string
		ctrl             byte
		src              uint32
		size             uint16
		wantOK           bool
		wantStart, wantE uint32
	}{
		{"forward 8 KB", 0x01, 0xE08460, 0x2000, true, 0xE08460, 0xE0A45F},
		{"fixed source has no span", 0x09, 0xE08460, 0x2000, false, 0, 0},
		{"B->A has no source span", 0x81, 0xE08460, 0x2000, false, 0, 0},
		{"decrementing", 0x11, 0xE0A45F, 0x2000, true, 0xE08460, 0xE0A45F},
		{"decrement past the bank start is refused", 0x11, 0xE00010, 0x2000, false, 0, 0},
		{"forward past the bank end is refused", 0x01, 0xE0FF00, 0x0000, false, 0, 0},
		{"exactly filling a bank is fine", 0x01, 0xE00000, 0x0000, true, 0xE00000, 0xE0FFFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := Decode(0, regs(tt.ctrl, 0x18,
				byte(tt.src), byte(tt.src>>8), byte(tt.src>>16),
				byte(tt.size), byte(tt.size>>8)))
			start, end, ok := c.SourceSpan()
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v (start $%06X end $%06X)", ok, tt.wantOK, start, end)
			}
			if ok && (start != tt.wantStart || end != tt.wantE) {
				t.Errorf("span = $%06X-$%06X, want $%06X-$%06X", start, end, tt.wantStart, tt.wantE)
			}
		})
	}
}

func TestTargetClassification(t *testing.T) {
	tests := []struct {
		dest             uint8
		vram, cgram, oam bool
	}{
		{0x18, true, false, false},
		{0x19, true, false, false},
		{0x22, false, true, false},
		{0x04, false, false, true},
		{0x30, false, false, false}, // $2130, a colour-math register
	}
	for _, tt := range tests {
		c, _ := Decode(0, regs(0x01, tt.dest, 0, 0, 0, 0, 0))
		if c.TargetsVRAM() != tt.vram || c.TargetsCGRAM() != tt.cgram || c.TargetsOAM() != tt.oam {
			t.Errorf("dest $%02X: vram=%v cgram=%v oam=%v", tt.dest,
				c.TargetsVRAM(), c.TargetsCGRAM(), c.TargetsOAM())
		}
	}
}

func TestEnabled(t *testing.T) {
	tests := []struct {
		mask uint8
		want []int
	}{
		{0x00, nil},
		{0x01, []int{0}},
		{0x81, []int{0, 7}},
		{0xFF, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{0x0A, []int{1, 3}},
	}
	for _, tt := range tests {
		if got := Enabled(tt.mask); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Enabled($%02X) = %v, want %v", tt.mask, got, tt.want)
		}
	}
}

func TestDecodeRejectsBadInput(t *testing.T) {
	if _, err := Decode(-1, regs(0, 0, 0, 0, 0, 0, 0)); err == nil {
		t.Error("negative channel index should fail")
	}
	if _, err := Decode(ChannelCount, regs(0, 0, 0, 0, 0, 0, 0)); err == nil {
		t.Error("channel index past the last should fail")
	}
	if _, err := Decode(0, []byte{0, 1, 2}); err == nil {
		t.Error("a short register block should fail rather than read past it")
	}
	if _, err := Decode(0, nil); err == nil {
		t.Error("a nil register block should fail")
	}
}

func FuzzDecode(f *testing.F) {
	f.Add([]byte{0x01, 0x18, 0x80, 0xC1, 0x7E, 0x00, 0x20, 0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, b []byte) {
		for i := 0; i < ChannelCount; i++ {
			c, err := Decode(i, b)
			if err != nil {
				continue
			}
			if n := c.ByteCount(); n < 1 || n > 65536 {
				t.Fatalf("ByteCount %d out of range for regs %x", n, b)
			}
			if p := TransferPattern(c.TransferMode); len(p) == 0 || len(p) > 4 {
				t.Fatalf("TransferPattern(%d) = %v", c.TransferMode, p)
			}
			if start, end, ok := c.SourceSpan(); ok {
				if end < start {
					t.Fatalf("span $%06X-$%06X is inverted", start, end)
				}
				if end-start+1 != uint32(c.ByteCount()) {
					t.Fatalf("span $%06X-$%06X does not cover %d bytes", start, end, c.ByteCount())
				}
			}
		}
	})
}
