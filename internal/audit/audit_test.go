package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "manifests/experiments.json",
		`{"schema_version":"1.0","experiments":[{"id":"EXP-0001","status":"completed","record":"docs/experiments/EXP-0001-x.md"}]}`)
	writeFile(t, root, "docs/experiments/EXP-0001-x.md", "# EXP-0001\n")
	writeFile(t, root, "indexes/EXPERIMENTS.md", "| EXP-0001 |\n")
	return root
}

func TestCheckManifestsClean(t *testing.T) {
	fs, err := CheckManifests(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 0 {
		t.Fatalf("want no findings, got %v", fs)
	}
}

func TestCheckManifestsBadJSONAndMissingRecord(t *testing.T) {
	root := fixture(t)
	writeFile(t, root, "manifests/broken.json", "{not json")
	writeFile(t, root, "manifests/experiments.json",
		`{"schema_version":"1.0","experiments":[{"id":"EXP-0002","status":"completed","record":"docs/experiments/missing.md"}]}`)
	fs, err := CheckManifests(root)
	if err != nil {
		t.Fatal(err)
	}
	var invalid, missing bool
	for _, f := range fs {
		if strings.Contains(f.Message, "invalid JSON") {
			invalid = true
		}
		if strings.Contains(f.Message, "not found") {
			missing = true
		}
	}
	if !invalid || !missing {
		t.Fatalf("want invalid-JSON and missing-record findings, got %v", fs)
	}
}

func TestCheckExperimentRecordsInManifest(t *testing.T) {
	root := fixture(t)
	if fs, err := CheckExperimentRecordsInManifest(root); err != nil || len(fs) != 0 {
		t.Fatalf("clean fixture: want no findings, got %v, %v", fs, err)
	}
	writeFile(t, root, "docs/experiments/EXP-0002-orphan.md", "# EXP-0002\n")
	fs, err := CheckExperimentRecordsInManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 1 || !strings.Contains(fs[0].Message, "EXP-0002-orphan.md") {
		t.Fatalf("want one orphan-record finding, got %v", fs)
	}
}

func TestCheckExperimentIndexSync(t *testing.T) {
	root := fixture(t)
	writeFile(t, root, "indexes/EXPERIMENTS.md", "# empty\n")
	fs, err := CheckExperimentIndexSync(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 1 || !strings.Contains(fs[0].Message, "EXP-0001") {
		t.Fatalf("want one missing-EXP-0001 finding, got %v", fs)
	}
}

func TestCheckMarkdownLinks(t *testing.T) {
	root := fixture(t)
	writeFile(t, root, "docs/a.md", "[ok](../manifests/experiments.json) [bad](nope.md) [web](https://x.example)\n")
	fs, err := CheckMarkdownLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 1 || !strings.Contains(fs[0].Message, "nope.md") {
		t.Fatalf("want exactly the nope.md finding, got %v", fs)
	}
}

func TestGenerateExperimentIndex(t *testing.T) {
	got, err := GenerateExperimentIndex(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "| EXP-0001 |") || !strings.Contains(got, "EXP-0001-x") {
		t.Fatalf("generated index missing entry:\n%s", got)
	}
}
