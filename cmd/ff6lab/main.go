package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/audio/brr"
	"github.com/bensabler/ff6-decompile/internal/audit"
	"github.com/bensabler/ff6-decompile/internal/census"
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
			return censusSync(".")
		}
		return fmt.Errorf("usage: ff6lab indexes generate")
	case "census":
		if len(args) > 1 && args[1] == "validate" {
			return censusValidate(".")
		}
		if len(args) > 1 && args[1] == "sync" {
			if err := censusSync("."); err != nil {
				return err
			}
			return censusValidate(".")
		}
		return fmt.Errorf("usage: ff6lab census validate|sync")
	case "coverage":
		return coverageCmd(args[1:])
	case "rom":
		if len(args) > 1 && args[1] == "gaps" {
			return romGaps(".")
		}
		return fmt.Errorf("usage: ff6lab rom gaps")
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
	case "brr":
		if len(args) > 2 && args[1] == "info" {
			return brrInfo(args[2], os.Stdout)
		}
		return fmt.Errorf("usage: ff6lab brr info <hexdump-file>")
	case "extract":
		return runExtract(args[1:], os.Stdout)
	case "archive":
		return runArchive(args[1:], os.Stdout)
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

// censusSync regenerates every census-derived file.
func censusSync(root string) error {
	c, r, err := census.Load(root)
	if err != nil {
		return err
	}
	for name, content := range map[string]string{
		"indexes/CONTENT_CENSUS.md": census.GenerateCensusIndex(c),
		"indexes/ROM_REGIONS.md":    census.GenerateRegionIndex(r),
		"dashboards/COVERAGE.md":    census.GenerateCoverageDashboard(c, r),
	} {
		if err := os.WriteFile(root+"/"+name, []byte(content), 0o644); err != nil {
			return err
		}
		fmt.Println(name, "regenerated")
	}
	return nil
}

func censusValidate(root string) error {
	issues, err := census.Validate(root)
	if err != nil {
		return err
	}
	for _, i := range issues {
		fmt.Println(i)
	}
	if len(issues) > 0 {
		return fmt.Errorf("census: %d issue(s)", len(issues))
	}
	fmt.Println("census: clean")
	return nil
}

func coverageCmd(args []string) error {
	c, _, err := census.Load(".")
	if err != nil {
		return err
	}
	sub := "summary"
	if len(args) > 0 {
		sub = args[0]
	}
	switch sub {
	case "summary":
		fmt.Println("Domain      Registered  Located  Decoded  Extracted  Implemented  RuntimeVerified")
		for _, d := range census.Summary(c) {
			fmt.Printf("%-11s %10d %8d %8d %10d %12d %16d\n",
				d.Domain, d.Registered, d.Located, d.Decoded, d.Extracted, d.Implemented, d.RuntimeVerified)
		}
		return nil
	case "gaps":
		for _, e := range c.Entries {
			if !census.AtLeast(census.ReconstructionLadder, e.ReconstructionStatus, "LOCATED") {
				fmt.Printf("%s [%s/%s] %s — next: %s\n",
					e.ID, e.ReconstructionStatus, e.RuntimeStatus, e.Name, e.NextAction)
			}
		}
		return nil
	case "domain":
		if len(args) < 2 {
			return fmt.Errorf("usage: ff6lab coverage domain <name>")
		}
		for _, e := range c.Entries {
			if !strings.EqualFold(e.Domain, args[1]) {
				continue
			}
			fmt.Printf("%s [%s/%s] %s\n  next: %s\n", e.ID, e.ReconstructionStatus, e.RuntimeStatus, e.Name, e.NextAction)
		}
		return nil
	}
	return fmt.Errorf("usage: ff6lab coverage summary|gaps|domain <name>")
}

func romGaps(root string) error {
	_, r, err := census.Load(root)
	if err != nil {
		return err
	}
	st := census.Analyze(r)
	known := st.KnownBytes
	fmt.Printf("ROM %d bytes; known %d (%.2f%%); candidate-only %d; %d gaps\n",
		st.RomSize, known, float64(known)*100/float64(st.RomSize), st.CandidateBytes, len(st.Gaps))
	fmt.Println("largest unknown spans:")
	for _, g := range census.TopGaps(st, 10) {
		fmt.Printf("  0x%06X-0x%06X (%d bytes)\n", g.Start, g.Start+g.Size-1, g.Size)
	}
	return nil
}

// brrInfo decodes a BRR stream from a hex dump and prints its block
// headers and PCM shape (count, min/max, leading samples).
func brrInfo(path string, w *os.File) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	raw, err := parseHexBytes(string(data))
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	pcm, err := brr.Decode(raw)
	if err != nil {
		return err
	}
	blocks := len(pcm) / 16
	for i := 0; i < blocks; i++ {
		h, err := brr.ParseHeader(raw[i*brr.BlockSize])
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "block %d: range=%d filter=%d loop=%t end=%t\n", i, h.Range, h.Filter, h.Loop, h.End)
	}
	lo, hi := pcm[0], pcm[0]
	for _, s := range pcm {
		if s < lo {
			lo = s
		}
		if s > hi {
			hi = s
		}
	}
	fmt.Fprintf(w, "samples=%d min=%d max=%d\n", len(pcm), lo, hi)
	n := len(pcm)
	if n > 8 {
		n = 8
	}
	fmt.Fprintf(w, "pcm[0:%d]=%v\n", n, pcm[:n])
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
