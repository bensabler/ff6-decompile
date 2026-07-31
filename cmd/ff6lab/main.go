package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/audit"
	"github.com/bensabler/ff6-decompile/internal/game/attackdata"
	"github.com/bensabler/ff6-decompile/internal/graphics/tile2bpp"
	"github.com/bensabler/ff6-decompile/internal/project"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ff6lab:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Print(project.Help())
		return nil
	}
	switch args[0] {
	case "version":
		fmt.Println(project.Version)
		return nil
	case "audit":
		findings, err := audit.Run(".")
		if err != nil {
			return err
		}
		for _, f := range findings {
			fmt.Println(f)
		}
		if len(findings) > 0 {
			return fmt.Errorf("audit: %d finding(s)", len(findings))
		}
		fmt.Println("audit: clean")
		return nil
	case "indexes":
		if len(args) > 1 && args[1] == "generate" {
			out, err := audit.GenerateExperimentIndex(".")
			if err != nil {
				return err
			}
			if err := os.WriteFile("indexes/EXPERIMENTS.md", []byte(out), 0o644); err != nil {
				return err
			}
			fmt.Println("indexes/EXPERIMENTS.md regenerated")
			return nil
		}
		return fmt.Errorf("usage: ff6lab indexes generate")
	case "attackdata":
		if len(args) > 2 && args[1] == "scan" {
			return scanAttackTable(args[2], os.Stdout)
		}
		return fmt.Errorf("usage: ff6lab attackdata scan <hexdump-file>")
	case "tiles":
		if len(args) > 3 && args[1] == "decode2bpp" {
			return decodeTile2bpp(args[2], args[3], os.Stdout)
		}
		return fmt.Errorf("usage: ff6lab tiles decode2bpp <hexdump-file> <tile-index>")
	case "help", "-h", "--help":
		fmt.Print(project.Help())
		return nil
	default:
		return fmt.Errorf("unknown command %q; run ff6lab help", args[0])
	}
}

// scanAttackTable decodes a bridge hex dump of the $C46AC0 table
// (space-separated hex bytes, as `read cpu` emits) and prints one line
// per record: index, raw bytes, and the verified accessor fields. The
// dump file stays in local_artifacts; only derived text leaves it.
func scanAttackTable(path string, w *os.File) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	table, err := parseHexBytes(string(data))
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	n := len(table) / attackdata.RecordSize
	fmt.Fprintf(w, "records=%d (from %d bytes)\n", n, len(table))
	for i := 0; i < n; i++ {
		r, err := attackdata.RecordAt(table, i)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "idx=%3d raw=% X elem=%02X flags2=%02X mode=%02X power=%3d mp=%t phys=%t\n",
			i, r[:], r.Element(), r.Flags2(), r.Mode(), r.Power(), r.TargetsMP(), r.PhysicalFormula())
	}
	return nil
}

// decodeTile2bpp prints one 8x8 tile from a hex dump as palette-index
// digits, one row per line. index counts 16-byte tiles from the dump
// start (e.g., a BG3 chr or ROM font dump).
func decodeTile2bpp(path, indexArg string, w *os.File) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	raw, err := parseHexBytes(string(data))
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	index, err := strconv.Atoi(indexArg)
	if err != nil || index < 0 {
		return fmt.Errorf("tile index %q must be a non-negative integer", indexArg)
	}
	off := index * tile2bpp.EncodedSize
	if off+tile2bpp.EncodedSize > len(raw) {
		return fmt.Errorf("tile %d out of range (dump holds %d tiles)", index, len(raw)/tile2bpp.EncodedSize)
	}
	t, err := tile2bpp.Decode(raw[off:])
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "tile %d (bytes %d-%d)\n", index, off, off+tile2bpp.EncodedSize-1)
	for _, row := range t {
		for _, px := range row {
			fmt.Fprintf(w, "%d", px)
		}
		fmt.Fprintln(w)
	}
	return nil
}

// parseHexBytes converts whitespace-separated hex byte tokens into a
// byte slice, rejecting malformed tokens with their position.
func parseHexBytes(s string) ([]byte, error) {
	fields := strings.Fields(s)
	out := make([]byte, 0, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseUint(f, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("token %d %q: %w", i, f, err)
		}
		out = append(out, byte(v))
	}
	return out, nil
}
