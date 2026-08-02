package content

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/game/battledata"
)

// Manifest entries for the archived battle tables.
const (
	MonstersAssetID       = "ASSET-RAW-0001"
	FormationsAssetID     = "ASSET-RAW-0002"
	FormationFlagsAssetID = "ASSET-RAW-0003"
)

// BattleTables are the decoded formation and monster tables.
type BattleTables struct {
	Formations *battledata.Table
	Monsters   *battledata.Table
	// FormationFlags is the parallel 4-byte-stride table at
	// ROMCPU:$CF5900. Its location is Confirmed and its fields are Unknown,
	// so it is exposed as records rather than accessors.
	FormationFlags *battledata.Table
}

// LoadBattleTables reads the archived tables, hash-verified.
//
// They are archived as raw bytes rather than decoded JSON because most of
// each record is still undecoded: writing out only the mapped fields would
// silently discard the bytes a future unit needs.
func LoadBattleTables(a *Archive) (*BattleTables, error) {
	load := func(id string, stride int) (*battledata.Table, error) {
		raw, err := a.Read(id)
		if err != nil {
			return nil, err
		}
		t, err := battledata.NewTable(raw, stride)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		return t, nil
	}

	formations, err := load(FormationsAssetID, battledata.FormationSize)
	if err != nil {
		return nil, err
	}
	monsters, err := load(MonstersAssetID, battledata.MonsterSize)
	if err != nil {
		return nil, err
	}
	flags, err := load(FormationFlagsAssetID, 4)
	if err != nil {
		return nil, err
	}
	return &BattleTables{Formations: formations, Monsters: monsters, FormationFlags: flags}, nil
}

// Formation decodes one formation record.
func (b *BattleTables) Formation(id int) (battledata.Formation, error) {
	raw, err := b.Formations.Raw(id)
	if err != nil {
		return battledata.Formation{}, fmt.Errorf("formation %d: %w", id, err)
	}
	return battledata.DecodeFormation(raw)
}

// Monster decodes one monster record.
func (b *BattleTables) Monster(id int) (battledata.Monster, error) {
	raw, err := b.Monsters.Raw(id)
	if err != nil {
		return battledata.Monster{}, fmt.Errorf("monster %d: %w", id, err)
	}
	return battledata.DecodeMonster(raw)
}
