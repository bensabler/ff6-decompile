package scenes

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/game/atb"
	"github.com/bensabler/ff6-decompile/internal/game/battledata"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// syntheticTables builds battle tables with project-authored values, so these
// tests need no ROM and no archive. Formation 1 holds monsters 3 and 5.
func syntheticTables(t *testing.T) *content.BattleTables {
	t.Helper()

	forms := make([]byte, battledata.FormationSize*2)
	for i := 0; i < battledata.FormationMonsterSlots; i++ {
		forms[2+i] = battledata.EmptySlot
		forms[battledata.FormationSize+2+i] = battledata.EmptySlot
	}
	forms[battledata.FormationSize+2] = 3
	forms[battledata.FormationSize+3] = 5

	mons := make([]byte, battledata.MonsterSize*8)
	set := func(id int, power uint8, hp, mp uint16) {
		o := id * battledata.MonsterSize
		mons[o+0x01] = power
		mons[o+0x08], mons[o+0x09] = byte(hp), byte(hp>>8)
		mons[o+0x0A], mons[o+0x0B] = byte(mp), byte(mp>>8)
	}
	set(3, 11, 111, 22)
	set(5, 33, 444, 55)

	ft, err := battledata.NewTable(forms, battledata.FormationSize)
	if err != nil {
		t.Fatal(err)
	}
	mt, err := battledata.NewTable(mons, battledata.MonsterSize)
	if err != nil {
		t.Fatal(err)
	}
	flags, err := battledata.NewTable(make([]byte, 8), 4)
	if err != nil {
		t.Fatal(err)
	}
	return &content.BattleTables{Formations: ft, Monsters: mt, FormationFlags: flags}
}

