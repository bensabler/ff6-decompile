package snesaddr

import (
	"errors"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/rom"
)

// TestImageSizeMatchesROM guards the deliberate duplication: internal/platform
// models hardware and must not import the project's ROM loader, so ImageSize
// is a copy. This test is the only thing keeping the copy honest.
func TestImageSizeMatchesROM(t *testing.T) {
	if ImageSize != rom.Size {
		t.Fatalf("ImageSize = %d, rom.Size = %d", ImageSize, rom.Size)
	}
}

func TestROMFile(t *testing.T) {
	tests := []struct {
		name string
		cpu  uint32
		want int
		win  Window
	}{
		// The Confirmed window. These are real tracked addresses.
		{"attack data table", 0xC46AC0, 0x046AC0, WindowHiROMUpper},
		{"HUD font block start", 0xC47FB0, 0x047FB0, WindowHiROMUpper},
		{"HUD font anchor", 0xC46FC0, 0x046FC0, WindowHiROMUpper},
		{"SFX pack", 0xC51EC9, 0x051EC9, WindowHiROMUpper},
		{"monster records", 0xCF0000, 0x0F0000, WindowHiROMUpper},
		{"formation table", 0xCF6200, 0x0F6200, WindowHiROMUpper},
		{"formation flags", 0xCF5900, 0x0F5900, WindowHiROMUpper},
		{"spell names", 0xE6F567, 0x26F567, WindowHiROMUpper},
		{"esper names", 0xE6F6E1, 0x26F6E1, WindowHiROMUpper},
		{"event flag decoder", 0xC0BAED, 0x00BAED, WindowHiROMUpper},
		{"event interpreter tail", 0xC09B5C, 0x009B5C, WindowHiROMUpper},
		{"battle ATB gate", 0xC21124, 0x021124, WindowHiROMUpper},
		{"first byte of the image", 0xC00000, 0x000000, WindowHiROMUpper},
		{"last byte of the image", 0xC00000 + ImageSize - 1, ImageSize - 1, WindowHiROMUpper},

		// Strong-hypothesis windows.
		{"lower window", 0x40BAED, 0x00BAED, WindowHiROMLower},
		// Banks $40-$6F are what a 3 MiB image reaches through the lower
		// window; $70-$7D exist in the layout but address nothing here.
		{"lower window last mapped byte", 0x6FFFFF, 0x2FFFFF, WindowHiROMLower},
		{"mirror window bank $00", 0x00FFEE, 0x00FFEE, WindowHiROMMirror},
		{"mirror window bank $80", 0x80FFEE, 0x00FFEE, WindowHiROMMirror},
		{"mirror window bank $2F", 0x2F8000, 0x2F8000, WindowHiROMMirror},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, win, err := ROMFile(tt.cpu)
			if err != nil {
				t.Fatalf("ROMFile(ROMCPU:$%06X): %v", tt.cpu, err)
			}
			if got != tt.want {
				t.Errorf("ROMFile(ROMCPU:$%06X) = ROMFILE:0x%06X, want ROMFILE:0x%06X", tt.cpu, got, tt.want)
			}
			if win != tt.win {
				t.Errorf("window = %v, want %v", win, tt.win)
			}
		})
	}
}

func TestROMFileRejects(t *testing.T) {
	tests := []struct {
		name string
		cpu  uint32
	}{
		{"WRAM bank $7E", 0x7E2EB5},
		{"WRAM bank $7F", 0x7F0000},
		{"system area in bank $00", 0x002140},
		{"system area in bank $80", 0x804218},
		{"system area at the boundary", 0x007FFF},
		{"beyond 24 bits", 0x1000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, win, err := ROMFile(tt.cpu)
			var e *Error
			if !errors.As(err, &e) {
				t.Fatalf("ROMFile(ROMCPU:$%06X) error = %v, want *snesaddr.Error", tt.cpu, err)
			}
			if win != WindowNone {
				t.Errorf("window = %v, want WindowNone", win)
			}
		})
	}
}

