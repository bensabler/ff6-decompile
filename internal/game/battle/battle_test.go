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
			p := BattleSlots{}
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
			p := BattleSlots{}
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
	p := BattleSlots{
		CurrentHP: [BattleSlotCount]uint16{10, 20, 30, 40, 24, 35},
		MaxHP:     [BattleSlotCount]uint16{99, 99, 99, 99, 99, 99},
	}
	p.ApplyHPDelta(0, -10)
	if p.CurrentHP != [BattleSlotCount]uint16{0, 20, 30, 40, 24, 35} {
		t.Errorf("unexpected cross-slot mutation: %v", p.CurrentHP)
	}
}

// Enemy entries (4-9) go through the identical engine paths — observed
// live in EXP-0005 (enemy HP at entries 4-5 damaged and zeroed by the
// same stores that serve the party).
func TestApplyHPDeltaEnemySlots(t *testing.T) {
	p := BattleSlots{}
	p.CurrentHP[5] = 35 // EXP-0005 enemy entry 5 initial value
	p.MaxHP[5] = 35

	if death := p.ApplyHPDelta(5, -10); death {
		t.Fatal("partial damage must not fire the death event")
	}
	if p.CurrentHP[5] != 25 {
		t.Fatalf("CurrentHP[5] = %d, want 25", p.CurrentHP[5])
	}
	if death := p.ApplyHPDelta(5, -25); !death {
		t.Fatal("reaching zero must fire the death event for enemy entries")
	}
	if p.CurrentHP[5] != 0 {
		t.Fatalf("CurrentHP[5] = %d, want 0", p.CurrentHP[5])
	}
}

func TestAccumulatePending(t *testing.T) {
	tests := []struct {
		name            string
		current, amount uint16
		want            uint16
	}{
		{name: "sentinel starts from zero", current: PendingDeltaNone, amount: 4, want: 4},
		{name: "accumulates onto existing pending", current: 100, amount: 250, want: 350},
		{name: "caps at 9999", current: 9000, amount: 1500, want: PendingDeltaCap},
		{name: "exactly 9999 stays", current: 9998, amount: 1, want: PendingDeltaCap},
		{name: "exactly 10000 clamps", current: 9999, amount: 1, want: PendingDeltaCap},
		{name: "16-bit overflow clamps", current: 0xFFFE, amount: 0xFFFE, want: PendingDeltaCap},
		{name: "sentinel plus zero is zero", current: PendingDeltaNone, amount: 0, want: 0},
		{name: "EXP-0004 live capture value", current: PendingDeltaNone, amount: 4, want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AccumulatePending(tt.current, tt.amount); got != tt.want {
				t.Errorf("AccumulatePending(%#x, %d) = %d, want %d", tt.current, tt.amount, got, tt.want)
			}
		})
	}
}

func TestApplyElementResponse(t *testing.T) {
	r := ElementResponse{FlipToHeal: 0x01, Nullify: 0x02, Halve: 0x04, Double: 0x08}
	tests := []struct {
		name       string
		amount     uint16
		heal       bool
		elems      uint8
		nullified  uint8
		wantAmount uint16
		wantHeal   bool
	}{
		{name: "non-elemental unchanged", amount: 346, elems: 0x00, wantAmount: 346},
		{name: "all elements nullified zeroes", amount: 346, elems: 0x10, nullified: 0xFF, wantAmount: 0},
		{name: "flip-to-heal match", amount: 100, elems: 0x01, wantAmount: 100, wantHeal: true},
		{name: "flip back to damage when already heal", amount: 100, heal: true, elems: 0x01, wantAmount: 100, wantHeal: false},
		{name: "nullify match zeroes", amount: 100, elems: 0x02, wantAmount: 0},
		{name: "halve match floors", amount: 101, elems: 0x04, wantAmount: 50},
		{name: "double match", amount: 100, elems: 0x08, wantAmount: 200},
		{name: "double overflow guard", amount: 0x8000, elems: 0x08, wantAmount: 0x8000},
		{name: "double just under guard", amount: 0x7FFF, elems: 0x08, wantAmount: 0xFFFE},
		{name: "flip wins over nullify", amount: 100, elems: 0x03, wantAmount: 100, wantHeal: true},
		{name: "nullify wins over halve", amount: 100, elems: 0x06, wantAmount: 0},
		{name: "halve wins over double", amount: 100, elems: 0x0C, wantAmount: 50},
		{name: "masks tested against full elems not nullify-masked", amount: 100, elems: 0x03, nullified: 0x01, wantAmount: 100, wantHeal: true},
		{name: "no mask match leaves amount", amount: 100, elems: 0x10, wantAmount: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmount, gotHeal := ApplyElementResponse(tt.amount, tt.heal, tt.elems, tt.nullified, r)
			if gotAmount != tt.wantAmount || gotHeal != tt.wantHeal {
				t.Errorf("ApplyElementResponse() = (%d, %v), want (%d, %v)", gotAmount, gotHeal, tt.wantAmount, tt.wantHeal)
			}
		})
	}
}

