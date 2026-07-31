// Package census implements the breadth-tracking layer: the content
// census (what exists, per docs/research/CONTENT_TAXONOMY.md) and the
// ROM ownership ledger. It loads and validates the two manifests,
// computes coverage summaries and ROM gaps, and generates the derived
// indexes and dashboard.
package census

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Entry is one content-census record (manifests/content-census.json).
type Entry struct {
	SchemaVersion          string   `json:"schema_version"`
	ID                     string   `json:"id"`
	Domain                 string   `json:"domain"`
	Category               string   `json:"category"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	ReconstructionStatus   string   `json:"reconstruction_status"`
	RuntimeStatus          string   `json:"runtime_status"`
	RomRevision            string   `json:"rom_revision"`
	RomLocations           []string `json:"rom_locations"`
	WramLocations          []string `json:"wram_locations"`
	CpuRoutines            []string `json:"cpu_routines"`
	SpcLocations           []string `json:"spc_locations"`
	RecordCountExpected    *int     `json:"record_count_expected"`
	RecordCountFound       *int     `json:"record_count_found"`
	RecordSize             *int     `json:"record_size"`
	PointerTables          []string `json:"pointer_tables"`
	Consumers              []string `json:"consumers"`
	Producers              []string `json:"producers"`
	RelatedExperiments     []string `json:"related_experiments"`
	RelatedDiscoveries     []string `json:"related_discoveries"`
	RelatedImplementations []string `json:"related_implementations"`
	RelatedTests           []string `json:"related_tests"`
	RuntimeScenarios       []string `json:"runtime_scenarios"`
	Evidence               []string `json:"evidence"`
	Confidence             string   `json:"confidence"`
	UnknownFields          []string `json:"unknown_fields"`
	Contradictions         []string `json:"contradictions"`
	NextAction             string   `json:"next_action"`
	Notes                  string   `json:"notes"`
}

// Census is the manifest wrapper.
type Census struct {
	SchemaVersion string  `json:"schema_version"`
	Entries       []Entry `json:"entries"`
}

// Region is one ROM ownership record (manifests/rom-regions.json).
type Region struct {
	ID             string   `json:"id"`
	Start          int      `json:"start"`
	Size           int      `json:"size"`
	Bank           string   `json:"bank"`
	Classification string   `json:"classification"`
	Owner          string   `json:"owner"`
	Format         string   `json:"format"`
	Confidence     string   `json:"confidence"`
	Evidence       []string `json:"evidence"`
	Consumers      []string `json:"consumers"`
	Overlaps       []string `json:"overlaps"`
	Status         string   `json:"status"`
	Notes          string   `json:"notes"`
}

// Regions is the ledger wrapper.
type Regions struct {
	SchemaVersion string   `json:"schema_version"`
	RomSize       int      `json:"rom_size"`
	RomSHA256     string   `json:"rom_sha256"`
	Regions       []Region `json:"regions"`
}

// Issue is one validation defect (adapted to audit.Finding by the CLI).
type Issue struct {
	Check   string
	Message string
}

func (i Issue) String() string { return i.Check + ": " + i.Message }

// Controlled vocabularies. Order defines the ladder rank.
var ReconstructionLadder = []string{
	"UNMAPPED", "OBSERVED", "CANDIDATE_LOCATION", "LOCATED",
	"FORMAT_PARTIAL", "FORMAT_DECODED", "EXTRACTED_PARTIAL",
	"EXTRACTED_COMPLETE", "IMPLEMENTED_PARTIAL", "IMPLEMENTED_COMPLETE",
	"REGRESSION_TESTED", "DIFFERENTIALLY_VERIFIED",
}

var RuntimeLadder = []string{
	"NOT_ENCOUNTERED", "ENCOUNTERED", "INPUTS_CAPTURED",
	"OUTPUTS_CAPTURED", "NORMAL_PATH_VERIFIED", "ALTERNATE_PATH_VERIFIED",
	"EDGE_CASES_PARTIAL", "EDGE_CASES_COMPLETE", "EXHAUSTIVELY_VERIFIED",
}

var Domains = []string{
	"BATTLE", "CHAR", "MAGIC", "MONSTER", "ITEM", "WORLD",
	"EVENT", "MENU", "GFX", "AUDIO", "SAVE", "QUIRK",
}

var Classifications = []string{
	"CODE", "POINTER_TABLE", "DATA_TABLE", "TEXT", "GRAPHICS",
	"PALETTE", "TILEMAP", "ANIMATION", "MUSIC", "SOUND", "BRR",
	"EVENT_SCRIPT", "AI_SCRIPT", "HEADER", "PADDING", "UNKNOWN",
}

var Confidences = []string{
	"confirmed", "strong-hypothesis", "tentative-hypothesis", "unknown",
}

func rank(ladder []string, v string) int {
	for i, s := range ladder {
		if s == v {
			return i
		}
	}
	return -1
}

func has(set []string, v string) bool { return rank(set, v) >= 0 }

// AtLeast reports whether status v is at or above threshold t on the
// given ladder; unknown statuses always fail.
func AtLeast(ladder []string, v, t string) bool {
	rv, rt := rank(ladder, v), rank(ladder, t)
	return rv >= 0 && rt >= 0 && rv >= rt
}

// Load reads both manifests from the repository root.
func Load(root string) (Census, Regions, error) {
	var c Census
	var r Regions
	cb, err := os.ReadFile(filepath.Join(root, "manifests", "content-census.json"))
	if err != nil {
		return c, r, fmt.Errorf("content-census.json unreadable: %w", err)
	}
	if err := json.Unmarshal(cb, &c); err != nil {
		return c, r, fmt.Errorf("content-census.json invalid: %w", err)
	}
	rb, err := os.ReadFile(filepath.Join(root, "manifests", "rom-regions.json"))
	if err != nil {
		return c, r, fmt.Errorf("rom-regions.json unreadable: %w", err)
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return c, r, fmt.Errorf("rom-regions.json invalid: %w", err)
	}
	return c, r, nil
}

// idSet extracts the "id" field from a manifest list under key.
func idSet(root, file, key string) (map[string]bool, error) {
	data, err := os.ReadFile(filepath.Join(root, "manifests", file))
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	out := map[string]bool{}
	list, _ := doc[key].([]any)
	for _, it := range list {
		if m, ok := it.(map[string]any); ok {
			if id, ok := m["id"].(string); ok {
				out[id] = true
			}
		}
	}
	return out, nil
}

// Validate runs every census and ledger integrity check.
func Validate(root string) ([]Issue, error) {
	c, r, err := Load(root)
	if err != nil {
		return []Issue{{"census", err.Error()}}, nil
	}
	var issues []Issue
	add := func(check, format string, args ...any) {
		issues = append(issues, Issue{check, fmt.Sprintf(format, args...)})
	}

	expIDs, err := idSet(root, "experiments.json", "experiments")
	if err != nil {
		add("census", "experiments manifest: %v", err)
	}
	discIDs, err := idSet(root, "discoveries.json", "discoveries")
	if err != nil {
		add("census", "discoveries manifest: %v", err)
	}

	censusIDs := map[string]bool{}
	ownedDiscoveries := map[string]bool{}
	for _, e := range c.Entries {
		if censusIDs[e.ID] {
			add("census", "%s: duplicate id", e.ID)
		}
		censusIDs[e.ID] = true
		if !has(Domains, e.Domain) {
			add("census", "%s: invalid domain %q", e.ID, e.Domain)
		}
		if !has(ReconstructionLadder, e.ReconstructionStatus) {
			add("census", "%s: invalid reconstruction_status %q", e.ID, e.ReconstructionStatus)
		}
		if !has(RuntimeLadder, e.RuntimeStatus) {
			add("census", "%s: invalid runtime_status %q", e.ID, e.RuntimeStatus)
		}
		if !has(Confidences, e.Confidence) {
			add("census", "%s: invalid confidence %q", e.ID, e.Confidence)
		}
		if e.Name == "" || e.Category == "" || e.NextAction == "" {
			add("census", "%s: name, category, and next_action are required", e.ID)
		}
		if AtLeast(ReconstructionLadder, e.ReconstructionStatus, "OBSERVED") && len(e.Evidence) == 0 {
			add("census", "%s: status %s requires evidence", e.ID, e.ReconstructionStatus)
		}
		if AtLeast(ReconstructionLadder, e.ReconstructionStatus, "IMPLEMENTED_PARTIAL") && len(e.RelatedImplementations) == 0 {
			add("census", "%s: status %s requires related_implementations", e.ID, e.ReconstructionStatus)
		}
		if AtLeast(ReconstructionLadder, e.ReconstructionStatus, "REGRESSION_TESTED") && len(e.RelatedTests) == 0 {
			add("census", "%s: status %s requires related_tests", e.ID, e.ReconstructionStatus)
		}
		for _, x := range e.RelatedExperiments {
			if expIDs != nil && !expIDs[x] {
				add("census", "%s: unknown experiment %s", e.ID, x)
			}
		}
		for _, d := range e.RelatedDiscoveries {
			if discIDs != nil && !discIDs[d] {
				add("census", "%s: unknown discovery %s", e.ID, d)
			}
			ownedDiscoveries[d] = true
		}
	}
	for d := range discIDs {
		if !ownedDiscoveries[d] {
			add("census", "discovery %s has no census ownership", d)
		}
	}

	regionIDs := map[string]bool{}
	for _, g := range r.Regions {
		if regionIDs[g.ID] {
			add("rom-regions", "%s: duplicate id", g.ID)
		}
		regionIDs[g.ID] = true
		if !has(Classifications, g.Classification) {
			add("rom-regions", "%s: invalid classification %q", g.ID, g.Classification)
		}
		if !has(Confidences, g.Confidence) {
			add("rom-regions", "%s: invalid confidence %q", g.ID, g.Confidence)
		}
		if g.Status != "known" && g.Status != "candidate" {
			add("rom-regions", "%s: invalid status %q", g.ID, g.Status)
		}
		if g.Size <= 0 || g.Start < 0 || g.Start+g.Size > r.RomSize {
			add("rom-regions", "%s: span [%#x,+%d) outside the %d-byte ROM", g.ID, g.Start, g.Size, r.RomSize)
		}
		if g.Status == "known" && len(g.Evidence) == 0 {
			add("rom-regions", "%s: known region requires evidence", g.ID)
		}
	}

	// Undeclared overlap detection.
	sorted := append([]Region(nil), r.Regions...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start < sorted[j].Start })
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			a, b := sorted[i], sorted[j]
			if b.Start >= a.Start+a.Size {
				break
			}
			declared := has(a.Overlaps, b.ID) && has(b.Overlaps, a.ID)
			if !declared {
				add("rom-regions", "%s and %s overlap without declaring it", a.ID, b.ID)
			}
		}
	}

	// Cross-links: census rom_locations naming ROM ids must exist, and
	// region owners naming CEN ids must exist.
	for _, e := range c.Entries {
		for _, loc := range e.RomLocations {
			if len(loc) > 4 && loc[:4] == "ROM-" && !regionIDs[loc] {
				add("census", "%s: rom_locations names unknown region %s", e.ID, loc)
			}
		}
	}
	for _, g := range r.Regions {
		if len(g.Owner) > 4 && g.Owner[:4] == "CEN-" && !censusIDs[g.Owner] {
			add("rom-regions", "%s: owner names unknown census entry %s", g.ID, g.Owner)
		}
	}

	// Data inventories must parse and follow the shared shape.
	inv, err := filepath.Glob(filepath.Join(root, "data", "census", "*.json"))
	if err == nil {
		for _, p := range inv {
			data, err := os.ReadFile(p)
			if err != nil {
				add("data-census", "%s unreadable: %v", filepath.Base(p), err)
				continue
			}
			var doc map[string]any
			if err := json.Unmarshal(data, &doc); err != nil {
				add("data-census", "%s: invalid JSON: %v", filepath.Base(p), err)
				continue
			}
			for _, k := range []string{"schema_version", "domain", "description", "records", "next_extraction_step"} {
				if _, ok := doc[k]; !ok {
					add("data-census", "%s: missing %s", filepath.Base(p), k)
				}
			}
			if d, ok := doc["domain"].(string); ok && !has(Domains, d) {
				add("data-census", "%s: invalid domain %q", filepath.Base(p), d)
			}
		}
	}

	// Generated files must match regeneration output (staleness).
	for name, gen := range map[string]func(Census, Regions) string{
		"indexes/CONTENT_CENSUS.md": func(c Census, _ Regions) string { return GenerateCensusIndex(c) },
		"indexes/ROM_REGIONS.md":    func(_ Census, r Regions) string { return GenerateRegionIndex(r) },
		"dashboards/COVERAGE.md":    GenerateCoverageDashboard,
	} {
		want := gen(c, r)
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			add("census-sync", "%s missing (run `ff6lab census sync`)", name)
		} else if string(got) != want {
			add("census-sync", "%s is stale (run `ff6lab census sync`)", name)
		}
	}

	return issues, nil
}
