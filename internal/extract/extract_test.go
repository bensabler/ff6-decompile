package extract

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/rom"
)

// loadROM returns the verified ROM, or skips: the ROM is user-supplied and
// never committed, so ROM-dependent tests are opt-in via FF6_ROM.
func loadROM(t *testing.T) *rom.ROM {
	t.Helper()
	p := os.Getenv(rom.EnvVar)
	if p == "" {
		t.Skipf("set %s to run ROM-dependent extraction tests", rom.EnvVar)
	}
	r, err := rom.Load(p)
	if err != nil {
		t.Fatalf("loading ROM: %v", err)
	}
	return r
}

func TestRegistryIsWellFormed(t *testing.T) {
	valid := map[string]bool{}
	for _, c := range Categories {
		valid[c] = true
	}
	ids := map[string]bool{}
	for _, ex := range registry() {
		if !valid[ex.Category()] {
			t.Errorf("extractor %s has category %q, not in Categories", ex.ID(), ex.Category())
		}
		if ids[ex.ID()] {
			t.Errorf("duplicate extractor ID %q", ex.ID())
		}
		ids[ex.ID()] = true
		if ex.Version() == "" {
			t.Errorf("extractor %s has no version", ex.ID())
		}
	}
}

// TestExtractionIsDeterministic is the property the whole architecture
// rests on: a fresh clone plus the correct ROM must regenerate identical
// bytes, or the tracked hashes are meaningless.
func TestExtractionIsDeterministic(t *testing.T) {
	r := loadROM(t)
	first := map[string]string{}
	for _, ex := range registry() {
		outs, err := ex.Extract(r)
		if err != nil {
			t.Fatalf("extractor %s: %v", ex.ID(), err)
		}
		for _, o := range outs {
			first[o.asset.OutputPath] = hashBytes(o.data)
		}
	}
	if len(first) == 0 {
		t.Fatal("no outputs produced")
	}
	for i := 0; i < 3; i++ {
		for _, ex := range registry() {
			outs, err := ex.Extract(r)
			if err != nil {
				t.Fatalf("extractor %s (pass %d): %v", ex.ID(), i, err)
			}
			for _, o := range outs {
				if got := hashBytes(o.data); got != first[o.asset.OutputPath] {
					t.Fatalf("%s not deterministic: pass %d hash %s, first %s", o.asset.OutputPath, i, got, first[o.asset.OutputPath])
				}
			}
		}
	}
}

// TestAssetsMatchTrackedManifest proves the committed manifest still
// describes what the extractors produce.
func TestAssetsMatchTrackedManifest(t *testing.T) {
	r := loadROM(t)
	root := repoRoot(t)
	man, err := LoadManifest(root)
	if err != nil {
		t.Fatalf("loading manifest: %v", err)
	}
	if man.ROMRevision != r.SHA256() {
		t.Fatalf("manifest ROM %s, supplied ROM %s", man.ROMRevision, r.SHA256())
	}
	want := map[string]string{}
	for _, a := range man.Assets {
		want[a.OutputPath] = a.SHA256
	}
	for _, ex := range registry() {
		outs, err := ex.Extract(r)
		if err != nil {
			t.Fatalf("extractor %s: %v", ex.ID(), err)
		}
		for _, o := range outs {
			w, ok := want[o.asset.OutputPath]
			if !ok {
				t.Errorf("%s produces %s, absent from the manifest", ex.ID(), o.asset.OutputPath)
				continue
			}
			if got := hashBytes(o.data); got != w {
				t.Errorf("%s: extractor hash %s, manifest %s", o.asset.OutputPath, got, w)
			}
			delete(want, o.asset.OutputPath)
		}
	}
	for p := range want {
		t.Errorf("manifest claims %s but no extractor produces it", p)
	}
}

// TestManifestEntriesAreComplete checks every tracked entry carries the
// provenance the architecture promises — without needing a ROM.
func TestManifestEntriesAreComplete(t *testing.T) {
	man, err := LoadManifest(repoRoot(t))
	if err != nil {
		t.Fatalf("loading manifest: %v", err)
	}
	if len(man.Assets) == 0 {
		t.Fatal("manifest has no assets")
	}
	if man.ArchiveRoot != ArchiveRoot {
		t.Errorf("manifest archive root %q, want %q", man.ArchiveRoot, ArchiveRoot)
	}
	valid := map[string]bool{}
	for _, c := range Categories {
		valid[c] = true
	}
	seen := map[string]bool{}
	for _, a := range man.Assets {
		switch {
		case a.ID == "":
			t.Error("asset with no ID")
		case a.Name == "":
			t.Errorf("%s: no name", a.ID)
		case !valid[a.Category]:
			t.Errorf("%s: category %q not in Categories", a.ID, a.Category)
		case a.ROMRevision != rom.SupportedSHA256:
			t.Errorf("%s: ROM revision %q is not the supported image", a.ID, a.ROMRevision)
		case a.ROMSource == "":
			t.Errorf("%s: no ROM source range", a.ID)
		case a.ExtractorID == "" || a.ExtractorVer == "":
			t.Errorf("%s: incomplete extractor identity", a.ID)
		case len(a.SHA256) != 64:
			t.Errorf("%s: sha256 %q is not a full digest", a.ID, a.SHA256)
		case a.OutputFormat == "":
			t.Errorf("%s: no output format", a.ID)
		case a.Verification == "":
			t.Errorf("%s: no verification status", a.ID)
		}
		if seen[a.ID] {
			t.Errorf("duplicate asset ID %s", a.ID)
		}
		seen[a.ID] = true
		// Outputs must live in the ignored archive: the manifest is the
		// public record, the bytes never are.
		if got := filepath.ToSlash(filepath.Clean(a.OutputPath)); !hasPrefix(got, ArchiveRoot+"/") {
			t.Errorf("%s: output %s is outside %s", a.ID, a.OutputPath, ArchiveRoot)
		}
	}
}

func hasPrefix(s, p string) bool { return len(s) >= len(p) && s[:len(p)] == p }

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not locate module root")
	return ""
}
