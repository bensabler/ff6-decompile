package extract

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/rom"
)

// ManifestPath is the tracked asset manifest, relative to the repo root.
const ManifestPath = "manifests/assets.json"

// Manifest is the tracked public record of every locally generated asset.
type Manifest struct {
	SchemaVersion string  `json:"schema_version"`
	ROMRevision   string  `json:"rom_revision"`
	ArchiveRoot   string  `json:"archive_root"`
	Note          string  `json:"note"`
	Assets        []Asset `json:"assets"`
}

const manifestNote = "Locally generated archival assets. The files themselves are never " +
	"committed (see docs/legal/ASSET_POLICY.md); this manifest records their identity, " +
	"provenance and hashes so a fresh clone with the verified ROM can regenerate them " +
	"byte-for-byte via `ff6lab extract all` and prove it with `ff6lab archive verify`."

// LoadManifest reads the tracked manifest.
func LoadManifest(repoRoot string) (*Manifest, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, ManifestPath))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", ManifestPath, err)
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", ManifestPath, err)
	}
	return &m, nil
}

// WriteManifest writes assets to the tracked manifest in stable ID order.
func WriteManifest(repoRoot, romRevision string, assets []Asset) error {
	sorted := append([]Asset(nil), assets...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	b, err := jsonBytes(Manifest{
		SchemaVersion: "1.0",
		ROMRevision:   romRevision,
		ArchiveRoot:   ArchiveRoot,
		Note:          manifestNote,
		Assets:        sorted,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(repoRoot, ManifestPath), b, 0o644)
}

// VerifyReport is the outcome of `archive verify`.
type VerifyReport struct {
	OK      []string // regenerates identically and matches the file on disk
	Missing []string // in the manifest, absent from the local archive
	Changed []string // on-disk bytes differ from the manifest hash
	Drifted []string // current extractor output differs from the manifest hash
	Unknown []string // present in the archive, claimed by no manifest entry
}

// Problems reports whether anything needs operator attention.
func (v *VerifyReport) Problems() int {
	return len(v.Missing) + len(v.Changed) + len(v.Drifted) + len(v.Unknown)
}

// Verify regenerates every asset in memory and compares three ways: the
// extractor's current output, the tracked manifest hash, and the bytes on
// disk. It also reports archive files no manifest entry claims.
//
// Nothing is written; verification is read-only by construction.
func Verify(r *rom.ROM, repoRoot string) (*VerifyReport, error) {
	man, err := LoadManifest(repoRoot)
	if err != nil {
		return nil, err
	}
	if man.ROMRevision != "" && man.ROMRevision != r.SHA256() {
		return nil, fmt.Errorf("manifest records ROM %s but the supplied ROM is %s", man.ROMRevision, r.SHA256())
	}

	// Regenerate every extractor's output in memory.
	generated := map[string][]byte{}
	for _, ex := range registry() {
		outs, err := ex.Extract(r)
		if err != nil {
			return nil, fmt.Errorf("extractor %s: %w", ex.ID(), err)
		}
		for _, o := range outs {
			generated[o.asset.OutputPath] = o.data
		}
	}

	rep := &VerifyReport{}
	claimed := map[string]bool{}
	for _, a := range man.Assets {
		claimed[filepath.Clean(a.OutputPath)] = true

		gen, ok := generated[a.OutputPath]
		if !ok {
			rep.Drifted = append(rep.Drifted, fmt.Sprintf("%s (%s): no extractor produces %s any more", a.ID, a.Name, a.OutputPath))
			continue
		}
		if h := hashBytes(gen); h != a.SHA256 {
			rep.Drifted = append(rep.Drifted, fmt.Sprintf("%s (%s): extractor output %s, manifest %s", a.ID, a.Name, h[:12], a.SHA256[:12]))
			continue
		}
		onDisk, err := os.ReadFile(filepath.Join(repoRoot, a.OutputPath))
		if err != nil {
			rep.Missing = append(rep.Missing, fmt.Sprintf("%s (%s): %s not generated locally", a.ID, a.Name, a.OutputPath))
			continue
		}
		if h := hashBytes(onDisk); h != a.SHA256 {
			rep.Changed = append(rep.Changed, fmt.Sprintf("%s (%s): on-disk %s, manifest %s", a.ID, a.Name, h[:12], a.SHA256[:12]))
			continue
		}
		rep.OK = append(rep.OK, fmt.Sprintf("%s (%s)", a.ID, a.Name))
	}

	// Archive files nothing claims.
	root := filepath.Join(repoRoot, ArchiveRoot)
	_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil //nolint:nilerr // a missing archive is reported via Missing
		}
		if d.Name() == ".gitkeep" || d.Name() == "README.md" {
			return nil
		}
		rel, rerr := filepath.Rel(repoRoot, p)
		if rerr != nil {
			return nil
		}
		rel = filepath.ToSlash(filepath.Clean(rel))
		if !claimed[rel] {
			rep.Unknown = append(rep.Unknown, rel)
		}
		return nil
	})

	sort.Strings(rep.OK)
	sort.Strings(rep.Missing)
	sort.Strings(rep.Changed)
	sort.Strings(rep.Drifted)
	sort.Strings(rep.Unknown)
	return rep, nil
}

// Inventory summarizes the tracked manifest without needing a ROM.
func Inventory(repoRoot string) (string, error) {
	man, err := LoadManifest(repoRoot)
	if err != nil {
		return "", err
	}
	byCat := map[string][]Asset{}
	for _, a := range man.Assets {
		byCat[a.Category] = append(byCat[a.Category], a)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "archive inventory (manifest %s)\n", ManifestPath)
	fmt.Fprintf(&b, "  rom revision: %s\n", man.ROMRevision)
	fmt.Fprintf(&b, "  archive root: %s (ignored; regenerate with `ff6lab extract all`)\n", man.ArchiveRoot)
	fmt.Fprintf(&b, "  assets: %d\n", len(man.Assets))
	for _, cat := range Categories {
		as := byCat[cat]
		if len(as) == 0 {
			note := CategoryNotes[cat]
			if note == "" {
				note = "no extractors registered"
			}
			fmt.Fprintf(&b, "\n  %s: none - %s\n", cat, note)
			continue
		}
		fmt.Fprintf(&b, "\n  %s (%d):\n", cat, len(as))
		for _, a := range as {
			fmt.Fprintf(&b, "    %s  %-46s %s\n", a.ID, a.Name, a.OutputPath)
			fmt.Fprintf(&b, "      %s | %s\n", a.ROMSource, a.OutputFormat)
			if a.Properties != "" {
				fmt.Fprintf(&b, "      %s\n", a.Properties)
			}
			fmt.Fprintf(&b, "      sha256 %s\n", a.SHA256)
		}
	}
	return b.String(), nil
}
