// Package extract implements the project's extractor-first archival
// architecture: deterministic generators that read the user-supplied
// verified ROM and materialize the substantial-asset archive under an
// ignored local path, while the public repository tracks only manifests
// (identity, provenance, hashes) plus structured reconstruction data
// permitted by docs/legal/ASSET_POLICY.md.
//
// Every extractor is registered here with a category and a version; a
// fresh clone plus the correct ROM regenerates the identical archive, and
// `ff6lab archive verify` proves it by recomputing every output in memory
// and comparing hashes against the tracked manifest and the on-disk files.
package extract

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/assetmanifest"
	"github.com/bensabler/ff6-decompile/internal/rom"
)

// ArchiveRoot is the ignored local root the archive is generated under,
// relative to the repository root.
const ArchiveRoot = assetmanifest.ArchiveRoot

// Categories in fixed presentation order. Categories may be registered
// with no extractors yet; running them reports that honestly instead of
// fabricating output.
var Categories = []string{"data", "graphics", "audio", "dialogue", "maps", "raw"}

// CategoryNotes explains categories that currently have no extractors.
var CategoryNotes = map[string]string{
	"dialogue": "no confirmed source located yet (dialogue/text table unlocated - CEN-EVENT-0002)",
	"maps":     "no confirmed source located yet (map headers/tilesets unlocated - CEN-WORLD-0004)",
}

// Asset is one generated file's public manifest entry. Field order is the
// tracked JSON order; do not reorder casually.
// Asset aliases the manifest model. The type lives in
// internal/assetmanifest so the demo can read the manifest without linking
// ROM-reading code; extraction keeps using it under the familiar name.
type Asset = assetmanifest.Asset

// output is one file an extractor produces.
type output struct {
	asset Asset
	data  []byte
}

// Extractor generates a deterministic set of outputs from the ROM.
type Extractor interface {
	// ID is the stable extractor identifier.
	ID() string
	// Category is one of Categories.
	Category() string
	// Version participates in every produced Asset; bump it whenever the
	// output bytes could change.
	Version() string
	// Extract computes all outputs in memory. It must be deterministic.
	Extract(r *rom.ROM) ([]output, error)
}

// registry returns all extractors in deterministic order.
func registry() []Extractor {
	ex := []Extractor{
		spellExtractor{},
		esperExtractor{},
		hudFontExtractor{},
		sfxPackExtractor{},
		rawSliceExtractor{},
	}
	sort.Slice(ex, func(i, j int) bool { return ex[i].ID() < ex[j].ID() })
	return ex
}

func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Result summarizes one Run.
type Result struct {
	Written   []string // paths written (new content)
	Unchanged []string // paths already present with identical bytes
	Skipped   []string // categories with no extractors
	Assets    []Asset
}

// Run executes every extractor whose category is in wanted (nil/["all"]
// means every category), writing outputs under repoRoot/ArchiveRoot.
//
// Overwrite safety: a target whose existing bytes differ from the
// deterministic output is an error unless force is set — regeneration must
// never silently destroy something an operator placed or preserved. Paths
// outside the archive root (or inside preserved experiment evidence) are
// refused outright.
func Run(r *rom.ROM, repoRoot string, wanted []string, force bool) (*Result, error) {
	want := map[string]bool{}
	all := len(wanted) == 0
	for _, w := range wanted {
		if w == "all" {
			all = true
			continue
		}
		valid := false
		for _, c := range Categories {
			if w == c {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("unknown category %q (valid: all %s)", w, strings.Join(Categories, " "))
		}
		want[w] = true
	}

	res := &Result{}
	seen := map[string]bool{}
	for _, cat := range Categories {
		if !all && !want[cat] {
			continue
		}
		ran := false
		for _, ex := range registry() {
			if ex.Category() != cat {
				continue
			}
			ran = true
			outs, err := ex.Extract(r)
			if err != nil {
				return nil, fmt.Errorf("extractor %s: %w", ex.ID(), err)
			}
			for _, o := range outs {
				if seen[o.asset.OutputPath] {
					return nil, fmt.Errorf("extractor %s: duplicate output path %s", ex.ID(), o.asset.OutputPath)
				}
				seen[o.asset.OutputPath] = true
				status, err := writeArchiveFile(repoRoot, o.asset.OutputPath, o.data, force)
				if err != nil {
					return nil, fmt.Errorf("extractor %s: %w", ex.ID(), err)
				}
				switch status {
				case "written":
					res.Written = append(res.Written, o.asset.OutputPath)
				case "unchanged":
					res.Unchanged = append(res.Unchanged, o.asset.OutputPath)
				}
				res.Assets = append(res.Assets, o.asset)
			}
		}
		if !ran {
			note := CategoryNotes[cat]
			if note == "" {
				note = "no extractors registered"
			}
			res.Skipped = append(res.Skipped, cat+": "+note)
		}
	}
	sort.Slice(res.Assets, func(i, j int) bool { return res.Assets[i].ID < res.Assets[j].ID })
	return res, nil
}

// writeArchiveFile writes data to repoRoot/relPath with the safety rules
// above. Returns "written" or "unchanged".
func writeArchiveFile(repoRoot, relPath string, data []byte, force bool) (string, error) {
	clean := filepath.Clean(relPath)
	if !strings.HasPrefix(clean, ArchiveRoot+string(filepath.Separator)) {
		return "", fmt.Errorf("refusing to write outside %s: %s", ArchiveRoot, relPath)
	}
	if strings.Contains(clean, "local_artifacts/experiments") {
		return "", fmt.Errorf("refusing to touch preserved experiment evidence: %s", relPath)
	}
	abs := filepath.Join(repoRoot, clean)
	if existing, err := os.ReadFile(abs); err == nil {
		if hashBytes(existing) == hashBytes(data) {
			return "unchanged", nil
		}
		if !force {
			return "", fmt.Errorf("%s exists with different content; re-run with -force to overwrite", clean)
		}
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(abs, data, 0o644); err != nil {
		return "", err
	}
	return "written", nil
}
