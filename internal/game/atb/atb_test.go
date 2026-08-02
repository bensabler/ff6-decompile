package atb

import "testing"

// TestStrideIsTwoNotTwenty guards the trap this layer is most likely to fall
// into: the ATB arrays are stride 2, while the HP/stat arrays of DISC-0001
// are stride $14.
func TestStrideIsTwoNotTwenty(t *testing.T) {
	if Stride != 2 {
		t.Fatalf("Stride = %d, want 2 (EXP-0043)", Stride)
	}
	if got := SlotIndex(9); got != 18 {
		t.Errorf("SlotIndex(9) = %d, want 18 ($12, the value ROMCPU:$C21141 loads into X)", got)
	}
	if got := SlotIndex(0); got != 0 {
		t.Errorf("SlotIndex(0) = %d, want 0", got)
	}
}

// TestPartySlotBoundaryMatchesTheBranch pins the party/enemy split to the
// comparison the ROM actually makes: CPX #$08 / BCC at ROMCPU:$C209F6, where
// X is the stride-2 index.
func TestPartySlotBoundaryMatchesTheBranch(t *testing.T) {
	for slot := 0; slot < SlotCount; slot++ {
		wantParty := SlotIndex(slot) < 0x08
		if got := IsPartySlot(slot); got != wantParty {
			t.Errorf("IsPartySlot(%d) = %v, want %v (X = $%02X)", slot, got, wantParty, SlotIndex(slot))
		}
	}
	if IsPartySlot(-1) || IsPartySlot(SlotCount) {
		t.Error("out-of-range slots are not party slots")
	}
}

func TestDecodeConfig(t *testing.T) {
	tests := []struct {
		name    string
		b       uint8
		mode    Mode
		speed   uint8
		display int
	}{
		// The project's recorded default (CURRENT_FOCUS, EXP-0041).
		{"default $2A", 0x2A, Wait, 2, 3},
		// EXP-0047's recorded starting state: Active, Battle Speed 6.
		{"EXP-0047 $25", 0x25, Active, 5, 6},
		{"all clear", 0x00, Active, 0, 1},
		{"speed max, active", 0x07, Active, 7, 8},
		{"wait only", 0x08, Wait, 0, 1},
		// Bits 4-7 are Msg.Speed and Cmd.Set; they must not leak in.
		{"upper bits ignored", 0xF0, Active, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := DecodeConfig(tt.b)
			if c.Mode != tt.mode {
				t.Errorf("Mode = %v, want %v", c.Mode, tt.mode)
			}
			if c.Speed != tt.speed {
				t.Errorf("Speed = %d, want %d", c.Speed, tt.speed)
			}
			if got := c.DisplaySpeed(); got != tt.display {
				t.Errorf("DisplaySpeed() = %d, want %d", got, tt.display)
			}
		})
	}
}

// TestSpeedByte pins 255 - 24*speed, including the two values EXP-0043
// measured enemy increments at.
func TestSpeedByte(t *testing.T) {
	want := map[uint8]uint8{0: 255, 1: 231, 2: 207, 3: 183, 4: 159, 5: 135}
	for speed, w := range want {
		if got := (Config{Speed: speed}).SpeedByte(); got != w {
			t.Errorf("SpeedByte(speed=%d) = %d, want %d", speed, got, w)
		}
	}
}

func TestSampleEntry(t *testing.T) {
	e := SampleEntry(Config{Mode: Wait, Speed: 3})
	if e.WaitFlag != 1 {
		t.Errorf("WaitFlag = %d, want 1", e.WaitFlag)
	}
	if e.SpeedByte != 183 {
		t.Errorf("SpeedByte = %d, want 183", e.SpeedByte)
	}
	if a := SampleEntry(Config{Mode: Active, Speed: 3}); a.WaitFlag != 0 {
		t.Errorf("Active WaitFlag = %d, want 0", a.WaitFlag)
	}
}

