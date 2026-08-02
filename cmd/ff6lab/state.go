package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/mesenstate"
)

const stateRegions = `wram|sram|vram|cgram|oam|aram`

const stateUsage = `usage:
  ff6lab state list <file.mss>
  ff6lab state ppu <file.mss>
  ff6lab state ` + stateRegions + ` <file.mss> [-o <out.bin>]
  ff6lab state read <file.mss> ` + stateRegions + ` <hexaddr> <length>
  ff6lab state diff <a.mss> <b.mss> [` + stateRegions + `]`

// runState reads emulated memory out of preserved Mesen savestates so
// earlier captures can be re-examined without an emulator.
func runState(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("%s", stateUsage)
	}
	switch args[0] {
	case "list":
		return stateList(args[1], out)
	case "ppu":
		return statePPU(args[1], out)
	case "wram", "sram", "vram", "cgram", "oam", "aram":
		rest, dest := takeOutputFlag(args[2:])
		if len(rest) != 0 {
			return fmt.Errorf("%s", stateUsage)
		}
		return stateDump(args[1], args[0], dest, out)
	case "read":
		if len(args) != 5 {
			return fmt.Errorf("%s", stateUsage)
		}
		return stateRead(args[1], args[2], args[3], args[4], out)
	case "diff":
		if len(args) < 3 || len(args) > 4 {
			return fmt.Errorf("%s", stateUsage)
		}
		region := "wram"
		if len(args) == 4 {
			region = args[3]
		}
		return stateDiff(args[1], args[2], region, out)
	default:
		return fmt.Errorf("%s", stateUsage)
	}
}

// takeOutputFlag pulls an "-o <path>" pair out of args.
func takeOutputFlag(args []string) (rest []string, dest string) {
	for i := 0; i < len(args); i++ {
		if args[i] == "-o" && i+1 < len(args) {
			dest = args[i+1]
			i++
			continue
		}
		rest = append(rest, args[i])
	}
	return rest, dest
}

// regionNames are the memory images `state` can read, in the order the
// usage text lists them.
var regionNames = []string{"wram", "sram", "vram", "cgram", "oam", "aram"}

// region returns the named memory image and the address prefix the
// project uses to write addresses inside it.
func region(st *mesenstate.State, name string) ([]byte, string, error) {
	switch strings.ToLower(name) {
	case "wram":
		b, err := st.WorkRAM()
		return b, "WRAM:+$", err
	case "sram":
		b, err := st.SaveRAM()
		return b, "SRAM:+$", err
	case "vram":
		b, err := st.VideoRAM()
		return b, "VRAM:$", err
	case "cgram":
		b, err := st.CGRAM()
		return b, "CGRAM:$", err
	case "oam":
		b, err := st.OAM()
		return b, "OAM:$", err
	case "aram":
		b, err := st.AudioRAM()
		return b, "ARAM:$", err
	default:
		return nil, "", fmt.Errorf("unknown region %q; want one of %s",
			name, strings.Join(regionNames, ", "))
	}
}

func stateList(path string, out io.Writer) error {
	st, err := mesenstate.Open(path)
	if err != nil {
		return err
	}
	for _, n := range st.Names() {
		b, _ := st.Block(n)
		fmt.Fprintf(out, "%-48s %8d\n", n, len(b))
	}
	return nil
}

