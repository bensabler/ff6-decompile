// Package snesdma decodes SNES DMA and HDMA channel registers.
//
// # Why this package exists
//
// The project has carried DMA *doctrine* since early on — a `/trace-dma`
// command, a `dma-tracer` skill, a `dma-researcher` agent and a TRACE_DMA
// playbook — with **no implementation beneath any of it**. Nothing read a DMA
// register, in Lua or in Go. The savestate reader added in Unit 12 gets close
// but only sees channel state at one instant, which EXP-0050 showed can be the
// leftovers of a finished transfer rather than the one that mattered.
//
// A real trace needs two halves: a probe that captures the registers when a
// transfer is kicked, and something that turns those bytes into a statement
// about what moved where. This package is the second half, and it is the half
// that can be tested without an emulator.
//
// # Address domain
//
// Channel n's registers occupy $43n0-$43nA in the CPU's I/O page. Source
// addresses are 24-bit CPU addresses (bank in A1Bn); use platform/snesaddr to
// convert a ROM source to ROMFILE. Destination is the low byte of a $21xx
// register, so $18 means $2118, the VRAM data port.
package snesdma

import "fmt"

// RegisterCount is how many registers per channel this package decodes,
// $43n0 through $43nA.
const RegisterCount = 11

// ChannelCount is the number of DMA/HDMA channels the SNES provides.
const ChannelCount = 8

// Well-known destination registers, as the low byte of $21xx.
const (
	DestVRAMLow  uint8 = 0x18 // $2118
	DestVRAMHigh uint8 = 0x19 // $2119
	DestCGRAM    uint8 = 0x22 // $2122
	DestOAM      uint8 = 0x04 // $2104
)

// Direction of a transfer, as bit 7 of the control register.
type Direction uint8

const (
	// AToB is CPU memory to the PPU/APU register, the ordinary case.
	AToB Direction = iota
	// BToA reads a register back into CPU memory.
	BToA
)

func (d Direction) String() string {
	if d == BToA {
		return "B->A"
	}
	return "A->B"
}

// Channel is one decoded DMA channel configuration.
type Channel struct {
	Index int

	// TransferMode selects the register write pattern, control bits 0-2.
	TransferMode uint8
	// Fixed holds the source address still instead of stepping it.
	Fixed bool
	// Decrement steps the source address backward. Ignored when Fixed.
	Decrement bool
	// Direction is which way bytes move.
	Direction Direction
	// HDMAIndirect selects indirect HDMA addressing, control bit 6.
	HDMAIndirect bool

	// Dest is the low byte of the $21xx register written.
	Dest uint8

	// SrcBank and SrcAddr form the 24-bit source.
	SrcBank uint8
	SrcAddr uint16

	// Size is the DMA byte count. Zero means 65536 — the hardware treats
	// the counter as "decrement then test", so a zero count transfers the
	// whole 64 KB. Use ByteCount rather than reading this directly.
	Size uint16

	// HDMA table fields. Meaningful only for channels $420C enabled.
	HDMAIndirectBank uint8
	HDMATableAddr    uint16
	HDMALineCounter  uint8
}

// Decode reads one channel's registers, $43n0 through $43nA in order.
func Decode(index int, regs []byte) (Channel, error) {
	if index < 0 || index >= ChannelCount {
		return Channel{}, fmt.Errorf("snesdma: channel %d out of range 0-%d", index, ChannelCount-1)
	}
	if len(regs) < RegisterCount {
		return Channel{}, fmt.Errorf("snesdma: channel %d: got %d registers, want %d", index, len(regs), RegisterCount)
	}
	ctrl := regs[0]
	c := Channel{
		Index:            index,
		TransferMode:     ctrl & 0x07,
		Fixed:            ctrl&0x08 != 0,
		Decrement:        ctrl&0x10 != 0,
		HDMAIndirect:     ctrl&0x40 != 0,
		Dest:             regs[1],
		SrcAddr:          uint16(regs[2]) | uint16(regs[3])<<8,
		SrcBank:          regs[4],
		Size:             uint16(regs[5]) | uint16(regs[6])<<8,
		HDMAIndirectBank: regs[7],
		HDMATableAddr:    uint16(regs[8]) | uint16(regs[9])<<8,
		HDMALineCounter:  regs[10],
	}
	if ctrl&0x80 != 0 {
		c.Direction = BToA
	}
	return c, nil
}

