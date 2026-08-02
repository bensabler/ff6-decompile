package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/platform/snesaddr"
)

// addrCmd translates between the two address spaces the project's records use.
// It reads no ROM: the mapping is arithmetic, and being able to check an
// address without a ROM in hand is the point.
func addrCmd(args []string, w io.Writer) error {
	const usage = "usage: ff6lab addr cpu <ROMCPU addr> | ff6lab addr file <ROMFILE offset>"
	if len(args) < 2 {
		return fmt.Errorf("%s", usage)
	}
	n, err := parseHexish(args[1])
	if err != nil {
		return err
	}

	switch args[0] {
	case "cpu":
		if n > 0xFFFFFF {
			return fmt.Errorf("ROMCPU:$%X exceeds the 24-bit CPU address space", n)
		}
		off, win, err := snesaddr.ROMFile(uint32(n))
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "ROMCPU:$%06X  ->  ROMFILE:0x%06X\n", n, off)
		fmt.Fprintf(w, "  window:     %s\n", win)
		fmt.Fprintf(w, "  confidence: %s\n", confidenceOf(win))
		if canon, err := snesaddr.ROMCPU(off); err == nil && uint32(n) != canon {
			fmt.Fprintf(w, "  canonical:  ROMCPU:$%06X (same bytes, Confirmed window)\n", canon)
		}
		return nil

	case "file":
		cpu, err := snesaddr.ROMCPU(int(n))
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "ROMFILE:0x%06X  ->  ROMCPU:$%06X\n", n, cpu)
		fmt.Fprintf(w, "  window:     %s\n", snesaddr.WindowHiROMUpper)
		fmt.Fprintf(w, "  confidence: %s\n", confidenceOf(snesaddr.WindowHiROMUpper))
		return nil

	default:
		return fmt.Errorf("%s", usage)
	}
}

func confidenceOf(w snesaddr.Window) string {
	if w.Confirmed() {
		return "Confirmed (18/18 Mesen ROM captures, CORR-0001)"
	}
	return "Strong hypothesis (standard HiROM; unobserved in FF6 evidence)"
}

// parseHexish accepts the notations the project's records actually use:
// $C46AC0, 0xC46AC0, and a bare hex value.
func parseHexish(s string) (uint64, error) {
	t := strings.TrimSpace(s)
	t = strings.TrimPrefix(t, "ROMCPU:")
	t = strings.TrimPrefix(t, "ROMFILE:")
	t = strings.TrimPrefix(t, "$")
	t = strings.TrimPrefix(t, "0x")
	t = strings.TrimPrefix(t, "0X")
	t = strings.ReplaceAll(t, ":", "") // bank-separated form, $C4:6AC0
	n, err := strconv.ParseUint(t, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing address %q: want hex, e.g. $C46AC0 or 0x046AC0", s)
	}
	return n, nil
}
