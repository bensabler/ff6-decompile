package mesenstate

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

// u16 and u8 build the little-endian scalar blocks Mesen writes.
func u16(v uint16) []byte {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	return b[:]
}

func u8(v uint8) []byte { return []byte{v} }

func flag(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

// hardwareSection builds a savestate carrying every memory image and every
// scalar the PPU and DMA decoders read. All values are project-authored, so
// this fixture contains no ROM-derived data and is safe in CI.
func hardwareSection() []byte {
	var b bytes.Buffer
	b.Write(block("cpu.k", []byte{0xC0}))
	b.Write(block(WorkRAMBlock, fill(WorkRAMSize, 1)))
	b.Write(block(VideoRAMBlock, fill(VideoRAMSize, 3)))
	b.Write(block(CGRAMBlock, fill(CGRAMSize, 4)))
	b.Write(block(OAMBlock, fill(OAMSize, 5)))
	b.Write(block(AudioRAMBlock, fill(AudioRAMSize, 6)))

	b.Write(block("ppu.bgMode", u8(1)))
	b.Write(block("ppu.screenBrightness", u8(15)))
	b.Write(block("ppu.forcedBlank", flag(false)))
	b.Write(block("ppu.mainScreenLayers", u8(0x17)))
	b.Write(block("ppu.subScreenLayers", u8(0x13)))
	b.Write(block("ppu.oamMode", u8(3)))
	b.Write(block("ppu.oamBaseAddress", u16(0x6000)))

	// Word addresses, which the decoder doubles into byte offsets.
	chr := []uint16{0x0000, 0x0000, 0x3000, 0x0000}
	tilemap := []uint16{0x4C00, 0x5400, 0x5800, 0x5400}
	for i := 0; i < LayerCount; i++ {
		p := "ppu.layers[" + string(rune('0'+i)) + "]."
		b.Write(block(p+"chrAddress", u16(chr[i])))
		b.Write(block(p+"tilemapAddress", u16(tilemap[i])))
		b.Write(block(p+"hscroll", u16(uint16(100*(i+1)))))
		b.Write(block(p+"vscroll", u16(uint16(10*(i+1)))))
		b.Write(block(p+"largeTiles", flag(i == 2)))
		b.Write(block(p+"doubleWidth", flag(i == 0)))
		b.Write(block(p+"doubleHeight", flag(false)))
	}

	dest := []uint8{0x18, 0x22, 0x04, 0x30, 0x30, 0x30, 0x30, 0x30}
	for i := 0; i < DMAChannelCount; i++ {
		p := "dmaController.channel[" + string(rune('0'+i)) + "]."
		b.Write(block(p+"destAddress", u8(dest[i])))
		b.Write(block(p+"transferMode", u8(uint8(i%8))))
		b.Write(block(p+"srcBank", u8(0x7E)))
		b.Write(block(p+"hdmaBank", u8(0xC0)))
		b.Write(block(p+"srcAddress", u16(uint16(0xC180+i))))
		b.Write(block(p+"transferSize", u16(uint16(0x1000*i))))
		b.Write(block(p+"hdmaTableAddress", u16(0x2000)))
		b.Write(block(p+"fixedTransfer", flag(false)))
		b.Write(block(p+"decrement", flag(false)))
		b.Write(block(p+"invertDirection", flag(false)))
		b.Write(block(p+"dmaActive", flag(i == 0)))
		b.Write(block(p+"hdmaIndirectAddressing", flag(false)))
		b.Write(block(p+"hdmaFinished", flag(true)))
	}
	return b.Bytes()
}

func hardwareState(t *testing.T) *State {
	t.Helper()
	st, err := Parse(savestate(t, hardwareSection()))
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestMemoryImages(t *testing.T) {
	st := hardwareState(t)

	tests := []struct {
		name string
		read func() ([]byte, error)
		want int
	}{
		{"WorkRAM", st.WorkRAM, WorkRAMSize},
		{"VideoRAM", st.VideoRAM, VideoRAMSize},
		{"CGRAM", st.CGRAM, CGRAMSize},
		{"OAM", st.OAM, OAMSize},
		{"AudioRAM", st.AudioRAM, AudioRAMSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.read()
			if err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}
			if len(b) != tt.want {
				t.Errorf("%s is %d bytes, want %d", tt.name, len(b), tt.want)
			}
		})
	}
}

// TestMissingImagesReportTheBlockName matters because the corpus is not
// uniform: an older savestate may lack a block, and the error has to say
// which one rather than returning a short slice.
func TestMissingImagesReportTheBlockName(t *testing.T) {
	st, err := Parse(savestate(t, block(WorkRAMBlock, fill(WorkRAMSize, 1))))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		read  func() ([]byte, error)
		block string
	}{
		{st.VideoRAM, VideoRAMBlock},
		{st.CGRAM, CGRAMBlock},
		{st.OAM, OAMBlock},
		{st.AudioRAM, AudioRAMBlock},
	} {
		_, err := tc.read()
		if err == nil {
			t.Errorf("%s: want an error when the block is absent", tc.block)
			continue
		}
		if !strings.Contains(err.Error(), tc.block) {
			t.Errorf("error %q does not name the missing block %q", err, tc.block)
		}
	}
}

