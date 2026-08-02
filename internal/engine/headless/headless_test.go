package headless

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// tick draws a moving block, so successive frames differ.
type tick struct {
	n    int
	quit uint64
}

func (s *tick) Update(ctx *engine.Context) {
	s.n++
	if s.quit != 0 && ctx.Frame+1 >= s.quit {
		ctx.Stack.Pop()
	}
}

func (s *tick) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	dst.Rect(s.n%framebuf.Width, 0, 4, 4, uint8(s.n%254)+1)
}

func TestRunSimulatesExactlyNFrames(t *testing.T) {
	sc := &tick{}
	m := engine.New(engine.NoInput, nil, sc)
	res, err := Run(m, Options{Frames: 30})
	if err != nil {
		t.Fatal(err)
	}
	if res.Frames != 30 || m.Frame() != 30 {
		t.Errorf("frames = %d (machine %d), want 30", res.Frames, m.Frame())
	}
	if len(res.Hashes) != 30 {
		t.Errorf("hashes = %d, want 30 (one per frame by default)", len(res.Hashes))
	}
	if res.StoppedEarly {
		t.Error("StoppedEarly should be false when no scene popped")
	}
}

func TestRunStopsWhenTheStackEmpties(t *testing.T) {
	m := engine.New(engine.NoInput, nil, &tick{quit: 10})
	res, err := Run(m, Options{Frames: 1000})
	if err != nil {
		t.Fatal(err)
	}
	if !res.StoppedEarly {
		t.Error("StoppedEarly should be true after the last scene pops")
	}
	if res.Frames != 10 {
		t.Errorf("frames = %d, want 10", res.Frames)
	}
}

func TestRunIsReproducible(t *testing.T) {
	build := func() *engine.Machine {
		return engine.New(engine.SourceFunc(func(f uint64) snespad.State {
			if f%5 == 0 {
				return snespad.State(0).With(snespad.A)
			}
			return 0
		}), nil, &tick{})
	}
	a, err := Run(build(), Options{Frames: 50})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Run(build(), Options{Frames: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Hashes) != len(b.Hashes) {
		t.Fatalf("hash counts differ: %d vs %d", len(a.Hashes), len(b.Hashes))
	}
	for i := range a.Hashes {
		if a.Hashes[i] != b.Hashes[i] {
			t.Fatalf("frame %d differs between identical runs", i+1)
		}
	}
}

func TestHashEvery(t *testing.T) {
	m := engine.New(engine.NoInput, nil, &tick{})
	res, err := Run(m, Options{Frames: 30, HashEvery: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Hashes) != 3 {
		t.Errorf("hashes = %d, want 3", len(res.Hashes))
	}
}

func TestCapture(t *testing.T) {
	dir := t.TempDir()
	var progress bytes.Buffer
	m := engine.New(engine.NoInput, nil, &tick{})
	res, err := Run(m, Options{Frames: 20, CaptureEvery: 10, CaptureLast: true, OutDir: dir, Progress: &progress})
	if err != nil {
		t.Fatal(err)
	}
	// Frames 10 and 20; 20 is both a multiple and the last, and must not be
	// written twice.
	if len(res.Captured) != 2 {
		t.Fatalf("captured %v, want 2 files", res.Captured)
	}
	for _, p := range res.Captured {
		st, err := os.Stat(p)
		if err != nil {
			t.Errorf("stat %s: %v", p, err)
			continue
		}
		if st.Size() == 0 {
			t.Errorf("%s is empty", p)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "frame-000020.png")); err != nil {
		t.Errorf("expected the final frame to be captured: %v", err)
	}
	if progress.Len() == 0 {
		t.Error("Progress should have received a line per capture")
	}
}

func TestCaptureNeedsAnOutDir(t *testing.T) {
	m := engine.New(engine.NoInput, nil, &tick{})
	res, err := Run(m, Options{Frames: 5, CaptureEvery: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Captured) != 0 {
		t.Error("capture without OutDir should write nothing")
	}
}

func TestRunRejectsNegativeFrames(t *testing.T) {
	m := engine.New(engine.NoInput, nil, &tick{})
	if _, err := Run(m, Options{Frames: -1}); err == nil {
		t.Error("a negative frame count should be an error")
	}
}

func TestZeroFramesIsValid(t *testing.T) {
	m := engine.New(engine.NoInput, nil, &tick{})
	res, err := Run(m, Options{Frames: 0})
	if err != nil {
		t.Fatal(err)
	}
	if res.Frames != 0 || len(res.Hashes) != 0 {
		t.Errorf("zero frames produced %+v", res)
	}
}
