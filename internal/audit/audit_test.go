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

func TestCheckBattleExperimentConfig(t *testing.T) {
	const cfg = `"battle_config":{"battle_mode":"Wait","battle_speed":3,"source":"memory-read"}`
	tests := []struct {
		name  string
		entry string
		want  string // substring of the expected finding; empty means none
	}{
		{
			name:  "battle experiment with fingerprint",
			entry: `{"id":"EXP-0041","status":"completed","question":"Where is battle config stored?","starting_state":{` + cfg + `}}`,
		},
		{
			name:  "battle experiment missing fingerprint",
			entry: `{"id":"EXP-0041","status":"completed","question":"Where is battle config stored?","starting_state":{}}`,
			want:  "battle_config is missing",
		},
		{
			name:  "battle experiment detected by record path",
			entry: `{"id":"EXP-0042","status":"planned","question":"Which routine ticks?","record":"docs/experiments/EXP-0042-battle-entry.md"}`,
			want:  "battle_config is missing",
		},
		{
			name:  "ATB in question counts as battle",
			entry: `{"id":"EXP-0043","status":"planned","question":"Where do ATB gauges live?"}`,
			want:  "battle_config is missing",
		},
		{
			name:  "fingerprint without source",
			entry: `{"id":"EXP-0041","status":"completed","question":"battle config?","starting_state":{"battle_config":{"battle_mode":"Wait"}}}`,
			want:  "needs a source",
		},
		{
			name:  "explicit non-battle domain opts out",
			entry: `{"id":"EXP-0044","status":"planned","domain":"graphics","question":"Battle background provenance?"}`,
		},
		{
			name:  "explicit battle domain opts in",
			entry: `{"id":"EXP-0045","status":"planned","domain":"battle","question":"Unrelated wording"}`,
			want:  "battle_config is missing",
		},
		{
			name:  "non-battle experiment is exempt",
			entry: `{"id":"EXP-0046","status":"planned","question":"Where is the map header table?"}`,
		},
		{
			name:  "battle experiment before the cutover is exempt",
			entry: `{"id":"EXP-0040","status":"blocked","question":"Can the Whelk battle be won?"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "manifests/experiments.json",
				`{"schema_version":"1.0","experiments":[`+tt.entry+`]}`)
			fs, err := CheckBattleExperimentConfig(root)
			if err != nil {
				t.Fatal(err)
			}
			if tt.want == "" {
				if len(fs) != 0 {
					t.Fatalf("want no findings, got %v", fs)
				}
				return
			}
			if len(fs) != 1 || !strings.Contains(fs[0].Message, tt.want) {
				t.Fatalf("want one finding containing %q, got %v", tt.want, fs)
			}
		})
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