func TestScalarReads(t *testing.T) {
	st := hardwareState(t)

	if v, err := st.Uint8("ppu.bgMode"); err != nil || v != 1 {
		t.Errorf("Uint8(ppu.bgMode) = %d, %v; want 1, nil", v, err)
	}
	if v, err := st.Uint16("ppu.oamBaseAddress"); err != nil || v != 0x6000 {
		t.Errorf("Uint16(ppu.oamBaseAddress) = $%04X, %v; want $6000, nil", v, err)
	}
	if v, err := st.Bool("ppu.forcedBlank"); err != nil || v {
		t.Errorf("Bool(ppu.forcedBlank) = %v, %v; want false, nil", v, err)
	}

	// A width mismatch must be an error, not a silent truncation: reading a
	// two-byte field as one byte would produce a plausible wrong address.
	if _, err := st.Uint8("ppu.oamBaseAddress"); err == nil {
		t.Error("reading a 2-byte block as Uint8 should fail")
	}
	if _, err := st.Uint16("ppu.bgMode"); err == nil {
		t.Error("reading a 1-byte block as Uint16 should fail")
	}
	if _, err := st.Uint32("nope"); err == nil {
		t.Error("reading an absent block should fail")
	}
}

func TestPPUState(t *testing.T) {
	st := hardwareState(t)
	p, err := st.PPUState()
	if err != nil {
		t.Fatal(err)
	}

	if p.BGMode != 1 || p.Brightness != 15 || p.ForcedBlank {
		t.Errorf("mode/brightness/blank = %d/%d/%v", p.BGMode, p.Brightness, p.ForcedBlank)
	}
	// Word address $6000 is byte offset $C000. Getting this wrong halves
	// every address and points the reader at the wrong half of VRAM.
	if p.OAMBaseAddress != 0xC000 {
		t.Errorf("OAMBaseAddress = $%04X, want $C000", p.OAMBaseAddress)
	}
	if p.Layers[2].ChrAddress != 0x6000 {
		t.Errorf("BG3 ChrAddress = $%04X, want $6000 (word $3000 doubled)", p.Layers[2].ChrAddress)
	}
	if p.Layers[1].TilemapAddress != 0xA800 {
		t.Errorf("BG2 TilemapAddress = $%04X, want $A800", p.Layers[1].TilemapAddress)
	}
	if !p.Layers[2].LargeTiles || p.Layers[0].LargeTiles {
		t.Error("largeTiles decoded onto the wrong layer")
	}
	if p.Layers[0].HScroll != 100 || p.Layers[3].VScroll != 40 {
		t.Errorf("scroll decoded wrong: BG1 h=%d BG4 v=%d", p.Layers[0].HScroll, p.Layers[3].VScroll)
	}

	// $17 is BG1, BG2, BG3 and sprites; BG4 is off.
	for i, want := range []bool{true, true, true, false} {
		if got := p.LayerEnabled(i); got != want {
			t.Errorf("LayerEnabled(%d) = %v, want %v (main screen $%02X)", i, got, want, p.MainScreenLayers)
		}
	}
	if p.LayerEnabled(-1) || p.LayerEnabled(LayerCount) {
		t.Error("LayerEnabled must reject out-of-range layers rather than index memory")
	}
}

func TestDMAChannels(t *testing.T) {
	st := hardwareState(t)
	ch, err := st.DMAChannels()
	if err != nil {
		t.Fatal(err)
	}
	if len(ch) != DMAChannelCount {
		t.Fatalf("got %d channels, want %d", len(ch), DMAChannelCount)
	}

	if got := ch[0].FullSource(); got != 0x7EC180 {
		t.Errorf("channel 0 FullSource = $%06X, want $7EC180", got)
	}
	if !ch[0].TargetsVRAM() || ch[0].TargetsCGRAM() || ch[0].TargetsOAM() {
		t.Error("channel 0 writes $2118 and should classify as VRAM only")
	}
	if !ch[1].TargetsCGRAM() {
		t.Error("channel 1 writes $2122 and should classify as CGRAM")
	}
	if !ch[2].TargetsOAM() {
		t.Error("channel 2 writes $2104 and should classify as OAM")
	}
	if ch[3].TargetsVRAM() || ch[3].TargetsCGRAM() || ch[3].TargetsOAM() {
		t.Error("channel 3 writes $2130 and should classify as none of the three")
	}
	if !ch[0].Active || ch[1].Active {
		t.Error("dmaActive decoded onto the wrong channel")
	}
	for i, c := range ch {
		if c.Index != i {
			t.Errorf("channel %d reports Index %d", i, c.Index)
		}
	}
}

func TestDMAChannelsReportMissingFields(t *testing.T) {
	st, err := Parse(savestate(t, block(WorkRAMBlock, fill(WorkRAMSize, 1))))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.DMAChannels(); err == nil {
		t.Error("want an error when the channel blocks are absent")
	} else if !strings.Contains(err.Error(), "channel 0") {
		t.Errorf("error %q should name the channel it failed on", err)
	}
	if _, err := st.PPUState(); err == nil {
		t.Error("want an error when the PPU blocks are absent")
	}
}