// Source returns the 24-bit CPU source address.
func (c Channel) Source() uint32 { return uint32(c.SrcBank)<<16 | uint32(c.SrcAddr) }

// ByteCount returns how many bytes the transfer moves.
//
// A Size of zero means 65536, not zero: the hardware decrements the counter
// before testing it. Reading Size directly is the classic way to record a
// 64 KB VRAM upload as an empty transfer, which is exactly what the mines
// savestate's channel 0 looks like at a glance.
func (c Channel) ByteCount() int {
	if c.Size == 0 {
		return 65536
	}
	return int(c.Size)
}

// TransferPattern returns the destination-register offsets the mode writes,
// in order, repeating until the byte count is exhausted.
//
// Mode 0 writes one register once; mode 1 alternates two ($2118/$2119, the
// pattern a VRAM upload uses); mode 4 walks four. Modes 6 and 7 duplicate 2
// and 3 on real hardware.
func TransferPattern(mode uint8) []uint8 {
	switch mode & 0x07 {
	case 0:
		return []uint8{0}
	case 1:
		return []uint8{0, 1}
	case 2, 6:
		return []uint8{0, 0}
	case 3, 7:
		return []uint8{0, 0, 1, 1}
	case 4:
		return []uint8{0, 1, 2, 3}
	case 5:
		return []uint8{0, 1, 0, 1}
	}
	return nil // unreachable: mode is masked to 0-7
}

// TargetsVRAM reports whether the channel writes the VRAM data port.
func (c Channel) TargetsVRAM() bool { return c.Dest == DestVRAMLow || c.Dest == DestVRAMHigh }

// TargetsCGRAM reports whether the channel writes the palette data port.
func (c Channel) TargetsCGRAM() bool { return c.Dest == DestCGRAM }

// TargetsOAM reports whether the channel writes the sprite data port.
func (c Channel) TargetsOAM() bool { return c.Dest == DestOAM }

// SourceSpan returns the inclusive CPU address range the transfer reads, and
// whether that range is a contiguous span at all.
//
// ok is false in three cases, each for a different reason:
//
//   - **Fixed** reads one address repeatedly, so there is no span.
//   - **B->A** writes CPU memory rather than reading it.
//   - **The transfer crosses a bank boundary.** A1Tn is 16 bits and the bank
//     register does not increment, so a transfer running past $FFFF wraps to
//     $0000 *within the same bank*. The bytes it reads are therefore not
//     contiguous, and reporting a span that runs into the next bank would
//     point a provenance search at a ROM region the hardware never touched.
//
// Returning ok=false rather than a plausible-looking range is the point: a
// span is what a provenance search consumes, and a wrong one is worse than
// none.
func (c Channel) SourceSpan() (start, end uint32, ok bool) {
	if c.Fixed || c.Direction == BToA {
		return 0, 0, false
	}
	n := uint32(c.ByteCount())
	bank := uint32(c.SrcBank) << 16
	off := uint32(c.SrcAddr)
	if c.Decrement {
		if n-1 > off {
			return 0, 0, false // wraps back to $FFFF within the bank
		}
		return bank | (off - (n - 1)), bank | off, true
	}
	if off+(n-1) > 0xFFFF {
		return 0, 0, false // wraps forward to $0000 within the bank
	}
	return bank | off, bank | (off + (n - 1)), true
}

// Enabled reports which channels a $420B (MDMAEN) or $420C (HDMAEN) write
// turns on, lowest channel first.
func Enabled(mask uint8) []int {
	var out []int
	for i := 0; i < ChannelCount; i++ {
		if mask&(1<<uint(i)) != 0 {
			out = append(out, i)
		}
	}
	return out
}
