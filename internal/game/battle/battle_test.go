package battle

import "testing"

func TestApplyHPDelta(t *testing.T) {
	tests := []struct {
		name      string
		hp, max   uint16
		status    uint16
		delta     int32
		wantHP    uint16
		wantDeath bool
	}{
		{name: "damage partial", hp: 100, max: 100, delta: -30, wantHP: 70},
		{name: "damage to exactly zero fires death", hp: 30, max: 100, delta: -30, wantHP: 0, wantDeath: true},
		{name: "damage underflow floors at zero and fires death", hp: 30, max: 100, delta: -31, wantHP: 0, wantDeath: true},
		{name: "status bit 1 suppresses death event", hp: 30, max: 100, status: 0x0002, delta: -30, wantHP: 0, wantDeath: false},
		{name: "other status bits do not suppress", hp: 30, max: 100, status: 0xFFFD, delta: -30, wantHP: 0, wantDeath: true},
		{name: "heal partial", hp: 50, max: 100, delta: 20, wantHP: 70},
		{name: "heal to exactly max", hp: 50, max: 100, delta: 50, wantHP: 100},
		{name: "heal beyond max clamps", hp: 50, max: 100, delta: 51, wantHP: 100},
		{name: "heal 16-bit overflow clamps to max", hp: 0xFFF0, max: 0xFFFA, delta: 0x0100, wantHP: 0xFFFA},
		{name: "zero delta heals by nothing", hp: 50, max: 100, delta: 0, wantHP: 50},
		{name: "damage magnitude clamps to 16 bits", hp: 0xFFFF, max: 0xFFFF, delta: -0x20000, wantHP: 0, wantDeath: true},
		{name: "heal magnitude clamps to 16 bits", hp: 1, max: 0xFFFF, delta: 0x20000, wantHP: 0xFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PartySlots{}
			p.CurrentHP[1] = tt.hp
			p.MaxHP[1] = tt.max
			p.UnknownStatus3EE4[1] = tt.status

			gotDeath := p.ApplyHPDelta(1, tt.delta)

			if p.CurrentHP[1] != tt.wantHP {
				t.Errorf("CurrentHP = %d, want %d", p.CurrentHP[1], tt.wantHP)
			}
			if gotDeath != tt.wantDeath {
				t.Errorf("death event = %v, want %v", gotDeath, tt.wantDeath)
			}
		})
	}
}

func TestApplyMPDelta(t *testing.T) {
	tests := []struct {
		name         string
		mp, maxMP    uint16
		hp           uint16
		diesAtZeroMP bool
		status       uint16
		delta        int32
		wantMP       uint16
		wantHP       uint16
		wantDeath    bool
	}{
		{name: "spend partial", mp: 24, maxMP: 24, hp: 50, delta: -5, wantMP: 19, wantHP: 50},
		{name: "MP to zero without flag is harmless", mp: 5, maxMP: 24, hp: 50, delta: -5, wantMP: 0, wantHP: 50},
		{name: "MP to zero with flag zeroes HP and fires death", mp: 5, maxMP: 24, hp: 50, diesAtZeroMP: true, delta: -5, wantMP: 0, wantHP: 0, wantDeath: true},
		{name: "MP underflow with flag fires death", mp: 5, maxMP: 24, hp: 50, diesAtZeroMP: true, delta: -6, wantMP: 0, wantHP: 0, wantDeath: true},
		{name: "MP death suppressed by status bit 1 still zeroes HP", mp: 5, maxMP: 24, hp: 50, diesAtZeroMP: true, status: 0x0002, delta: -5, wantMP: 0, wantHP: 0, wantDeath: false},
		{name: "MP heal clamps to max", mp: 20, maxMP: 24, hp: 50, delta: 10, wantMP: 24, wantHP: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PartySlots{}
			p.CurrentMP[2] = tt.mp
			p.MaxMP[2] = tt.maxMP
			p.CurrentHP[2] = tt.hp
			p.DiesAtZeroMP[2] = tt.diesAtZeroMP
			p.UnknownStatus3EE4[2] = tt.status

			gotDeath := p.ApplyMPDelta(2, tt.delta)

			if p.CurrentMP[2] != tt.wantMP {
				t.Errorf("CurrentMP = %d, want %d", p.CurrentMP[2], tt.wantMP)
			}
			if p.CurrentHP[2] != tt.wantHP {
				t.Errorf("CurrentHP = %d, want %d", p.CurrentHP[2], tt.wantHP)
			}
			if gotDeath != tt.wantDeath {
				t.Errorf("death event = %v, want %v", gotDeath, tt.wantDeath)
			}
		})
	}
}

// Other slots must never be touched by a delta applied to one slot.
func TestApplyHPDeltaIsolation(t *testing.T) {
	p := PartySlots{
		CurrentHP: [PartySlotCount]uint16{10, 20, 30, 40},
		MaxHP:     [PartySlotCount]uint16{99, 99, 99, 99},
	}
	p.ApplyHPDelta(0, -10)
	if p.CurrentHP != [PartySlotCount]uint16{0, 20, 30, 40} {
		t.Errorf("unexpected cross-slot mutation: %v", p.CurrentHP)
	}
}