// TestPausedIsAnAnd is the ACTIVE/WAIT gate, ROMCPU:$C21124.
// The pause requires *both* a submenu and Wait mode — which is why the
// command window does not pause under WAIT, and why EXP-0040's "queued
// actions resolved out of order" was real system behaviour.
func TestPausedIsAnAnd(t *testing.T) {
	tests := []struct {
		name          string
		submenu, wait uint8
		want          bool
	}{
		{"active, no submenu", 0, 0, false},
		{"active, submenu open", 1, 0, false},
		{"wait, no submenu", 0, 1, false},
		{"wait, submenu open", 1, 1, true},
		// It is an AND of bits, not of booleans.
		{"disjoint bits", 0x02, 0x01, false},
		{"overlapping bits", 0x03, 0x01, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Paused(tt.submenu, tt.wait); got != tt.want {
				t.Errorf("Paused(%d,%d) = %v, want %v", tt.submenu, tt.wait, got, tt.want)
			}
		})
	}
}

// TestAdvanceSlotMatchesTheMeasuredIncrement reproduces EXP-0043's headline
// number: a slot with increment $9C advanced by exactly $4E per frame.
func TestAdvanceSlotMatchesTheMeasuredIncrement(t *testing.T) {
	var s State
	s.Increments[6] = 0x9C
	s.Flags[6] = FlagActive

	for i := 1; i <= 3; i++ {
		filled, err := s.AdvanceSlot(6)
		if err != nil {
			t.Fatal(err)
		}
		if filled {
			t.Fatalf("gauge should not fill at frame %d", i)
		}
		if want := uint16(0x4E * i); s.Gauges[6] != want {
			t.Fatalf("frame %d gauge = $%04X, want $%04X", i, s.Gauges[6], want)
		}
	}
}

// TestAdvanceSlotWraps covers the case EXP-0043 caught live: slot 8 going
// $00B6 -> $0004. That is a 16-bit wrap, not a clamp.
func TestAdvanceSlotWraps(t *testing.T) {
	var s State
	s.Gauges[8] = 0xFFB6
	s.Increments[8] = 0x9C // >>1 = $4E; $FFB6 + $4E = $10004
	s.Flags[8] = FlagActive

	filled, err := s.AdvanceSlot(8)
	if err != nil {
		t.Fatal(err)
	}
	if !filled {
		t.Error("the carry out of the 16-bit add signals a full gauge")
	}
	if s.Gauges[8] != 0x0004 {
		t.Errorf("gauge = $%04X, want $0004 (wrapped, not clamped)", s.Gauges[8])
	}
	if !s.Ready(8) {
		t.Error("filling should set FlagReady, as ROMCPU:$C211B2 ORAs #$20")
	}
}

// TestInactiveSlotIsSkipped covers the LSR / BCC at ROMCPU:$C2114E.
func TestInactiveSlotIsSkipped(t *testing.T) {
	var s State
	s.Increments[0] = 0xFFFF
	s.Flags[0] = 0 // bit 0 clear
	filled, err := s.AdvanceSlot(0)
	if err != nil {
		t.Fatal(err)
	}
	if filled || s.Gauges[0] != 0 {
		t.Errorf("an inactive slot must not advance: gauge=$%04X filled=%v", s.Gauges[0], filled)
	}
}

func TestAdvanceSlotRejectsBadSlots(t *testing.T) {
	var s State
	for _, slot := range []int{-1, SlotCount, 999} {
		if _, err := s.AdvanceSlot(slot); err == nil {
			t.Errorf("AdvanceSlot(%d) should be an error", slot)
		}
	}
}

// TestFrameFreezesEverythingWhenGated is EXP-0044's central finding: the gate
// skips the entire update, so the tick counter stops with the gauges.
func TestFrameFreezesEverythingWhenGated(t *testing.T) {
	var s State
	s.Entry = SampleEntry(Config{Mode: Wait, Speed: 3})
	s.Increments[6] = 0x9C
	s.Flags[6] = FlagActive

	// Submenu open under WAIT: nothing moves.
	if _, err := s.Frame(1); err != nil {
		t.Fatal(err)
	}
	if s.Tick != 0 || s.Gauges[6] != 0 {
		t.Errorf("gated frame advanced state: tick=%d gauge=$%04X", s.Tick, s.Gauges[6])
	}

	// Submenu closed: it resumes.
	if _, err := s.Frame(0); err != nil {
		t.Fatal(err)
	}
	if s.Tick != 1 || s.Gauges[6] != 0x4E {
		t.Errorf("ungated frame: tick=%d gauge=$%04X, want 1 and $004E", s.Tick, s.Gauges[6])
	}
}

