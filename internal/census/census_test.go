package census

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

const minimalEntry = `{
  "schema_version": "1.0", "id": "CEN-BATTLE-0001", "domain": "BATTLE",
  "category": "damage formulas", "name": "Damage pipeline",
  "description": "d", "reconstruction_status": "LOCATED",
  "runtime_status": "ENCOUNTERED", "confidence": "confirmed",
  "evidence": ["docs/experiments/EXP-0001-x.md"],
  "related_experiments": ["EXP-0001"], "related_discoveries": ["DISC-0001"],
  "rom_locations": ["ROM-0001"], "next_action": "n"
}`

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "manifests/experiments.json",
		`{"schema_version":"1.0","experiments":[{"id":"EXP-0001","status":"completed"}]}`)
	writeFile(t, root, "manifests/discoveries.json",
		`{"schema_version":"1.0","discoveries":[{"id":"DISC-0001"}]}`)
	writeFile(t, root, "manifests/content-census.json",
		`{"schema_version":"1.0","entries":[`+minimalEntry+`]}`)
	writeFile(t, root, "manifests/rom-regions.json",
		`{"schema_version":"1.0","rom_size":1000,"rom_sha256":"x","regions":[
		  {"id":"ROM-0001","start":0,"size":100,"classification":"CODE","status":"known","confidence":"confirmed","evidence":["e"]},
		  {"id":"ROM-0002","start":200,"size":100,"classification":"DATA_TABLE","status":"candidate","confidence":"tentative-hypothesis"}
		]}`)
	// generated files in sync
	c, r, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "indexes/CONTENT_CENSUS.md", GenerateCensusIndex(c))
	writeFile(t, root, "indexes/ROM_REGIONS.md", GenerateRegionIndex(r))
	writeFile(t, root, "dashboards/COVERAGE.md", GenerateCoverageDashboard(c, r))
	return root
}

func TestValidateClean(t *testing.T) {
	issues, err := Validate(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Fatalf("want no issues, got %v", issues)
	}
}

func TestValidateCatchesDefects(t *testing.T) {
	root := fixture(t)
	// duplicate id + unknown experiment + missing evidence + undeclared overlap + stale index
	writeFile(t, root, "manifests/content-census.json",
		`{"schema_version":"1.0","entries":[`+minimalEntry+`,`+minimalEntry+`,
		 {"schema_version":"1.0","id":"CEN-MAGIC-0001","domain":"MAGIC","category":"c","name":"n",
		  "description":"d","reconstruction_status":"OBSERVED","runtime_status":"ENCOUNTERED",
		  "confidence":"unknown","related_experiments":["EXP-9999"],"next_action":"n"}]}`)
	writeFile(t, root, "manifests/rom-regions.json",
		`{"schema_version":"1.0","rom_size":1000,"rom_sha256":"x","regions":[
		  {"id":"ROM-0001","start":0,"size":100,"classification":"CODE","status":"known","confidence":"confirmed","evidence":["e"]},
		  {"id":"ROM-0003","start":50,"size":100,"classification":"TEXT","status":"known","confidence":"confirmed","evidence":["e"]}
		]}`)
	issues, err := Validate(root)
	if err != nil {
		t.Fatal(err)
	}
	text := ""
	for _, i := range issues {
		text += i.String() + "\n"
	}
	for _, want := range []string{
		"duplicate id", "unknown experiment EXP-9999",
		"requires evidence", "overlap without declaring",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in:\n%s", want, text)
		}
	}
	if !strings.Contains(text, "stale") {
		t.Fatalf("expected staleness findings, got:\n%s", text)
	}
}

func TestDiscoveryOwnership(t *testing.T) {
	root := fixture(t)
	writeFile(t, root, "manifests/discoveries.json",
		`{"schema_version":"1.0","discoveries":[{"id":"DISC-0001"},{"id":"DISC-0002"}]}`)
	issues, err := Validate(root)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, i := range issues {
		if strings.Contains(i.Message, "DISC-0002 has no census ownership") {
			found = true
		}
	}
	if !found {
		t.Fatalf("unowned discovery not reported: %v", issues)
	}
}

func TestAnalyzeGapsAndOverlapDedup(t *testing.T) {
	r := Regions{RomSize: 1000, Regions: []Region{
		{ID: "ROM-0001", Start: 0, Size: 100, Classification: "CODE", Status: "known"},
		{ID: "ROM-0002", Start: 50, Size: 100, Classification: "DATA_TABLE", Status: "known"},
		{ID: "ROM-0003", Start: 400, Size: 100, Classification: "TEXT", Status: "candidate"},
	}}
	st := Analyze(r)
	if st.KnownBytes != 150 {
		t.Fatalf("known bytes = %d, want 150 (overlap deduplicated)", st.KnownBytes)
	}
	if st.CandidateBytes != 100 {
		t.Fatalf("candidate bytes = %d, want 100", st.CandidateBytes)
	}
	wantGaps := []Gap{{150, 250}, {500, 500}}
	if len(st.Gaps) != 2 || st.Gaps[0] != wantGaps[0] || st.Gaps[1] != wantGaps[1] {
		t.Fatalf("gaps = %v, want %v", st.Gaps, wantGaps)
	}
}

func TestAtLeast(t *testing.T) {
	if !AtLeast(ReconstructionLadder, "REGRESSION_TESTED", "LOCATED") {
		t.Fatal("REGRESSION_TESTED should be >= LOCATED")
	}
	if AtLeast(ReconstructionLadder, "OBSERVED", "LOCATED") {
		t.Fatal("OBSERVED should be < LOCATED")
	}
	if AtLeast(ReconstructionLadder, "bogus", "LOCATED") {
		t.Fatal("unknown status must fail")
	}
}

func TestSummaryCountsAllDomains(t *testing.T) {
	c := Census{Entries: []Entry{
		{Domain: "MAGIC", ReconstructionStatus: "FORMAT_DECODED", RuntimeStatus: "NORMAL_PATH_VERIFIED"},
	}}
	sum := Summary(c)
	if len(sum) != len(Domains) {
		t.Fatalf("summary rows = %d, want %d", len(sum), len(Domains))
	}
	for _, d := range sum {
		if d.Domain == "MAGIC" {
			if d.Registered != 1 || d.Decoded != 1 || d.RuntimeVerified != 1 || d.Implemented != 0 {
				t.Fatalf("MAGIC row wrong: %+v", d)
			}
		}
	}
}
