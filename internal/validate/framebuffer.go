package validate

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/bgr555"
)

// Two comparisons, and the difference between them is not academic.
//
// framebuf.Indexed.Sum256 hashes palette *indices*. That is the right identity
// for a composition golden — it makes a colour change and a layout change fail
// differently, and it lets goldens be committed without tracking ROM-derived
// colour data.
//
// It also cannot see anything a viewer sees. Deviation D0 is the proof: for
// seven units the demo drew every bright string onto palette entries the
// palette had never defined, 32 % of drawn ink resolved to a visible colour,
// and the frame goldens passed the whole time. A frame can hash correctly and
// be blank on screen.
//
// So this file provides both, named so they cannot be confused, and a caller
// that wants "does it look right" is pushed to the one that can answer.

// IndexedDiff is the result of comparing two index planes.
type IndexedDiff struct {
	// Pixels is how many pixels the two planes disagree on.
	Pixels int
	// Total is how many pixels were compared.
	Total int
	// FirstX and FirstY locate the first disagreement, -1 when equal.
	FirstX, FirstY int
	// Histogram maps got-index -> want-index -> count, for the first
	// mismatching combinations. It is what turns "1,412 pixels differ" into
	// "every 3 became a 7", which is a palette-base error rather than a
	// layout error.
	Histogram map[[2]uint8]int
}

// Equal reports whether the planes agree everywhere.
func (d IndexedDiff) Equal() bool { return d.Pixels == 0 }

func (d IndexedDiff) String() string {
	if d.Equal() {
		return fmt.Sprintf("index planes identical (%d pixels)", d.Total)
	}
	return fmt.Sprintf("index planes differ on %d of %d pixels (%.2f%%), first at (%d,%d)",
		d.Pixels, d.Total, 100*float64(d.Pixels)/float64(d.Total), d.FirstX, d.FirstY)
}

// CompareIndexed compares two framebuffers by palette index.
//
// This answers "is the composition the same". It does **not** answer whether
// the two look the same, because two different indices can hold the same
// colour and one index can be undefined. Use CompareResolved for that.
func CompareIndexed(got, want *framebuf.Indexed) IndexedDiff {
	d := IndexedDiff{FirstX: -1, FirstY: -1, Histogram: map[[2]uint8]int{}}
	if got == nil || want == nil {
		return d
	}
	n := min(len(got.Pix), len(want.Pix))
	d.Total = n
	for i := 0; i < n; i++ {
		g, w := got.Pix[i], want.Pix[i]
		if g == w {
			continue
		}
		if d.Pixels == 0 {
			d.FirstX, d.FirstY = i%framebuf.Width, i/framebuf.Width
		}
		d.Pixels++
		d.Histogram[[2]uint8{g, w}]++
	}
	return d
}

// ResolvedDiff is the result of comparing two frames as a viewer sees them.
type ResolvedDiff struct {
	// Pixels is how many pixels differ in colour.
	Pixels int
	// Total is how many pixels were compared.
	Total int
	// FirstX, FirstY locate the first difference, -1 when equal.
	FirstX, FirstY int

	// Drawn is how many pixels of got are non-transparent.
	Drawn int
	// Visible is how many of those resolve to a colour other than the
	// frame's background colour — the number D0 moved from 32 % to 63 %.
	Visible int
}

// Equal reports whether the resolved images agree everywhere.
func (d ResolvedDiff) Equal() bool { return d.Pixels == 0 }

// VisibleFraction returns the share of drawn ink that a viewer can actually
// see, in [0,1]. It is 0 when nothing was drawn.
func (d ResolvedDiff) VisibleFraction() float64 {
	if d.Drawn == 0 {
		return 0
	}
	return float64(d.Visible) / float64(d.Drawn)
}

func (d ResolvedDiff) String() string {
	vis := fmt.Sprintf("%d of %d drawn pixels visible (%.0f%%)", d.Visible, d.Drawn, 100*d.VisibleFraction())
	if d.Equal() {
		return fmt.Sprintf("resolved images identical (%d pixels); %s", d.Total, vis)
	}
	return fmt.Sprintf("resolved images differ on %d of %d pixels (%.2f%%), first at (%d,%d); %s",
		d.Pixels, d.Total, 100*float64(d.Pixels)/float64(d.Total), d.FirstX, d.FirstY, vis)
}

// CompareResolved compares two frames after resolving each through its own
// palette, which is what a viewer sees.
//
// Two frames with different indices can be identical here — that is the
// point. It is also the only comparison that can fail when a frame renders
// blank, which CompareIndexed and Sum256 structurally cannot.
func CompareResolved(got *framebuf.Indexed, gotPal *framebuf.Palette,
	want *framebuf.Indexed, wantPal *framebuf.Palette) ResolvedDiff {

	d := ResolvedDiff{FirstX: -1, FirstY: -1}
	if got == nil || want == nil || gotPal == nil || wantPal == nil {
		return d
	}
	n := min(len(got.Pix), len(want.Pix))
	d.Total = n
	bg := gotPal[framebuf.Transparent]
	for i := 0; i < n; i++ {
		gi, wi := got.Pix[i], want.Pix[i]
		gc, wc := gotPal[gi], wantPal[wi]
		if gi != framebuf.Transparent {
			d.Drawn++
			if gc != bg {
				d.Visible++
			}
		}
		if gc == wc {
			continue
		}
		if d.Pixels == 0 {
			d.FirstX, d.FirstY = i%framebuf.Width, i/framebuf.Width
		}
		d.Pixels++
	}
	return d
}

// Visibility measures one frame without comparing it to anything, for the
// case where there is no reference yet but "did it render at all" still needs
// an answer.
func Visibility(fb *framebuf.Indexed, pal *framebuf.Palette) ResolvedDiff {
	return CompareResolved(fb, pal, fb, pal)
}

// PaletteUse reports the highest palette index a frame actually draws.
//
// A frame drawing above the entries its palette defines is the D0 defect, and
// it is invisible to any index-based check.
func PaletteUse(fb *framebuf.Indexed) (highest uint8, distinct int) {
	if fb == nil {
		return 0, 0
	}
	var seen [256]bool
	for _, v := range fb.Pix {
		if !seen[v] {
			seen[v] = true
			distinct++
			if v > highest {
				highest = v
			}
		}
	}
	return highest, distinct
}

// DecodeCGRAM converts a captured 512-byte CGRAM image into a Palette, so a
// demo frame can be resolved through the colours the original was using.
//
// This is the bridge that makes a real colour comparison possible at all:
// until readiness F2 locates the palette tables in ROM, the only true FF6
// palette available is the one in a preserved savestate.
func DecodeCGRAM(cgram []byte) (*framebuf.Palette, error) {
	const want = 512
	if len(cgram) != want {
		return nil, fmt.Errorf("validate: CGRAM is %d bytes, want %d", len(cgram), want)
	}
	var p framebuf.Palette
	for i := 0; i < 256; i++ {
		p[i] = bgr555.Decode(uint16(cgram[i*2]) | uint16(cgram[i*2+1])<<8)
	}
	return &p, nil
}
