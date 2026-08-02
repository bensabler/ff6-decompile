// Command ff6demo is the playable Final Fantasy VI reconstruction
// (DEMO-0001).
//
// It loads every asset from the locally generated archive under
// local_artifacts/archive/, which the operator produces from their own
// verified ROM:
//
//	export FF6_ROM=/path/to/your/verified.sfc
//	go run ./cmd/ff6lab extract all
//
// No ROM-derived bytes ship in this binary, and it cannot read a ROM: the
// import boundary that guarantees that is checked by `ff6lab audit`.
//
// # Headless mode
//
// `-headless` runs without a display and optionally writes PNG frames. It is
// the authoritative mode for automated comparison, because a windowed host
// may skip Draw or pause Update and so cannot be frame-exact.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/engine/headless"
	"github.com/bensabler/ff6-decompile/internal/game/scenes"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// defaultOutDir keeps rendered frames inside the ignored artifact tree.
// Rendered frames are ROM-derived images, so writing them anywhere else needs
// an explicit opt-in.
const defaultOutDir = "local_artifacts/demo-frames"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ff6demo:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("ff6demo", flag.ContinueOnError)
	var (
		headlessMode = fs.Bool("headless", false, "run without a display")
		frames       = fs.Int("frames", 120, "headless: frames to simulate")
		captureEvery = fs.Int("capture-every", 0, "headless: write a PNG every n frames (0 = none)")
		captureLast  = fs.Bool("capture-last", false, "headless: write a PNG of the final frame")
		outDir       = fs.String("out", defaultOutDir, "headless: directory for captured frames")
		forceOut     = fs.Bool("force-out", false, "allow -out outside local_artifacts/")
		archiveRoot  = fs.String("archive", "", "repository root holding the generated archive")
		script       = fs.String("input", "", "headless: scripted input, e.g. 'A@10,Right@20-40'")
		printHashes  = fs.Bool("hashes", false, "headless: print per-frame framebuffer hashes")
		formationID  = fs.Int("formation", 44, "formation id the battle scene shows (A opens it)")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	archive, err := content.Open(*archiveRoot)
	if err != nil {
		return err
	}
	font, err := content.LoadHUDFont(archive)
	if err != nil {
		return err
	}

	src := engine.NoInput
	if *script != "" {
		src, err = parseInputScript(*script)
		if err != nil {
			return err
		}
	}

	// The battle tables are optional: a demo that cannot open a battle is
	// still a demo, and a missing table should not stop it launching.
	boot := scenes.NewBoot(font)
	if tables, err := content.LoadBattleTables(archive); err == nil {
		boot = boot.WithBattle(tables, *formationID)
	} else {
		fmt.Fprintln(os.Stderr, "ff6demo: battle scene unavailable:", err)
	}

	m := engine.New(src, nil, boot)

	if !*headlessMode {
		// A live keyboard is not reproducible, so the windowed host supplies
		// its own Source unless a script was given explicitly.
		if *script == "" {
			m = engine.New(windowedInput(), nil, boot)
		}
		return runWindowed(m)
	}

	capture := *captureEvery > 0 || *captureLast
	dir := ""
	if capture {
		if dir, err = safeOutDir(*outDir, *forceOut); err != nil {
			return err
		}
	}
	res, err := headless.Run(m, headless.Options{
		Frames:       *frames,
		CaptureEvery: *captureEvery,
		CaptureLast:  *captureLast,
		OutDir:       dir,
		Progress:     os.Stdout,
	})
	if err != nil {
		return err
	}
	fmt.Printf("simulated %d frames", res.Frames)
	if res.StoppedEarly {
		fmt.Print(" (exited early)")
	}
	fmt.Println()
	if *printHashes {
		for i, h := range res.Hashes {
			fmt.Printf("frame %6d %x\n", i+1, h)
		}
	}
	return nil
}

// safeOutDir refuses to scatter ROM-derived frames across the filesystem.
// The same rule the extractor applies to the archive applies here.
func safeOutDir(dir string, force bool) (string, error) {
	clean := filepath.Clean(dir)
	if force {
		return clean, nil
	}
	slash := filepath.ToSlash(clean)
	if slash != "local_artifacts" && !strings.HasPrefix(slash, "local_artifacts/") {
		return "", fmt.Errorf("refusing to write rendered frames to %s: they are ROM-derived images and belong under local_artifacts/ (pass -force-out to override)", dir)
	}
	return clean, nil
}

// parseInputScript builds a deterministic input Source from a compact
// notation, so a scenario run is reproducible from a command line:
//
//	A@10          press A on frame 10 only
//	Right@20-40   hold Right from frame 20 through 40
//	A@10,Start@90 comma-separated
func parseInputScript(s string) (engine.Source, error) {
	type span struct {
		btn      snespad.Button
		from, to uint64
	}
	byName := map[string]snespad.Button{}
	for _, b := range snespad.All {
		byName[strings.ToLower(b.String())] = b
	}

	var spans []span
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		name, rng, ok := strings.Cut(part, "@")
		if !ok {
			return nil, fmt.Errorf("input script %q: each entry needs Button@frame or Button@from-to", part)
		}
		btn, ok := byName[strings.ToLower(strings.TrimSpace(name))]
		if !ok {
			return nil, fmt.Errorf("input script %q: unknown button %q", part, name)
		}
		var from, to uint64
		if lo, hi, isRange := strings.Cut(rng, "-"); isRange {
			if _, err := fmt.Sscanf(lo, "%d", &from); err != nil {
				return nil, fmt.Errorf("input script %q: bad start frame", part)
			}
			if _, err := fmt.Sscanf(hi, "%d", &to); err != nil {
				return nil, fmt.Errorf("input script %q: bad end frame", part)
			}
		} else {
			if _, err := fmt.Sscanf(rng, "%d", &from); err != nil {
				return nil, fmt.Errorf("input script %q: bad frame", part)
			}
			to = from
		}
		if to < from {
			return nil, fmt.Errorf("input script %q: end frame %d precedes start %d", part, to, from)
		}
		spans = append(spans, span{btn, from, to})
	}

	return engine.SourceFunc(func(frame uint64) snespad.State {
		var st snespad.State
		for _, sp := range spans {
			if frame >= sp.from && frame <= sp.to {
				st = st.With(sp.btn)
			}
		}
		return st
	}), nil
}