// TestFrameUnderActiveNeverPauses is the other half of the matrix: with
// Bat.Mode = Active the AND is always zero, so a submenu changes nothing.
func TestFrameUnderActiveNeverPauses(t *testing.T) {
	var s State
	s.Entry = SampleEntry(Config{Mode: Active, Speed: 5})
	s.Increments[6] = 0x9C
	s.Flags[6] = FlagActive

	for i := 0; i < 3; i++ {
		if _, err := s.Frame(1); err != nil { // submenu open the whole time
			t.Fatal(err)
		}
	}
	if s.Tick != 3 {
		t.Errorf("tick = %d, want 3: ACTIVE mode never pauses", s.Tick)
	}
}

// TestFrameVisitsSlotsHighToLow pins the traversal order to the ROM's
// countdown loop (X = $12, DEX DEX, BPL at ROMCPU:$C21184).
func TestFrameVisitsSlotsHighToLow(t *testing.T) {
	var s State
	s.Entry = SampleEntry(Config{Mode: Active})
	// Make every slot fill on the first frame.
	for i := range s.Gauges {
		s.Gauges[i] = 0xFFFF
		s.Increments[i] = 0xFFFF
		s.Flags[i] = FlagActive
	}
	filled, err := s.Frame(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(filled) != SlotCount {
		t.Fatalf("filled %d slots, want %d", len(filled), SlotCount)
	}
	for i, slot := range filled {
		if want := SlotCount - 1 - i; slot != want {
			t.Fatalf("filled[%d] = slot %d, want %d (high to low)", i, slot, want)
		}
	}
}

// TestBattleSpeedScalesEnemiesOnly encodes EXP-0043's measurement: party
// increments were byte-identical at Battle Speed 3 and 6 while enemy
// increments went 240 -> 156.
//
// The increment formula itself is not implemented — ROMCPU:$C209F0 is only
// partly decoded — so this asserts the property the branch guarantees: the
// speed byte applies to enemy indices and not party ones.
func TestBattleSpeedScalesEnemiesOnly(t *testing.T) {
	for slot := 0; slot < SlotCount; slot++ {
		x := SlotIndex(slot)
		scaled := x >= 0x08 // CPX #$08 / BCC skips the multiply below this
		if scaled == IsPartySlot(slot) {
			t.Errorf("slot %d (X=$%02X): scaled=%v but IsPartySlot=%v; these must be opposites",
				slot, x, scaled, IsPartySlot(slot))
		}
	}
}

func TestModeString(t *testing.T) {
	if Active.String() != "Active" || Wait.String() != "Wait" {
		t.Error("Mode.String must name the setting as the Config screen does")
	}
}

// FuzzAdvanceSlot checks the 16-bit arithmetic never panics and never
// produces a value outside a uint16.
func FuzzAdvanceSlot(f *testing.F) {
	f.Add(uint16(0), uint16(0x9C), uint8(FlagActive))
	f.Add(uint16(0xFFFF), uint16(0xFFFF), uint8(FlagActive))
	f.Fuzz(func(t *testing.T, gauge, inc uint16, flags uint8) {
		var s State
		s.Gauges[3] = gauge
		s.Increments[3] = inc
		s.Flags[3] = flags
		filled, err := s.AdvanceSlot(3)
		if err != nil {
			t.Fatalf("slot 3 is valid: %v", err)
		}
		if flags&FlagActive == 0 {
			if s.Gauges[3] != gauge || filled {
				t.Fatal("an inactive slot must not change")
			}
			return
		}
		want := gauge + inc>>1 // Go's uint16 add wraps, like the 65816's
		if s.Gauges[3] != want {
			t.Fatalf("gauge = $%04X, want $%04X", s.Gauges[3], want)
		}
		if filled != (uint32(gauge)+uint32(inc>>1) > 0xFFFF) {
			t.Fatalf("filled = %v, disagrees with the carry out", filled)
		}
	})
}