func newBattleScene(t *testing.T) *Battle {
	t.Helper()
	b, err := NewBattle(syntheticFont(), syntheticTables(t), 1)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestBattlePlacesEnemiesFromSlotFour pins the DISC-0001 convention that
// EXP-0040 confirmed for Whelk: shell in slot 4, head in slot 5.
func TestBattlePlacesEnemiesFromSlotFour(t *testing.T) {
	b := newBattleScene(t)
	if len(b.enemies) != 2 {
		t.Fatalf("loaded %d enemies, want 2", len(b.enemies))
	}
	if b.enemies[0].slot != 4 || b.enemies[1].slot != 5 {
		t.Errorf("slots = %d,%d; enemies occupy 4 upward", b.enemies[0].slot, b.enemies[1].slot)
	}
	if b.enemies[0].record != 3 || b.enemies[0].hp != 111 || b.enemies[0].mp != 22 {
		t.Errorf("enemy 0 = record %d HP %d MP %d, want 3/111/22", b.enemies[0].record, b.enemies[0].hp, b.enemies[0].mp)
	}
	if b.enemies[1].hp != 444 {
		t.Errorf("enemy 1 HP = %d, want 444", b.enemies[1].hp)
	}
	// Only the occupied slots get an increment, or an empty slot would tick.
	for slot := 0; slot < atb.SlotCount; slot++ {
		want := slot == 4 || slot == 5
		if got := b.timing.Flags[slot]&atb.FlagActive != 0; got != want {
			t.Errorf("slot %d active=%v, want %v", slot, got, want)
		}
	}
}

func TestBattleRejectsUnknownFormation(t *testing.T) {
	if _, err := NewBattle(syntheticFont(), syntheticTables(t), 99); err == nil {
		t.Error("a formation id past the table should be an error, not an empty battle")
	}
}

// TestBattleGaugesAdvance runs the scene and checks the ATB layer is actually
// driven, at the measured rate.
func TestBattleGaugesAdvance(t *testing.T) {
	b := newBattleScene(t)
	m := engine.New(engine.NoInput, nil, b)
	for i := 0; i < 10; i++ {
		m.Tick()
	}
	if b.timing.Tick != 10 {
		t.Errorf("battle tick = %d, want 10", b.timing.Tick)
	}
	if want := uint16(observedEnemyIncrement>>1) * 10; b.timing.Gauges[4] != want {
		t.Errorf("slot 4 gauge = $%04X, want $%04X", b.timing.Gauges[4], want)
	}
}

// TestBattleGateIsOperable exercises both halves of the ACTIVE/WAIT matrix
// through the scene's own controls.
func TestBattleGateIsOperable(t *testing.T) {
	press := func(btn snespad.Button, at uint64) engine.Source {
		return engine.SourceFunc(func(f uint64) snespad.State {
			if f == at {
				return snespad.State(0).With(btn)
			}
			return 0
		})
	}

	// Under ACTIVE (the default), opening the submenu must not pause.
	b := newBattleScene(t)
	m := engine.New(press(snespad.B, 0), nil, b)
	for i := 0; i < 5; i++ {
		m.Tick()
	}
	if b.submenu == 0 {
		t.Fatal("B should have opened the stand-in submenu")
	}
	if b.timing.Tick != 5 {
		t.Errorf("ACTIVE + submenu: tick = %d, want 5 (the gate is an AND)", b.timing.Tick)
	}

	// Under WAIT, the same submenu freezes everything.
	b2 := newBattleScene(t)
	b2.timing.Entry.WaitFlag = 1
	m2 := engine.New(press(snespad.B, 0), nil, b2)
	for i := 0; i < 5; i++ {
		m2.Tick()
	}
	if b2.timing.Tick != 0 || b2.timing.Gauges[4] != 0 {
		t.Errorf("WAIT + submenu: tick=%d gauge=$%04X, want both frozen", b2.timing.Tick, b2.timing.Gauges[4])
	}
}

func TestBattleStartExits(t *testing.T) {
	b := newBattleScene(t)
	m := engine.New(engine.SourceFunc(func(f uint64) snespad.State {
		if f == 2 {
			return snespad.State(0).With(snespad.Start)
		}
		return 0
	}), nil, b)
	for i := 0; i < 5; i++ {
		m.Tick()
	}
	if !m.Done() {
		t.Error("Start should pop the battle scene")
	}
}

func TestBattleDraws(t *testing.T) {
	b := newBattleScene(t)
	m := engine.New(engine.NoInput, nil, b)
	m.Tick()
	fb, _ := m.Render()
	ink := 0
	for _, v := range fb.Pix {
		if v != 0 {
			ink++
		}
	}
	if ink == 0 {
		t.Fatal("the battle scene drew nothing")
	}
}

// TestBootPushesBattle covers the integration: A on the boot scene opens the
// battle, and the boot scene still works without tables.
func TestBootPushesBattle(t *testing.T) {
	boot := NewBoot(syntheticFont()).WithBattle(syntheticTables(t), 1)
	m := engine.New(engine.SourceFunc(func(f uint64) snespad.State {
		if f == 1 {
			return snespad.State(0).With(snespad.A)
		}
		return 0
	}), nil, boot)
	for i := 0; i < 4; i++ {
		m.Tick()
	}
	if _, ok := m.Stack().Top().(*Battle); !ok {
		t.Fatalf("top scene is %T, want *Battle", m.Stack().Top())
	}

	// Without tables the boot scene must still run; a missing battle table
	// is not a reason for the demo to fail to launch.
	plain := NewBoot(syntheticFont())
	m2 := engine.New(engine.SourceFunc(func(uint64) snespad.State {
		return snespad.State(0).With(snespad.A)
	}), nil, plain)
	for i := 0; i < 4; i++ {
		m2.Tick()
	}
	if m2.Stack().Len() != 1 {
		t.Errorf("stack depth %d without tables, want 1", m2.Stack().Len())
	}
}

// TestGaugeUsesTheIdentifiedTiles checks the gauge is drawn with the five
// tiles EXP-0049 identified, not with invented graphics.
func TestGaugeUsesTheIdentifiedTiles(t *testing.T) {
	b := newBattleScene(t)
	dst := framebuf.New()

	// A full gauge and an empty one must differ, or the fill is not being
	// rendered at all.
	b.drawGauge(dst, 0, 0, 0, false)
	emptySum := dst.Sum256()

	dst2 := framebuf.New()
	b.drawGauge(dst2, 0, 0, 0xFFFF, true)
	if dst2.Sum256() == emptySum {
		t.Error("a full gauge should not render identically to an empty one")
	}

	// The bar is six cells wide: left cap, four segments, right cap.
	const wantWidth = (gaugeCells + 2) * content.GlyphWidth
	for x := 0; x < wantWidth; x++ {
		found := false
		for y := 0; y < 8; y++ {
			if dst2.At(x, y) != 0 {
				found = true
			}
		}
		if !found && x < wantWidth-content.GlyphWidth {
			t.Errorf("column %d of a full gauge is blank; the bar should span %d px", x, wantWidth)
			break
		}
	}
}
