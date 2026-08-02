// Package content loads game data from the locally generated archive and
// adapts it to engine-ready forms.
//
// # The ROM boundary
//
// This package, and everything above it, must never read the ROM. The demo
// consumes only `local_artifacts/archive/`, which the operator regenerates
// from their own verified image with `ff6lab extract all`. That is the asset
// policy expressed as a code boundary, and `ff6lab audit` enforces it by
// rejecting an `internal/rom` import from here, from `internal/engine`, or
// from `cmd/ff6demo`.
//
// # Archive formats
//
// The on-disk format of each asset is chosen by its extractor, for archival
// fidelity and human inspectability — indexed PNG for graphics, JSON for
// structured tables, raw bytes for undecoded material. This package adapts
// those to in-memory forms the engine wants. It never re-derives an asset
// from the ROM at runtime, and never invents one that is missing.
//
// Every read is verified against the SHA-256 the manifest records, so a
// stale or hand-edited archive fails loudly instead of rendering wrong data.
package content

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bensabler/ff6-decompile/internal/assetmanifest"
)

// Sentinel errors, so a host can tell "you have not generated the archive"
// from "your archive is out of date" and print the right instruction.
var (
	// ErrNoArchive means no generated archive was found.
	ErrNoArchive = errors.New("no generated asset archive")
	// ErrAssetMissing means the manifest describes an asset the archive
	// does not contain.
	ErrAssetMissing = errors.New("asset missing from the archive")
	// ErrAssetStale means an archive file's contents no longer match the
	// hash the manifest records.
	ErrAssetStale = errors.New("archive asset does not match the tracked manifest")
)

// EnvVar names the environment variable that overrides archive discovery.
const EnvVar = "FF6_ARCHIVE"

// SetupHint is the operator-facing instruction printed whenever the archive
// is absent or stale. It names the exact commands, because "generate the
// assets" is not actionable.
const SetupHint = `The demo renders data extracted from your own verified ROM.
No ROM-derived bytes ship in this binary.

  export FF6_ROM=/path/to/your/verified.sfc
  go run ./cmd/ff6lab extract all

Then run the demo again.`

// Archive is a verified handle on a generated asset archive.
type Archive struct {
	root     string
	manifest *assetmanifest.Manifest
	byID     map[string]assetmanifest.Asset
}

// Open locates and validates an archive. Discovery order: an explicit root,
// then $FF6_ARCHIVE, then the repository containing the working directory.
func Open(explicit string) (*Archive, error) {
	root, err := resolveRoot(explicit)
	if err != nil {
		return nil, err
	}
	man, err := assetmanifest.Load(root)
	if err != nil {
		return nil, fmt.Errorf("%w: reading the manifest under %s: %v\n\n%s", ErrNoArchive, root, err, SetupHint)
	}
	if len(man.Assets) == 0 {
		return nil, fmt.Errorf("%w: the manifest under %s lists no assets\n\n%s", ErrNoArchive, root, SetupHint)
	}
	a := &Archive{root: root, manifest: man, byID: make(map[string]assetmanifest.Asset, len(man.Assets))}
	for _, as := range man.Assets {
		a.byID[as.ID] = as
	}
	// The archive directory itself must exist, or every read would fail one
	// at a time with a confusing message.
	if _, err := os.Stat(filepath.Join(root, man.ArchiveRoot)); err != nil {
		return nil, fmt.Errorf("%w: %s is not present\n\n%s", ErrNoArchive, man.ArchiveRoot, SetupHint)
	}
	return a, nil
}

// Root returns the repository root the archive was found under.
func (a *Archive) Root() string { return a.root }

// ROMRevision returns the ROM the archive was generated from.
func (a *Archive) ROMRevision() string { return a.manifest.ROMRevision }

// Asset returns the manifest entry for an asset ID.
func (a *Archive) Asset(id string) (assetmanifest.Asset, error) {
	as, ok := a.byID[id]
	if !ok {
		return assetmanifest.Asset{}, fmt.Errorf("%w: no manifest entry for %s", ErrAssetMissing, id)
	}
	return as, nil
}

// Read returns an asset's bytes, verified against the manifest hash.
func (a *Archive) Read(id string) ([]byte, error) {
	as, err := a.Asset(id)
	if err != nil {
		return nil, err
	}
	p := filepath.Join(a.root, filepath.FromSlash(as.OutputPath))
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s (%s)\n\n%s", ErrAssetMissing, as.OutputPath, id, SetupHint)
		}
		return nil, fmt.Errorf("reading %s: %w", as.OutputPath, err)
	}
	sum := sha256.Sum256(data)
	if got := hex.EncodeToString(sum[:]); got != as.SHA256 {
		return nil, fmt.Errorf("%w: %s has sha256 %s, manifest records %s; regenerate with `ff6lab extract all -force`",
			ErrAssetStale, as.OutputPath, got, as.SHA256)
	}
	return data, nil
}

// resolveRoot finds the repository root that owns the archive.
func resolveRoot(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if p := os.Getenv(EnvVar); p != "" {
		return p, nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, assetmanifest.Path)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("%w: no %s found in this directory or its parents; pass -archive or set %s\n\n%s",
		ErrNoArchive, assetmanifest.Path, EnvVar, SetupHint)
}
