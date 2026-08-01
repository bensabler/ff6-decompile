package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/extract"
	"github.com/bensabler/ff6-decompile/internal/rom"
)

// parseROMFlag pulls an optional `-rom <path>` (and `-force`) out of args,
// returning the remaining positional arguments.
func parseROMFlag(args []string) (romPath string, force bool, rest []string, err error) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-rom", "--rom":
			if i+1 >= len(args) {
				return "", false, nil, fmt.Errorf("-rom needs a path")
			}
			romPath = args[i+1]
			i++
		case "-force", "--force":
			force = true
		default:
			rest = append(rest, args[i])
		}
	}
	return romPath, force, rest, nil
}

// loadROM resolves and verifies the user-supplied ROM. Extraction refuses
// unsupported revisions outright rather than producing wrong assets.
func loadROM(explicit string) (*rom.ROM, error) {
	p, err := rom.ResolvePath(explicit)
	if err != nil {
		return nil, err
	}
	r, err := rom.Load(p)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// runExtract implements `ff6lab extract <all|category...> [-rom p] [-force]`.
func runExtract(args []string, w io.Writer) error {
	romPath, force, rest, err := parseROMFlag(args)
	if err != nil {
		return err
	}
	if len(rest) == 0 {
		return fmt.Errorf("usage: ff6lab extract <all|%s> [-rom <path>] [-force]", strings.Join(extract.Categories, "|"))
	}
	r, err := loadROM(romPath)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "rom verified: sha256 %s\n", r.SHA256())

	res, err := extract.Run(r, ".", rest, force)
	if err != nil {
		return err
	}
	for _, p := range res.Written {
		fmt.Fprintf(w, "  wrote      %s\n", p)
	}
	for _, p := range res.Unchanged {
		fmt.Fprintf(w, "  unchanged  %s\n", p)
	}
	for _, s := range res.Skipped {
		fmt.Fprintf(w, "  skipped    %s\n", s)
	}

	// Refresh the tracked manifest only for a full run; a partial run
	// would otherwise drop the categories it did not generate.
	full := false
	for _, c := range rest {
		if c == "all" {
			full = true
		}
	}
	if full {
		if err := extract.WriteManifest(".", r.SHA256(), res.Assets); err != nil {
			return err
		}
		fmt.Fprintf(w, "manifest updated: %s (%d assets)\n", extract.ManifestPath, len(res.Assets))
	} else {
		fmt.Fprintf(w, "manifest unchanged (partial run; use `extract all` to refresh %s)\n", extract.ManifestPath)
	}
	return nil
}

// runArchive implements `ff6lab archive <verify|inventory>`.
func runArchive(args []string, w io.Writer) error {
	romPath, _, rest, err := parseROMFlag(args)
	if err != nil {
		return err
	}
	if len(rest) == 0 {
		return fmt.Errorf("usage: ff6lab archive <verify|inventory> [-rom <path>]")
	}
	switch rest[0] {
	case "inventory":
		out, err := extract.Inventory(".")
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)
		return nil
	case "verify":
		r, err := loadROM(romPath)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "rom verified: sha256 %s\n", r.SHA256())
		rep, err := extract.Verify(r, ".")
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "  ok        %d\n", len(rep.OK))
		for _, s := range rep.Missing {
			fmt.Fprintf(w, "  MISSING   %s\n", s)
		}
		for _, s := range rep.Changed {
			fmt.Fprintf(w, "  CHANGED   %s\n", s)
		}
		for _, s := range rep.Drifted {
			fmt.Fprintf(w, "  DRIFTED   %s\n", s)
		}
		for _, s := range rep.Unknown {
			fmt.Fprintf(w, "  UNKNOWN   %s\n", s)
		}
		if n := rep.Problems(); n > 0 {
			return fmt.Errorf("archive verify: %d problem(s)", n)
		}
		fmt.Fprintln(w, "archive: clean")
		return nil
	default:
		return fmt.Errorf("usage: ff6lab archive <verify|inventory> [-rom <path>]")
	}
}