func TestScale256(t *testing.T) {
	tests := []struct {
		value  uint16
		factor uint8
		want   uint16
	}{
		{0, 255, 0},
		{256, 255, 255},
		{1000, 128, 500},
		{0xFFFF, 255, 0xFEFF},
		{346, 170, 229}, // the $AA (~2/3) path on the EXP-0007 Fire Beam value
		{12345, 0, 0},
	}
	for _, tt := range tests {
		if got := Scale256(tt.value, tt.factor); got != tt.want {
			t.Errorf("Scale256(%d, %d) = %d, want %d", tt.value, tt.factor, got, tt.want)
		}
		// Cross-check against the original's byte-composition form.
		alt := uint16(tt.factor)*uint16(tt.value>>8) + uint16(tt.factor)*uint16(tt.value&0xFF)>>8
		if alt != Scale256(tt.value, tt.factor) {
			t.Errorf("composition mismatch for (%d, %d): %d vs %d", tt.value, tt.factor, alt, Scale256(tt.value, tt.factor))
		}
	}
}

func FuzzScale256Composition(f *testing.F) {
	f.Add(uint16(346), uint8(170))
	f.Fuzz(func(t *testing.T, value uint16, factor uint8) {
		want := uint16(factor)*uint16(value>>8) + uint16(factor)*uint16(value&0xFF)>>8
		if got := Scale256(value, factor); got != want {
			t.Errorf("Scale256(%d, %d) = %d, byte-composition = %d", value, factor, got, want)
		}
	})
}

func TestApplyDefense(t *testing.T) {
	tests := []struct {
		amount  uint16
		defense uint8
		want    uint16
	}{
		{1000, NoDefense, 1000}, // $FF sentinel skips entirely (no +1)
		{1000, 0, 997},          // (1000*255)/256+1
		{1000, 128, 497},        // (1000*127)/256+1
		{1000, 254, 4},          // (1000*1)/256+1
		{0, 0, 1},               // minimum 1 after scaling
	}
	for _, tt := range tests {
		if got := ApplyDefense(tt.amount, tt.defense); got != tt.want {
			t.Errorf("ApplyDefense(%d, %d) = %d, want %d", tt.amount, tt.defense, got, tt.want)
		}
	}
}

func TestChainBoost(t *testing.T) {
	tests := []struct {
		amount uint16
		count  uint8
		want   uint16
	}{
		{100, 0, 100},
		{100, 1, 150},
		{100, 2, 225},
		{101, 1, 151},       // odd: 101 + 50
		{0xFFFF, 1, 0xFFFF}, // clamp
		{0xAAAB, 1, 0xFFFF}, // 0xAAAB + 0x5555 = 0x10000 -> clamp
		{0xAAAA, 1, 0xFFFF}, // 0xAAAA + 0x5555 = 0xFFFF exactly
	}
	for _, tt := range tests {
		if got := ChainBoost(tt.amount, tt.count); got != tt.want {
			t.Errorf("ChainBoost(%d, %d) = %d, want %d", tt.amount, tt.count, got, tt.want)
		}
	}
}

func TestBaseAmountStandard(t *testing.T) {
	tests := []struct {
		name                string
		power, statA, statB uint8
		want                uint16
	}{
		{name: "EXP-0015 live golden vector", power: 60, statA: 28, statB: 4, want: 450},
		{name: "statB zero skips the x4", power: 60, statA: 28, statB: 0, want: 60},
		{name: "zero power", power: 0, statA: 100, statB: 100, want: 0},
		{name: "max operands wrap at 16 bits", power: 255, statA: 255, statB: 255,
			// 255*4 + ((255*255*255)>>5)&0xFFFF, the original's 24-bit
			// product shifted with only 16 result bits retained
			want: uint16(1020 + ((255 * 255 * 255 >> 5) & 0xFFFF)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseAmountStandard(tt.power, tt.statA, tt.statB); got != tt.want {
				t.Errorf("BaseAmountStandard(%d, %d, %d) = %d, want %d", tt.power, tt.statA, tt.statB, got, tt.want)
			}
		})
	}
}

// The full damage chain on the EXP-0007/0014/0015 live evidence:
// base 450 (power 60, stats 28/4) through defense ~58 gives the
// observed 346 Fire Beam damage.
func TestDamageChainGolden(t *testing.T) {
	base := BaseAmountStandard(60, 28, 4)
	if base != 450 {
		t.Fatalf("base = %d, want 450", base)
	}
	got := ApplyDefense(base, 58)
	if got != 347 && got != 346 {
		t.Fatalf("ApplyDefense(450, 58) = %d, want ~346", got)
	}
}
