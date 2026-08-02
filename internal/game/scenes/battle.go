package scenes

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/game/atb"
	"github.com/bensabler/ff6-decompile/internal/game/battledata"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// Battle renders the enemy side of the battle HUD from extracted data and
// runs the Confirmed ATB model over it.
//
// # What is real here
//
// The formation composition, every enemy's HP/MP/power, the gauge tiles, the
// font, and the gauge arithmetic all come from verified sources. Nothing on
// screen is a mock-up of a battle.
//
// # What is not, and is therefore labelled
//
// Three things have no extractable source, and this scene shows that rather
// than inventing them. Each has a row in DEMO-0001-DEVIATIONS.md:
//
//   - **Monster names** (D2). The name table is unlocated, so enemies are
//     labelled by record id. A plausible-looking invented name would be
//     indistinguishable from recovered data, which is the failure mode the
//     deviations register exists to prevent.
//   - **ATB increments** (D3). ROMCPU:$C209F0 is only partly decoded, so the
//     increment is the value EXP-0043 *measured* for enemy slots rather than
//     one this project can compute.
//   - **The party side** (D4). Character initialisation is unlocated; there is
//     no table to read party HP, MP, or names from. The rows are drawn empty
//     and labelled.
type Battle struct {
	font    *content.Font
	tables  *content.BattleTables
	id      int
	enemies []enemy
	timing  atb.State
	frame   uint64
	submenu uint8
}

type enemy struct {
	slot   int
	record int
	hp     uint16
	mp     uint16
	power  uint8
}

// observedEnemyIncrement is the per-slot ATB increment EXP-0043 measured for
// the three enemies of formation 14 (all $9C, advancing gauges by $4E per
// frame). It stands in for the undecoded formula at ROMCPU:$C209F0.
//
// DEMO-0001 deviation D3. Replace it when the formula is decoded; the ATB
// package already takes increments as an input for exactly this reason.
const observedEnemyIncrement uint16 = 0x9C

// NewBattle builds the scene for a formation id.
//
// Enemies occupy slots 4 upward, the DISC-0001 convention that EXP-0040
// confirmed for Whelk (shell slot 4, head slot 5).
func NewBattle(font *content.Font, tables *content.BattleTables, formationID int) (*Battle, error) {
	f, err := tables.Formation(formationID)
	if err != nil {
		return nil, err
	}
	b := &Battle{font: font, tables: tables, id: formationID}
	b.timing.Entry = atb.SampleEntry(atb.Config{Mode: atb.Active, Speed: 2})

	for i, id := range f.OccupiedIDs() {
		slot := 4 + i
		if slot >= atb.SlotCount {
			break
		}
		m, err := tables.Monster(int(id))
		if err != nil {
			return nil, fmt.Errorf("formation %d slot %d: %w", formationID, slot, err)
		}
		b.enemies = append(b.enemies, enemy{
			slot: slot, record: int(id), hp: m.HP(), mp: m.MP(), power: m.Power(),
		})
		b.timing.Increments[slot] = observedEnemyIncrement
		b.timing.Flags[slot] = atb.FlagActive
	}
	return b, nil
}

// FormationID returns the formation being shown.
func (b *Battle) FormationID() int { return b.id }

// Timing exposes the ATB state, for tests.
func (b *Battle) Timing() *atb.State { return &b.timing }

func (b *Battle) Update(ctx *engine.Context) {
	b.frame = ctx.Frame

	// B toggles a stand-in submenu, so the ACTIVE/WAIT gate is operable and
	// visible. Under Active it changes nothing — which is the point: the
	// gate is an AND, and EXP-0044 found the pause narrower than assumed.
	if ctx.Input.JustPressed(snespad.B) {
		b.submenu ^= 1
	}
	// Y flips the mode, so both halves of the matrix can be seen.
	if ctx.Input.JustPressed(snespad.Y) {
		if b.timing.Entry.WaitFlag == 0 {
			b.timing.Entry.WaitFlag = 1
		} else {
			b.timing.Entry.WaitFlag = 0
		}
	}
	if ctx.Input.JustPressed(snespad.Start) {
		ctx.Stack.Pop()
	}

	// Errors here would mean a slot index out of range, which the
	// constructor already prevents.
	_, _ = b.timing.Frame(b.submenu)
}

func (b *Battle) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	const (
		white = 3
		gray  = 2
	)
	dst.Fill(0)

	text := func(x, y int, s string, pal uint8) {
		b.font.DrawString(dst, x, y, s, content.TextOptions{PaletteBase: pal})
	}

	text(8, 6, fmt.Sprintf("FORMATION %d", b.id), white)
	if !battledata.FormationVerified(b.id) {
		// The archive holds more formations than any experiment has
		// verified. Saying so on screen keeps the distinction visible.
		text(8, 18, "UNVERIFIED RECORD", gray)
	} else {
		text(8, 18, "VERIFIED RECORD", gray)
	}

	y := 40
	for _, e := range b.enemies {
		// No name table exists, so the label is the record id (D2).
		text(8, y, fmt.Sprintf("MONSTER %d", e.record), white)
		text(112, y, fmt.Sprintf("HP %d", e.hp), white)
		b.drawGauge(dst, 184, y, b.timing.Gauges[e.slot], b.timing.Ready(e.slot))
		y += 16
	}

	y += 8
	text(8, y, "PARTY SIDE NOT IMPLEMENTED", gray)
	y += 12
	text(8, y, "NO CHARACTER DATA SOURCE", gray)

	mode := "ACTIVE"
	if b.timing.Entry.WaitFlag != 0 {
		mode = "WAIT"
	}
	sub := "CLOSED"
	if b.submenu != 0 {
		sub = "OPEN"
	}
	text(8, 172, fmt.Sprintf("MODE %s  MENU %s", mode, sub), gray)
	text(8, 184, fmt.Sprintf("TICK %d", b.timing.Tick), gray)
	if atb.Paused(b.submenu, b.timing.Entry.WaitFlag) {
		text(8, 196, "GATED - GAUGES FROZEN", white)
	} else {
		text(8, 196, "B MENU  Y MODE  START EXIT", gray)
	}
}

// gaugeCells is the six-cell run EXP-0049 decoded from the HUD tilemap:
// a left cap, four segments, a right cap.
const gaugeCells = 4

// drawGauge renders the HP/ATB bar with the real tiles (CEN-GFX-0005).
//
// The **fill quantisation is Unknown** — EXP-0049 saw only a full bar and a
// nearly empty one, so "one partial cell at the boundary" is a reading of two
// samples, not a decoded rule. It is recorded as such rather than presented as
// the game's behaviour.
func (b *Battle) drawGauge(dst *framebuf.Indexed, x, y int, gauge uint16, ready bool) {
	const (
		leftCap  = 0xF9
		rightCap = 0xFA
		empty    = 0xF0
		partial  = 0xF1
		full     = 0xF8
	)
	cells := make([]byte, 0, gaugeCells+2)
	cells = append(cells, leftCap)

	filled := int(uint32(gauge) * gaugeCells / 0x10000)
	if ready {
		filled = gaugeCells
	}
	for i := 0; i < gaugeCells; i++ {
		switch {
		case i < filled:
			cells = append(cells, full)
		case i == filled && gauge != 0:
			cells = append(cells, partial)
		default:
			cells = append(cells, empty)
		}
	}
	cells = append(cells, rightCap)
	b.font.DrawBytes(dst, x, y, cells, content.TextOptions{})
}
