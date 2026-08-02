package mesenstate

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/platform/snesdma"
)

// This file decodes the scalar hardware state a savestate carries alongside
// its memory images.
//
// The memory images alone are not interpretable. A 64 KB VRAM dump does not
// say which of its bytes are tiles and which are tilemap, at what depth, or
// where each background layer reads from — the PPU registers say that, and
// they are in the same file. Likewise a DMA channel's source bank and address
// name the buffer a transfer came from, which is the first step backward from
// "these bytes were in VRAM" to "this ROM span produced them".
//
// Live DMA *tracing* does not exist in this project: no Lua probe and no Go
// path reads a DMA register during execution, despite /trace-dma, the
// dma-tracer skill, the dma-researcher agent and TRACE_DMA.md all existing.
// What is here is weaker and worth naming precisely — **channel configuration
// at the instant the state was written**, not a record of transfers over
// time. A channel that finished before the capture has whatever its last
// setup left behind, and a channel that never ran holds initialisation
// values. Treat it as a lead, never as a trace.

// Uint8 reads a one-byte block.
func (s *State) Uint8(name string) (uint8, error) {
	b, err := s.scalar(name, 1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// Uint16 reads a two-byte little-endian block.
func (s *State) Uint16(name string) (uint16, error) {
	b, err := s.scalar(name, 2)
	if err != nil {
		return 0, err
	}
	return uint16(b[0]) | uint16(b[1])<<8, nil
}

// Uint32 reads a four-byte little-endian block.
func (s *State) Uint32(name string) (uint32, error) {
	b, err := s.scalar(name, 4)
	if err != nil {
		return 0, err
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

// Bool reads a one-byte block as a flag. Mesen stores these as 0 or 1.
func (s *State) Bool(name string) (bool, error) {
	v, err := s.Uint8(name)
	return v != 0, err
}

// scalar returns a block of exactly want bytes.
func (s *State) scalar(name string, want int) ([]byte, error) {
	b, ok := s.blocks[name]
	if !ok {
		return nil, fmt.Errorf("mesenstate: savestate has no %q block", name)
	}
	if len(b) != want {
		return nil, fmt.Errorf("mesenstate: %q is %d bytes, want %d", name, len(b), want)
	}
	return b, nil
}

// LayerCount is the number of SNES background layers.
const LayerCount = 4

// Layer is one background layer's configuration.
//
// ChrAddress and TilemapAddress are byte offsets into the VRAM image, already
// converted from the word addresses the registers hold, so they index
// VideoRAM() directly.
type Layer struct {
	ChrAddress     int
	TilemapAddress int
	HScroll        uint16
	VScroll        uint16
	LargeTiles     bool
	DoubleWidth    bool
	DoubleHeight   bool
}

// PPU is the subset of PPU configuration needed to interpret a VRAM image.
type PPU struct {
	BGMode           uint8
	Brightness       uint8
	ForcedBlank      bool
	MainScreenLayers uint8
	SubScreenLayers  uint8
	OAMMode          uint8
	OAMBaseAddress   int
	Layers           [LayerCount]Layer
}

// PPUState decodes the PPU configuration.
func (s *State) PPUState() (*PPU, error) {
	var p PPU
	var err error

	get8 := func(name string, dst *uint8) {
		if err != nil {
			return
		}
		*dst, err = s.Uint8(name)
	}
	getBool := func(name string, dst *bool) {
		if err != nil {
			return
		}
		*dst, err = s.Bool(name)
	}
	get16 := func(name string, dst *uint16) {
		if err != nil {
			return
		}
		*dst, err = s.Uint16(name)
	}

	get8("ppu.bgMode", &p.BGMode)
	get8("ppu.screenBrightness", &p.Brightness)
	getBool("ppu.forcedBlank", &p.ForcedBlank)
	get8("ppu.mainScreenLayers", &p.MainScreenLayers)
	get8("ppu.subScreenLayers", &p.SubScreenLayers)
	get8("ppu.oamMode", &p.OAMMode)

	var oamBase uint16
	get16("ppu.oamBaseAddress", &oamBase)
	if err != nil {
		return nil, fmt.Errorf("decoding PPU state: %w", err)
	}
	p.OAMBaseAddress = int(oamBase) * 2

	for i := range p.Layers {
		var chr, tilemap uint16
		get16(fmt.Sprintf("ppu.layers[%d].chrAddress", i), &chr)
		get16(fmt.Sprintf("ppu.layers[%d].tilemapAddress", i), &tilemap)
		get16(fmt.Sprintf("ppu.layers[%d].hscroll", i), &p.Layers[i].HScroll)
		get16(fmt.Sprintf("ppu.layers[%d].vscroll", i), &p.Layers[i].VScroll)
		getBool(fmt.Sprintf("ppu.layers[%d].largeTiles", i), &p.Layers[i].LargeTiles)
		getBool(fmt.Sprintf("ppu.layers[%d].doubleWidth", i), &p.Layers[i].DoubleWidth)
		getBool(fmt.Sprintf("ppu.layers[%d].doubleHeight", i), &p.Layers[i].DoubleHeight)
		if err != nil {
			return nil, fmt.Errorf("decoding PPU layer %d: %w", i, err)
		}
		// Mesen stores these as VRAM word addresses; VideoRAM() is bytes.
		p.Layers[i].ChrAddress = int(chr) * 2
		p.Layers[i].TilemapAddress = int(tilemap) * 2
	}
	return &p, nil
}

// LayerEnabled reports whether layer i appears on the main screen.
func (p *PPU) LayerEnabled(i int) bool {
	if i < 0 || i >= LayerCount {
		return false
	}
	return p.MainScreenLayers&(1<<uint(i)) != 0
}

// DMAChannelCount is the number of SNES DMA/HDMA channels.
const DMAChannelCount = 8

// DMAChannel is one channel's configuration at capture time.
//
// It is **not** a transfer log. See the note at the top of this file: a
// finished channel retains its last setup, and an unused channel retains
// initialisation values. SrcBank and SrcAddress are a lead toward the buffer
// or ROM span a VRAM region came from, to be confirmed by other evidence.
type DMAChannel struct {
	Index int
	// DestAddress is the low byte of the $21xx register written, so $18 is
	// $2118, the VRAM data port.
	DestAddress     uint8
	TransferMode    uint8
	SrcBank         uint8
	SrcAddress      uint16
	TransferSize    uint16
	FixedTransfer   bool
	Decrement       bool
	InvertDirection bool
	Active          bool
	// HDMA fields. HDMATableAddress is meaningful only for channels the
	// $420C enable register turned on, which this struct does not carry.
	HDMABank         uint8
	HDMATableAddress uint16
	HDMAIndirect     bool
	HDMAFinished     bool
}

// FullSource returns the 24-bit source address as the project writes CPU
// addresses, bank in the high byte.
func (c DMAChannel) FullSource() uint32 {
	return uint32(c.SrcBank)<<16 | uint32(c.SrcAddress)
}

// ByteCount returns how many bytes the configured transfer would move.
//
// TransferSize is the raw register pair, and **zero means 65536**: the
// hardware decrements the counter before testing it. Reading the raw value is
// how a 64 KB VRAM upload gets recorded as an empty transfer, which is exactly
// how the mines savestate's channel 0 read before this method existed.
func (c DMAChannel) ByteCount() int {
	if c.TransferSize == 0 {
		return 65536
	}
	return int(c.TransferSize)
}

// TargetsVRAM reports whether the channel was configured to write the VRAM
// data port ($2118/$2119).
func (c DMAChannel) TargetsVRAM() bool {
	return c.DestAddress == snesdma.DestVRAMLow || c.DestAddress == snesdma.DestVRAMHigh
}

// TargetsCGRAM reports whether the channel was configured to write the
// palette data port ($2122).
func (c DMAChannel) TargetsCGRAM() bool { return c.DestAddress == snesdma.DestCGRAM }

// TargetsOAM reports whether the channel was configured to write the sprite
// data port ($2104).
func (c DMAChannel) TargetsOAM() bool { return c.DestAddress == snesdma.DestOAM }

// DMAChannels decodes all eight channel configurations.
func (s *State) DMAChannels() ([]DMAChannel, error) {
	out := make([]DMAChannel, 0, DMAChannelCount)
	for i := 0; i < DMAChannelCount; i++ {
		c := DMAChannel{Index: i}
		var err error
		field := func(n string) string { return fmt.Sprintf("dmaController.channel[%d].%s", i, n) }

		type u8 struct {
			name string
			dst  *uint8
		}
		for _, f := range []u8{
			{"destAddress", &c.DestAddress},
			{"transferMode", &c.TransferMode},
			{"srcBank", &c.SrcBank},
			{"hdmaBank", &c.HDMABank},
		} {
			if c2, e := s.Uint8(field(f.name)); e != nil {
				err = e
				break
			} else {
				*f.dst = c2
			}
		}
		if err != nil {
			return nil, fmt.Errorf("decoding DMA channel %d: %w", i, err)
		}

		type u16 struct {
			name string
			dst  *uint16
		}
		for _, f := range []u16{
			{"srcAddress", &c.SrcAddress},
			{"transferSize", &c.TransferSize},
			{"hdmaTableAddress", &c.HDMATableAddress},
		} {
			if v, e := s.Uint16(field(f.name)); e != nil {
				err = e
				break
			} else {
				*f.dst = v
			}
		}
		if err != nil {
			return nil, fmt.Errorf("decoding DMA channel %d: %w", i, err)
		}

		type ub struct {
			name string
			dst  *bool
		}
		for _, f := range []ub{
			{"fixedTransfer", &c.FixedTransfer},
			{"decrement", &c.Decrement},
			{"invertDirection", &c.InvertDirection},
			{"dmaActive", &c.Active},
			{"hdmaIndirectAddressing", &c.HDMAIndirect},
			{"hdmaFinished", &c.HDMAFinished},
		} {
			if v, e := s.Bool(field(f.name)); e != nil {
				err = e
				break
			} else {
				*f.dst = v
			}
		}
		if err != nil {
			return nil, fmt.Errorf("decoding DMA channel %d: %w", i, err)
		}
		out = append(out, c)
	}
	return out, nil
}
