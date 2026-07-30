package chardata

import "testing"

func TestCopyCharacterFields(t *testing.T) {
	tests := []struct {
		name         string
		src          CharacterFieldsSource
		mode         CopyMode
		want         [PartySlotCount]PartySlotRecord
		wantSlotMask uint8
	}{
		{
			name: "session 002 live battle capture, both flags clear",
			// Captured from WRAM during the Narshe intro battle:
			// HP 53/39/33, max 63/68/70, MP 24/24 slot 0, $2E98=8,8,8,
			// $2EA0 all zero; $628D=$00, $E9EF=$00, $61AD=$0F.
			src: CharacterFieldsSource{
				CurrentHP:      [PartySlotCount]uint16{0x0035, 0x0027, 0x0021, 0},
				CandidateMaxHP: [PartySlotCount]uint16{0x003F, 0x0044, 0x0046, 0},
				Unknown04:      [PartySlotCount]uint16{0x0018, 0, 0, 0},
				Unknown06:      [PartySlotCount]uint16{0x0018, 0, 0, 0},
				Unknown08:      [PartySlotCount]uint16{0x0008, 0x0008, 0x0008, 0},
				Unknown0A:      [PartySlotCount]uint16{0, 0, 0, 0},
			},
			mode: CopyMode{},
			want: [PartySlotCount]PartySlotRecord{
				{CurrentHP: 0x0035, CandidateMaxHP: 0x003F, Unknown04: 0x0018, Unknown06: 0x0018, Unknown08: 0x0008},
				{CurrentHP: 0x0027, CandidateMaxHP: 0x0044, Unknown08: 0x0008},
				{CurrentHP: 0x0021, CandidateMaxHP: 0x0046, Unknown08: 0x0008},
				{},
			},
			wantSlotMask: 0x0F,
		},
		{
			name: "masked mode ANDs +08 with 0x0038 and zeroes +0A",
			src: CharacterFieldsSource{
				Unknown08: [PartySlotCount]uint16{0xFFFF, 0x0038, 0x0007, 0x1240},
				Unknown0A: [PartySlotCount]uint16{0xFFFF, 0x1234, 0, 0x2000},
			},
			mode: CopyMode{UnknownFlagE9EF: true},
			want: [PartySlotCount]PartySlotRecord{
				{Unknown08: 0x0038},
				{Unknown08: 0x0038},
				{Unknown08: 0x0000},
				{Unknown08: 0x0000},
			},
			// Mask derives from the unmasked source: bit 13 set in
			// slots 0 ($FFFF) and 3 ($2000) clears their bits.
			wantSlotMask: 0x06,
		},
		{
			name: "flag $628D set suppresses masking even with $E9EF set",
			src: CharacterFieldsSource{
				Unknown08: [PartySlotCount]uint16{0xFFFF, 0, 0, 0},
				Unknown0A: [PartySlotCount]uint16{0xABCD, 0, 0, 0},
			},
			mode: CopyMode{UnknownFlag628D: true, UnknownFlagE9EF: true},
			want: [PartySlotCount]PartySlotRecord{
				{Unknown08: 0xFFFF, Unknown0A: 0xABCD},
				{},
				{},
				{},
			},
			// $ABCD has bit 13 set, so slot 0 is excluded.
			wantSlotMask: 0x0E,
		},
		{
			name: "distinct values catch slot or field transposition",
			src: CharacterFieldsSource{
				CurrentHP:      [PartySlotCount]uint16{0x1101, 0x1102, 0x1103, 0x1104},
				CandidateMaxHP: [PartySlotCount]uint16{0x2201, 0x2202, 0x2203, 0x2204},
				Unknown04:      [PartySlotCount]uint16{0x3301, 0x3302, 0x3303, 0x3304},
				Unknown06:      [PartySlotCount]uint16{0x4401, 0x4402, 0x4403, 0x4404},
				Unknown08:      [PartySlotCount]uint16{0x5501, 0x5502, 0x5503, 0x5504},
				Unknown0A:      [PartySlotCount]uint16{0x0601, 0x0602, 0x0603, 0x0604},
			},
			mode: CopyMode{},
			want: [PartySlotCount]PartySlotRecord{
				{CurrentHP: 0x1101, CandidateMaxHP: 0x2201, Unknown04: 0x3301, Unknown06: 0x4401, Unknown08: 0x5501, Unknown0A: 0x0601},
				{CurrentHP: 0x1102, CandidateMaxHP: 0x2202, Unknown04: 0x3302, Unknown06: 0x4402, Unknown08: 0x5502, Unknown0A: 0x0602},
				{CurrentHP: 0x1103, CandidateMaxHP: 0x2203, Unknown04: 0x3303, Unknown06: 0x4403, Unknown08: 0x5503, Unknown0A: 0x0603},
				{CurrentHP: 0x1104, CandidateMaxHP: 0x2204, Unknown04: 0x3304, Unknown06: 0x4404, Unknown08: 0x5504, Unknown0A: 0x0604},
			},
			wantSlotMask: 0x0F,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst [PartySlotCount]PartySlotRecord
			// Pre-populate the copied fields so all-zero expectations
			// prove fields are overwritten, not merely left at zero.
			for slot := range dst {
				dst[slot].CurrentHP = 0xDEAD
				dst[slot].CandidateMaxHP = 0xDEAD
				dst[slot].Unknown04 = 0xDEAD
				dst[slot].Unknown06 = 0xDEAD
				dst[slot].Unknown08 = 0xDEAD
				dst[slot].Unknown0A = 0xDEAD
			}

			gotMask := CopyCharacterFields(&tt.src, &dst, tt.mode)

			if dst != tt.want {
				t.Errorf("CopyCharacterFields() dst = %+v, want %+v", dst, tt.want)
			}
			if gotMask != tt.wantSlotMask {
				t.Errorf("CopyCharacterFields() slot mask = %#02x, want %#02x", gotMask, tt.wantSlotMask)
			}
		})
	}
}

// The routine writes only offsets +0x00..+0x0B of each record; bytes
// +0x0C..+0x1F were observed holding live data written by other code and
// must be left untouched.
func TestCopyCharacterFieldsPreservesUnknownBytes(t *testing.T) {
	src := CharacterFieldsSource{
		CurrentHP: [PartySlotCount]uint16{1, 2, 3, 4},
	}

	var dst [PartySlotCount]PartySlotRecord
	for slot := range dst {
		for i := range dst[slot].Unknown0C {
			dst[slot].Unknown0C[i] = 0xA5
		}
	}

	CopyCharacterFields(&src, &dst, CopyMode{})

	for slot := range dst {
		for i, b := range dst[slot].Unknown0C {
			if b != 0xA5 {
				t.Fatalf("slot %d Unknown0C[%#x] = %#x, want 0xA5 (byte was modified)", slot, i, b)
			}
		}
	}
}
