package validate

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/bgr555"
)

func fbOf(fill uint8, set map[int]uint8) *framebuf.Indexed {
	fb := framebuf.New()
	fb.Fill(fill)
	for i, v := range set {
		fb.Pix[i] = v
	}
	return fb
}

func TestCompareIndexed(t *testing.T) {
	a := fbOf(0, map[int]uint8{0: 3, 100: 2})
	b := fbOf(0, map[int]uint8{0: 3, 100: 2})
	if d := CompareIndexed(a, b); !d.Equal() {
		t.Errorf("identical planes reported as different: %s", d)
	}

	c := fbOf(0, map[int]uint8{0: 3, 100: 7})
	d := CompareIndexed(a, c)
	if d.Equal() || d.Pixels != 1 {
		t.Fatalf("want 1 differing pixel, got %s", d)
	}
	if d.FirstX != 100 || d.FirstY != 0 {
		t.Errorf("first difference at (%d,%d), want (100,0)", d.FirstX, d.FirstY)
	}
	// The histogram is what turns a pixel count into a diagnosis.
	if d.Histogram[[2]uint8{2, 7}] != 1 {
		t.Errorf("histogram should record 2->7, got %v", d.Histogram)
	}
	if d.Total != framebuf.Width*framebuf.Height {
		t.Errorf("Total = %d", d.Total)
	}
}

// The load-bearing test. Two frames with different indices that resolve to the
// same colours are identical to a viewer, and CompareIndexed must disagree
// with CompareResolved about that.
func TestIndexedAndResolvedAnswerDifferentQuestions(t *testing.T) {
	white := bgr555.Color{R: 0xFF, G: 0xFF, B: 0xFF}

	var palA framebuf.Palette
	palA[3] = white
	var palB framebuf.Palette
	palB[7] = white

	a := fbOf(0, map[int]uint8{50: 3})
	b := fbOf(0, map[int]uint8{50: 7})

	if CompareIndexed(a, b).Equal() {
		t.Error("CompareIndexed should see indices 3 and 7 as different")
	}
	if d := CompareResolved(a, &palA, b, &palB); !d.Equal() {
		t.Errorf("CompareResolved should see identical colours: %s", d)
	}
}

// The D0 case, as a test. A frame whose ink lands on undefined entries hashes
// fine and renders blank; only the resolved comparison can say so.
func TestResolvedCatchesInvisibleInk(t *testing.T) {
	var pal framebuf.Palette // every entry black, including index 0
	pal[3] = bgr555.Color{R: 0xFF, G: 0xFF, B: 0xFF}

	visible := fbOf(0, map[int]uint8{10: 3, 11: 3, 12: 3})
	invisible := fbOf(0, map[int]uint8{10: 9, 11: 9, 12: 9}) // entry 9 is black

	if a, b := visible.Sum256(), invisible.Sum256(); a == b {
		t.Fatal("the two frames should have different index hashes")
	}

	dv := Visibility(visible, &pal)
	if dv.Drawn != 3 || dv.Visible != 3 || dv.VisibleFraction() != 1 {
		t.Errorf("visible frame: %s", dv)
	}

	di := Visibility(invisible, &pal)
	if di.Drawn != 3 {
		t.Errorf("invisible frame should still count 3 drawn pixels, got %d", di.Drawn)
	}
	if di.Visible != 0 || di.VisibleFraction() != 0 {
		t.Errorf("invisible frame should have zero visible ink, got %s", di)
	}
	// And the point: comparing it to itself by index says nothing is wrong.
	if !CompareIndexed(invisible, invisible).Equal() {
		t.Error("an index self-comparison cannot fail, which is why it is not enough")
	}
}

func TestPaletteUse(t *testing.T) {
	fb := fbOf(0, map[int]uint8{1: 3, 2: 14, 3: 7})
	high, distinct := PaletteUse(fb)
	if high != 14 {
		t.Errorf("highest = %d, want 14", high)
	}
	if distinct != 4 { // 0, 3, 7, 14
		t.Errorf("distinct = %d, want 4", distinct)
	}
	if h, d := PaletteUse(nil); h != 0 || d != 0 {
		t.Errorf("nil frame should report 0,0, got %d,%d", h, d)
	}
}

func TestDecodeCGRAM(t *testing.T) {
	cgram := make([]byte, 512)
	// Entry 1 = BGR555 $7FFF, white.
	cgram[2], cgram[3] = 0xFF, 0x7F
	pal, err := DecodeCGRAM(cgram)
	if err != nil {
		t.Fatal(err)
	}
	if got := pal[1]; got.R != 0xFF || got.G != 0xFF || got.B != 0xFF {
		t.Errorf("entry 1 = %+v, want white", got)
	}
	if got := pal[0]; got.R != 0 || got.G != 0 || got.B != 0 {
		t.Errorf("entry 0 = %+v, want black", got)
	}

	for _, n := range []int{0, 256, 511, 513} {
		if _, err := DecodeCGRAM(make([]byte, n)); err == nil {
			t.Errorf("a %d-byte CGRAM should be rejected", n)
		}
	}
}

func TestComparisonsHandleNil(t *testing.T) {
	fb := framebuf.New()
	pal := framebuf.GrayPalette()
	if d := CompareIndexed(nil, fb); d.Total != 0 {
		t.Error("nil input should compare nothing rather than panic")
	}
	if d := CompareResolved(fb, nil, fb, pal); d.Total != 0 {
		t.Error("nil palette should compare nothing rather than panic")
	}
	if d := CompareResolved(nil, pal, nil, pal); !d.Equal() {
		t.Error("nil frames should not report differences")
	}
}