// TestROMFileRejectsPastImageEnd covers the case a bad constant is most likely
// to produce: an address in a valid window that names data the 3 MiB image
// does not contain. Banks $70-$7D and $F0-$FF are real HiROM windows that a
// 3 MiB cartridge simply does not fill.
func TestROMFileRejectsPastImageEnd(t *testing.T) {
	for _, cpu := range []uint32{0xC00000 + ImageSize, 0xFFFFFF, 0x700000, 0x7DFFFF} {
		if off, _, err := ROMFile(cpu); err == nil {
			t.Errorf("ROMFile(ROMCPU:$%06X) = ROMFILE:0x%06X, want an error", cpu, off)
		}
	}
}

func TestConfirmedOnlyCoversTheEvidencedWindow(t *testing.T) {
	if !WindowHiROMUpper.Confirmed() {
		t.Error("banks $C0-$FF carry 18/18 runtime evidence (CORR-0001) and must report Confirmed")
	}
	for _, w := range []Window{WindowNone, WindowHiROMLower, WindowHiROMMirror} {
		if w.Confirmed() {
			t.Errorf("%v has no FF6 runtime evidence and must not report Confirmed", w)
		}
	}
}

func TestROMCPURoundTrip(t *testing.T) {
	for _, off := range []int{0, 1, 0x046AC0, 0x047FB0, 0x0F0000, ImageSize - 1} {
		cpu, err := ROMCPU(off)
		if err != nil {
			t.Fatalf("ROMCPU(ROMFILE:0x%06X): %v", off, err)
		}
		back, win, err := ROMFile(cpu)
		if err != nil {
			t.Fatalf("ROMFile(ROMCPU:$%06X): %v", cpu, err)
		}
		if back != off {
			t.Errorf("round trip ROMFILE:0x%06X -> ROMCPU:$%06X -> ROMFILE:0x%06X", off, cpu, back)
		}
		if !win.Confirmed() {
			t.Errorf("ROMCPU must return an address in the Confirmed window, got %v", win)
		}
	}
}

func TestROMCPURejectsOutOfRange(t *testing.T) {
	for _, off := range []int{-1, ImageSize, ImageSize + 1} {
		if _, err := ROMCPU(off); err == nil {
			t.Errorf("ROMCPU(%d) succeeded, want an error", off)
		}
	}
}

// FuzzRoundTrip asserts the two directions compose for every offset, and that
// neither ever panics.
func FuzzRoundTrip(f *testing.F) {
	for _, seed := range []int{0, 1, 0x046AC0, ImageSize - 1} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, off int) {
		cpu, err := ROMCPU(off)
		if err != nil {
			return // out of range is a valid answer
		}
		back, win, err := ROMFile(cpu)
		if err != nil {
			t.Fatalf("ROMFILE:0x%06X round-tripped to an unmappable ROMCPU:$%06X: %v", off, cpu, err)
		}
		if back != off || !win.Confirmed() {
			t.Fatalf("round trip %d -> $%06X -> %d (window %v)", off, cpu, back, win)
		}
	})
}

// FuzzROMFileNeverPanics feeds arbitrary CPU addresses through the mapping.
func FuzzROMFileNeverPanics(f *testing.F) {
	for _, seed := range []uint32{0, 0xC46AC0, 0x7E2EB5, 0xFFFFFF, 0x1000000} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, cpu uint32) {
		off, win, err := ROMFile(cpu)
		if err != nil {
			return
		}
		if off < 0 || off >= ImageSize {
			t.Fatalf("ROMCPU:$%06X mapped to out-of-image ROMFILE:0x%06X with no error", cpu, off)
		}
		if win == WindowNone {
			t.Fatalf("ROMCPU:$%06X succeeded but reported WindowNone", cpu)
		}
	})
}
