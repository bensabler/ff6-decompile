// Package romorigin traces a captured memory image back to the ROM spans it
// was copied from.
//
// # Why this exists
//
// A preserved savestate holds what the PPU was drawing. The question that
// matters for reconstruction is where those bytes came from, and the cheapest
// discriminating answer is whether they appear in the ROM verbatim.
//
// If a region maps to a contiguous ROM span, it was copied — the format is
// "none", the extractor is a slice, and no decompressor is needed. If it does
// not, something transformed it: compression, runtime composition from
// smaller records, or generation. Knowing which of those two a region is in
// is the difference between a one-day extractor and an open research program,
// and until now the project had no way to ask.
//
// # What a match does and does not prove
//
// A match proves the bytes are identical, not that this span is the source.
// Short probes can collide, and identical data can appear more than once —
// the HUD font's blank tiles, for instance, are a run of zeros that matches
// everywhere. MinRun exists to keep collisions out of the report, and Trace
// skips probe windows with too little variety to be distinctive.
//
// A *non*-match proves less than it seems. It rules out verbatim copying, and
// nothing else. Bit-plane reordering, a different tile arrangement, or a
// single changed byte all defeat it. "Not verbatim" is a lead toward
// compression, not evidence of it.
package romorigin

import "bytes"

// Block is one region of the captured image and the ROM span it matches.
type Block struct {
	// ImageOffset is the byte offset within the captured image.
	ImageOffset int
	// ROMOffset is the matching offset in the ROM image (ROMFILE domain).
	ROMOffset int
	// Length is how far the two agree, extended in both directions from the
	// probe that found the match.
	Length int
}

// Options tune the search.
type Options struct {
	// Probe is the window size used to look for a match. Larger is more
	// selective and slower to fall back on when a region does not match.
	Probe int
	// MinRun is the shortest run reported. Runs below it are dropped as
	// likely collisions rather than provenance.
	MinRun int
	// MinDistinct is how many distinct byte values a probe window must hold
	// to be worth searching. A window of one repeated value matches
	// everywhere and tells you nothing.
	MinDistinct int
	// Limit stops the scan at this image offset. Zero means the whole image.
	Limit int
}

// DefaultOptions are tuned for SNES graphics: a 32-byte probe is exactly one
// 4bpp tile, and 64 bytes is two, which is long enough that an accidental
// match is not worth reporting.
func DefaultOptions() Options {
	return Options{Probe: 32, MinRun: 64, MinDistinct: 3}
}

// Trace maps image back onto rom.
//
// It walks the image, and at each position takes a probe window, searches the
// ROM for it, then extends the match backward and forward to its full run.
// The scan resumes past the run, so overlapping reports do not accumulate.
//
// Both arguments are read-only. Trace allocates only the returned slice.
func Trace(image, rom []byte, o Options) []Block {
	if o.Probe <= 0 {
		o.Probe = 32
	}
	if o.MinDistinct <= 0 {
		o.MinDistinct = 1
	}
	end := len(image)
	if o.Limit > 0 && o.Limit < end {
		end = o.Limit
	}

	var out []Block
	for pos := 0; pos+o.Probe <= end; {
		win := image[pos : pos+o.Probe]
		if distinct(win) < o.MinDistinct {
			pos += o.Probe
			continue
		}
		i := bytes.Index(rom, win)
		if i < 0 {
			pos += o.Probe
			continue
		}

		back := 0
		for pos-back > 0 && i-back > 0 && image[pos-back-1] == rom[i-back-1] {
			back++
		}
		fwd := 0
		for pos+fwd < len(image) && i+fwd < len(rom) && image[pos+fwd] == rom[i+fwd] {
			fwd++
		}

		b := Block{ImageOffset: pos - back, ROMOffset: i - back, Length: back + fwd}
		if b.Length >= o.MinRun {
			out = append(out, b)
		}
		// Always resume past the run, matched or not, so a long run that
		// fails MinRun cannot be rediscovered probe by probe.
		next := b.ImageOffset + b.Length
		if next <= pos {
			next = pos + o.Probe
		}
		pos = next
	}
	return out
}

// Coverage returns how many bytes of the image the blocks account for.
//
// Blocks from Trace never overlap, but Coverage tolerates it so a caller may
// pass a filtered or merged set.
func Coverage(blocks []Block) int {
	if len(blocks) == 0 {
		return 0
	}
	total, prevEnd := 0, -1
	for _, b := range blocks {
		start, end := b.ImageOffset, b.ImageOffset+b.Length
		if start < prevEnd {
			start = prevEnd
		}
		if end > start {
			total += end - start
			prevEnd = end
		}
	}
	return total
}

// Contiguous reports whether b continues a, in both the image and the ROM.
// A tileset uploaded as one transfer shows up as several probe hits that this
// predicate joins back into one span.
func Contiguous(a, b Block) bool {
	return a.ImageOffset+a.Length == b.ImageOffset && a.ROMOffset+a.Length == b.ROMOffset
}

// Merge joins runs that are contiguous in both domains.
func Merge(blocks []Block) []Block {
	if len(blocks) == 0 {
		return nil
	}
	out := []Block{blocks[0]}
	for _, b := range blocks[1:] {
		last := &out[len(out)-1]
		if Contiguous(*last, b) {
			last.Length += b.Length
			continue
		}
		out = append(out, b)
	}
	return out
}

func distinct(b []byte) int {
	var seen [256]bool
	n := 0
	for _, v := range b {
		if !seen[v] {
			seen[v] = true
			n++
		}
	}
	return n
}
