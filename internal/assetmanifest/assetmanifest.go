// Package assetmanifest is the data model for the tracked asset manifest.
//
// It is deliberately separate from internal/extract and imports nothing that
// touches a ROM. Extraction *writes* the manifest and needs a ROM to do it;
// the demo only *reads* it, and must not link ROM-reading code at all
// (docs/legal/ASSET_POLICY.md). Splitting the model from the extractor is
// what makes that boundary real rather than aspirational — `ff6lab audit`
// checks the import, and without this split internal/content would pull in
// internal/rom transitively.
package assetmanifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Path is the tracked asset manifest, relative to the repository root.
const Path = "manifests/assets.json"

// ArchiveRoot is where generated assets live. The directory is gitignored;
// the manifest is the tracked record of what belongs in it.
const ArchiveRoot = "local_artifacts/archive"

// Asset is one locally generated file, with the provenance needed to
// regenerate and verify it.
type Asset struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	ROMRevision     string   `json:"rom_revision"`
	ROMSource       string   `json:"rom_source"`
	ExtractorID     string   `json:"extractor_id"`
	ExtractorVer    string   `json:"extractor_version"`
	OutputPath      string   `json:"output_path"`
	OutputFormat    string   `json:"output_format"`
	Properties      string   `json:"properties,omitempty"`
	SHA256          string   `json:"sha256"`
	RuntimeConsumer string   `json:"runtime_consumer,omitempty"`
	CensusRefs      []string `json:"census_refs,omitempty"`
	ExperimentRefs  []string `json:"experiment_refs,omitempty"`
	Verification    string   `json:"verification_status"`
}

// Manifest is the tracked public record of every locally generated asset.
type Manifest struct {
	SchemaVersion string  `json:"schema_version"`
	ROMRevision   string  `json:"rom_revision"`
	ArchiveRoot   string  `json:"archive_root"`
	Note          string  `json:"note"`
	Assets        []Asset `json:"assets"`
}

// Load reads the tracked manifest from a repository root.
func Load(repoRoot string) (*Manifest, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, Path))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", Path, err)
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", Path, err)
	}
	return &m, nil
}
