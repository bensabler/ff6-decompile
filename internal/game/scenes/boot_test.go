package scenes

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/engine/headless"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// syntheticFont builds a font from project-authored pixels: glyph b is a
// hollow box whose thickness varies with b, so different text produces
// different frames.
//
// It is deliberately NOT the FF6 font. Frame hashes derived from it are safe
// to commit; hashes of real FF6 glyphs would be a hash of ROM-derived
// graphics, which the asset policy keeps out of the repository.
func syntheticFont() *content.Font {
	var tiles [256][8][8]uint8
	for i := range tiles {
		if i < 0x20 {
			continue // leave a blank range, as the real block has
		}
		v := uint8(i%3) + 1
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if y == 0 || y == 7 || x == 0 || x == 7 || (i%5 == 0 && x == y) {
					tiles[i][y][x] = v
				}
			}
		}
	}
	return content.NewFont(&tiles)
}

func newBootMachine(src engine.Source) *engine.Machine {
	return engine.New(src, nil, NewBoot(syntheticFont()))
}

func TestBootDrawsSomething(t *testing.T) {
	m := newBootMachine(engine.NoInput)
	m.Tick()
	fb, _ := m.Render()

	ink := 0
	for _, v := range fb.Pix {
		if v != 0 {
			ink++
		}
	}
	if ink == 0 {
		t.Fatal("the boot scene drew nothing")
	}
	// The border proves the full 256x224 extent is being composed, which is
	// what makes a host's letterboxing visible.
	for _, p := range [][2]int{{0, 0}, {framebuf.Width - 1, 0}, {0, framebuf.Height - 1}, {framebuf.Width - 1, framebuf.Height - 1}} {
		if fb.At(p[0], p[1]) == 0 {
			t.Errorf("border pixel (%d,%d) is blank", p[0], p[1])
		}
	}
}

func TestBootRespondsToInput(t *testing.T) {
	// Hold Right for 10 frames, then read the cursor's position.
	m := newBootMachine(engine.SourceFunc(func(f uint64) snespad.State {
		if f < 10 {
			return snespad.State(0).With(snespad.Right)
		}
		return 0
	}))
	b := m.Stack().Top().(*Boot)
	start := b.cursorX
	for i := 0; i < 10; i++ {
		m.Tick()
	}
	if b.cursorX != start+10 {
		t.Errorf("cursorX = %d after 10 frames of Right, want %d", b.cursorX, start+10)
	}
}

func TestBootCursorStaysOnScreen(t *testing.T) {
	m := newBootMachine(engine.SourceFunc(func(uint64) snespad.State {
		return snespad.State(0).With(snespad.Left).With(snespad.Up)
	}))
	b := m.Stack().Top().(*Boot)
	for i := 0; i < 500; i++ {
		m.Tick()
	}
	if b.cursorX != 0 || b.cursorY != 0 {
		t.Errorf("cursor = (%d,%d), want it clamped to (0,0)", b.cursorX, b.cursorY)
	}
	// And the same at the far corner.
	m2 := newBootMachine(engine.SourceFunc(func(uint64) snespad.State {
		return snespad.State(0).With(snespad.Right).With(snespad.Down)
	}))
	b2 := m2.Stack().Top().(*Boot)
	for i := 0; i < 500; i++ {
		m2.Tick()
	}
	if b2.cursorX != framebuf.Width-8 || b2.cursorY != framebuf.Height-8 {
		t.Errorf("cursor = (%d,%d), want it clamped to the far corner", b2.cursorX, b2.cursorY)
	}
}

func TestBootStartExits(t *testing.T) {
	m := newBootMachine(engine.SourceFunc(func(f uint64) snespad.State {
		if f == 5 {
			return snespad.State(0).With(snespad.Start)
		}
		return 0
	}))
	res, err := headless.Run(m, headless.Options{Frames: 100})
	if err != nil {
		t.Fatal(err)
	}
	if !res.StoppedEarly || res.Frames != 6 {
		t.Errorf("Start on frame 5 should end the run at frame 6; got %d frames, early=%v", res.Frames, res.StoppedEarly)
	}
}

// goldenPath holds frame hashes produced from the synthetic font. They are
// project-authored data, not ROM-derived.
const goldenPath = "testdata/boot-frames.json"

type goldens struct {
	Note   string   `json:"note"`
	Script string   `json:"input_script"`
	Frames int      `json:"frames"`
	Hashes []string `json:"hashes"`
}

// goldenScript exercises movement, a button press, and idle frames.
func goldenScript(f uint64) snespad.State {
	switch {
	case f >= 5 && f < 25:
		return snespad.State(0).With(snespad.Right)
	case f >= 30 && f < 40:
		return snespad.State(0).With(snespad.Down)
	case f == 50:
		return snespad.State(0).With(snespad.A)
	default:
		return 0
	}
}

// TestBootFrameGoldens is the demo's first regression test on rendered
// output. Regenerate with -update after an intended visual change, and read
// the diff before committing it.
func TestBootFrameGoldens(t *testing.T) {
	const frames = 60
	res, err := headless.Run(newBootMachine(engine.SourceFunc(goldenScript)), headless.Options{Frames: frames})
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(res.Hashes))
	for i, h := range res.Hashes {
		got[i] = hex.EncodeToString(h[:])
	}

	if os.Getenv("UPDATE_GOLDENS") != "" {
		g := goldens{
			Note: "Frame hashes for the DEMO-0001A boot scene, rendered with a " +
				"project-authored synthetic font (see syntheticFont). Not derived " +
				"from ROM graphics. Regenerate with UPDATE_GOLDENS=1 go test ./internal/game/scenes/",
			Script: "Right@5-24, Down@30-39, A@50",
			Frames: frames,
			Hashes: got,
		}
		raw, err := json.MarshalIndent(g, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, append(raw, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %d hashes to %s", len(got), goldenPath)
		return
	}

	raw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading goldens: %v (regenerate with UPDATE_GOLDENS=1)", err)
	}
	var want goldens
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	if len(want.Hashes) != len(got) {
		t.Fatalf("golden has %d hashes, run produced %d", len(want.Hashes), len(got))
	}
	for i := range got {
		if got[i] != want.Hashes[i] {
			t.Fatalf("frame %d differs\n  got  %s\n  want %s\nIf this change is intended, regenerate with UPDATE_GOLDENS=1", i+1, got[i], want.Hashes[i])
		}
	}
}
