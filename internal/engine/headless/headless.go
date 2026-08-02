// Package headless runs an engine.Machine with no display.
//
// It is the **authoritative** host for tests. A windowed host is at the mercy
// of the platform — it may skip Draw, coalesce or pause Update on focus loss,
// and run at whatever rate the monitor dictates — so it can never be the
// reference for a frame-exact comparison. This host calls Tick exactly once
// per frame and Render exactly when asked, which is what makes a recorded
// frame hash mean something.
//
// It uses only the standard library, so the demo is fully testable in CI with
// no display, no GPU, and no third-party dependency.
package headless

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/bensabler/ff6-decompile/internal/engine"
)

// Options configures a headless run.
type Options struct {
	// Frames is the number of ticks to simulate.
	Frames int
	// CaptureEvery writes a PNG every n frames when OutDir is set. Zero
	// disables periodic capture.
	CaptureEvery int
	// CaptureLast writes a PNG of the final frame when OutDir is set.
	CaptureLast bool
	// OutDir receives captured PNGs. Empty means capture nothing.
	OutDir string
	// HashEvery records a frame hash every n frames. Zero means every frame.
	HashEvery int
	// Progress receives one line per capture, for operator feedback.
	Progress io.Writer
}

// Result is what a run produced.
type Result struct {
	// Frames simulated.
	Frames int
	// Hashes of the framebuffer, in capture order.
	Hashes [][32]byte
	// Captured lists the PNG paths written.
	Captured []string
	// StoppedEarly is true if every scene was popped before Frames elapsed.
	StoppedEarly bool
}

// Run drives the machine and returns the frame hashes.
func Run(m *engine.Machine, o Options) (*Result, error) {
	if o.Frames < 0 {
		return nil, fmt.Errorf("headless run: negative frame count %d", o.Frames)
	}
	if o.HashEvery <= 0 {
		o.HashEvery = 1
	}
	if o.OutDir != "" {
		if err := os.MkdirAll(o.OutDir, 0o755); err != nil {
			return nil, fmt.Errorf("headless run: creating %s: %w", o.OutDir, err)
		}
	}

	res := &Result{}
	for i := 0; i < o.Frames; i++ {
		m.Tick()
		res.Frames++

		fb, pal := m.Render()
		if res.Frames%o.HashEvery == 0 {
			res.Hashes = append(res.Hashes, fb.Sum256())
		}

		last := i == o.Frames-1
		shouldCapture := o.OutDir != "" &&
			((o.CaptureEvery > 0 && res.Frames%o.CaptureEvery == 0) || (o.CaptureLast && last))
		if shouldCapture {
			p := filepath.Join(o.OutDir, fmt.Sprintf("frame-%06d.png", res.Frames))
			if err := writePNG(p, fb.Paletted(pal)); err != nil {
				return res, err
			}
			res.Captured = append(res.Captured, p)
			if o.Progress != nil {
				fmt.Fprintf(o.Progress, "captured %s\n", p)
			}
		}

		if m.Done() {
			res.StoppedEarly = true
			return res, nil
		}
	}
	return res, nil
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	defer f.Close()
	if err := (&png.Encoder{CompressionLevel: png.DefaultCompression}).Encode(f, img); err != nil {
		return fmt.Errorf("encoding %s: %w", path, err)
	}
	return f.Close()
}