// statePPU prints the configuration needed to interpret a VRAM image, plus
// the DMA channel setup that hints at where those bytes came from.
//
// A VRAM dump on its own says nothing about which bytes are tiles, at what
// depth, or which layer reads them. This is the other half of the evidence,
// and it is in the same file.
func statePPU(path string, out io.Writer) error {
	st, err := mesenstate.Open(path)
	if err != nil {
		return err
	}
	p, err := st.PPUState()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "BG mode %d   brightness %d/15   forced blank %v\n", p.BGMode, p.Brightness, p.ForcedBlank)
	fmt.Fprintf(out, "main screen layers $%02X   sub screen layers $%02X   OAM mode %d, base VRAM:$%04X\n",
		p.MainScreenLayers, p.SubScreenLayers, p.OAMMode, p.OAMBaseAddress)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "layer  on   chr       tilemap   scroll        tiles  size")
	for i, l := range p.Layers {
		on := " "
		if p.LayerEnabled(i) {
			on = "*"
		}
		tiles := "8x8"
		if l.LargeTiles {
			tiles = "16x16"
		}
		fmt.Fprintf(out, "BG%d    %s    VRAM:$%04X VRAM:$%04X (%5d,%5d) %-6s %dx%d screens\n",
			i+1, on, l.ChrAddress, l.TilemapAddress, l.HScroll, l.VScroll, tiles,
			1+b2i(l.DoubleWidth), 1+b2i(l.DoubleHeight))
	}

	ch, err := st.DMAChannels()
	if err != nil {
		return err
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, "DMA channel configuration at capture time. This is NOT a transfer log:")
	fmt.Fprintln(out, "a finished channel keeps its last setup and an unused one keeps its")
	fmt.Fprintln(out, "initialisation values. Treat a source address as a lead, not a trace.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "ch  dest    mode  source        size   target")
	for _, c := range ch {
		target := ""
		switch {
		case c.TargetsVRAM():
			target = "VRAM"
		case c.TargetsCGRAM():
			target = "CGRAM"
		case c.TargetsOAM():
			target = "OAM"
		}
		fmt.Fprintf(out, "%d   $21%02X   %d     $%06X      %5d  %s\n",
			c.Index, c.DestAddress, c.TransferMode, c.FullSource(), c.TransferSize, target)
	}
	return nil
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func stateDump(path, name, dest string, out io.Writer) error {
	st, err := mesenstate.Open(path)
	if err != nil {
		return err
	}
	b, _, err := region(st, name)
	if err != nil {
		return err
	}
	if dest == "" {
		return fmt.Errorf("%s: -o <out.bin> is required (the image is %d bytes of binary)", name, len(b))
	}
	if err := os.WriteFile(dest, b, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dest, err)
	}
	fmt.Fprintf(out, "%s: %d bytes -> %s\n", name, len(b), dest)
	return nil
}

func stateRead(path, name, addr, length string, out io.Writer) error {
	st, err := mesenstate.Open(path)
	if err != nil {
		return err
	}
	b, prefix, err := region(st, name)
	if err != nil {
		return err
	}
	start, err := parseHex(addr)
	if err != nil {
		return fmt.Errorf("address %q: %w", addr, err)
	}
	n, err := strconv.Atoi(length)
	if err != nil || n <= 0 {
		return fmt.Errorf("length %q: want a positive decimal count", length)
	}
	if start+n > len(b) {
		return fmt.Errorf("%s%04X+%d is past the end of %s (%d bytes)", prefix, start, n, name, len(b))
	}
	for off := 0; off < n; off += 16 {
		end := off + 16
		if end > n {
			end = n
		}
		fmt.Fprintf(out, "%s%04X  %s\n", prefix, start+off, hexBytes(b[start+off:start+end]))
	}
	return nil
}

// stateDiff reports where two states differ, grouping nearby bytes into
// runs so a single setting change reads as one line rather than many.
func stateDiff(pathA, pathB, name string, out io.Writer) error {
	stA, err := mesenstate.Open(pathA)
	if err != nil {
		return err
	}
	stB, err := mesenstate.Open(pathB)
	if err != nil {
		return err
	}
	a, prefix, err := region(stA, name)
	if err != nil {
		return err
	}
	b, _, err := region(stB, name)
	if err != nil {
		return err
	}
	if len(a) != len(b) {
		return fmt.Errorf("%s images differ in size: %d vs %d", name, len(a), len(b))
	}
	var diffs []int
	for i := range a {
		if a[i] != b[i] {
			diffs = append(diffs, i)
		}
	}
	if len(diffs) == 0 {
		fmt.Fprintf(out, "%s: identical (%d bytes)\n", name, len(a))
		return nil
	}
	// Join runs separated by up to three matching bytes; a settings
	// field spread over adjacent bytes stays one line.
	const gap = 4
	runs := 0
	start, prev := diffs[0], diffs[0]
	flush := func(lo, hi int) {
		runs++
		fmt.Fprintf(out, "%s%04X-%s%04X  A=%s  B=%s\n",
			prefix, lo, prefix, hi, hexBytes(a[lo:hi+1]), hexBytes(b[lo:hi+1]))
	}
	for _, i := range diffs[1:] {
		if i-prev < gap {
			prev = i
			continue
		}
		flush(start, prev)
		start, prev = i, i
	}
	flush(start, prev)
	fmt.Fprintf(out, "%s: %d differing bytes in %d runs\n", name, len(diffs), runs)
	return nil
}

func parseHex(s string) (int, error) {
	s = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(s, "+"), "$"), "0x")
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("want hexadecimal: %w", err)
	}
	return int(v), nil
}

func hexBytes(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = fmt.Sprintf("%02x", v)
	}
	return strings.Join(parts, " ")
}
