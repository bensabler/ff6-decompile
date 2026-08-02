// Package audit implements the repository-integrity checks behind
// `ff6lab audit`: JSON manifest validation, manifest/index
// synchronization, broken relative Markdown links, restricted-file
// boundaries, and required tracked sources. Checks operate on a
// repository root passed in explicitly so tests can run against
// fixture trees.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bensabler/ff6-decompile/internal/census"
)

// Finding is one audit defect.
type Finding struct {
	Check   string
	Message string
}

func (f Finding) String() string { return f.Check + ": " + f.Message }

// Run executes every repository check and returns the findings.
func Run(root string) ([]Finding, error) {
	var all []Finding
	for _, fn := range []func(string) ([]Finding, error){
		CheckManifests,
		CheckExperimentRecordsInManifest,
		CheckExperimentIndexSync,
		CheckBattleExperimentConfig,
		CheckCensus,
		CheckMarkdownLinks,
		CheckRestrictedTracked,
		CheckRequiredTracked,
		CheckImportBoundaries,
		CheckTrackedBinaries,
		CheckReadinessSummary,
	} {
		fs, err := fn(root)
		if err != nil {
			return all, err
		}
		all = append(all, fs...)
	}
	return all, nil
}

type experimentsManifest struct {
	SchemaVersion string `json:"schema_version"`
	Experiments   []struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		Record        string `json:"record"`
		Question      string `json:"question"`
		Domain        string `json:"domain"`
		StartingState struct {
			BattleConfig map[string]any `json:"battle_config"`
		} `json:"starting_state"`
	} `json:"experiments"`
}

// CheckManifests parses every manifests/*.json and validates the
// minimal shape the project relies on.
func CheckManifests(root string) ([]Finding, error) {
	var fs []Finding
	paths, err := filepath.Glob(filepath.Join(root, "manifests", "*.json"))
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var generic map[string]any
		if err := json.Unmarshal(data, &generic); err != nil {
			fs = append(fs, Finding{"manifests", fmt.Sprintf("%s: invalid JSON: %v", rel(root, p), err)})
			continue
		}
		if _, ok := generic["schema_version"]; !ok {
			fs = append(fs, Finding{"manifests", rel(root, p) + ": missing schema_version"})
		}
	}
	em, err := loadExperiments(root)
	if err != nil {
		fs = append(fs, Finding{"manifests", err.Error()})
		return fs, nil
	}
	for _, e := range em.Experiments {
		if e.ID == "" || e.Status == "" {
			fs = append(fs, Finding{"manifests", "experiments.json: entry missing id or status"})
		}
		if e.Record != "" {
			if _, err := os.Stat(filepath.Join(root, e.Record)); err != nil {
				fs = append(fs, Finding{"manifests", fmt.Sprintf("experiments.json: %s record %s not found", e.ID, e.Record)})
			}
		}
	}
	return fs, nil
}

func loadExperiments(root string) (experimentsManifest, error) {
	var em experimentsManifest
	data, err := os.ReadFile(filepath.Join(root, "manifests", "experiments.json"))
	if err != nil {
		return em, fmt.Errorf("experiments.json unreadable: %w", err)
	}
	if err := json.Unmarshal(data, &em); err != nil {
		return em, fmt.Errorf("experiments.json invalid: %w", err)
	}
	return em, nil
}

// battleConfigFirstID is the experiment number from which battle
// experiments must carry an in-game configuration fingerprint. Earlier
// records predate the requirement and are not retrofitted, matching the
// "applies from EXP-NNNN onward" cutover docs/research/EVIDENCE_LAYOUT.md
// already uses.
const battleConfigFirstID = 41

var (
	expNumRe     = regexp.MustCompile(`^EXP-(\d{4})$`)
	battleWordRe = regexp.MustCompile(`(?i)\bbattles?\b|\bATB\b|\bcombat\b`)
)

// CheckBattleExperimentConfig requires every battle experiment from
// EXP-0041 onward to record starting_state.battle_config. Without it a
// run cannot say which Battle Mode or Battle Speed was in force, so its
// timing observations cannot be compared against another run's — the
// gap that stopped EXP-0040.
func CheckBattleExperimentConfig(root string) ([]Finding, error) {
	em, err := loadExperiments(root)
	if err != nil {
		return []Finding{{"battle-config", err.Error()}}, nil
	}
	var fs []Finding
	for _, e := range em.Experiments {
		m := expNumRe.FindStringSubmatch(e.ID)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil || n < battleConfigFirstID {
			continue
		}
		if !isBattleExperiment(e.Domain, e.Question, e.Record) {
			continue
		}
		if len(e.StartingState.BattleConfig) == 0 {
			fs = append(fs, Finding{"battle-config", fmt.Sprintf(
				"experiments.json: %s is a battle experiment but starting_state.battle_config is missing", e.ID)})
			continue
		}
		if src, _ := e.StartingState.BattleConfig["source"].(string); src == "" {
			fs = append(fs, Finding{"battle-config", fmt.Sprintf(
				"experiments.json: %s battle_config needs a source (memory-read, screen-read, or mixed)", e.ID)})
		}
	}
	return fs, nil
}

// isBattleExperiment reports whether an experiment must carry a battle
// configuration fingerprint. An explicit domain wins, so a record can
// opt out deliberately; otherwise the question and record path decide,
// so a battle experiment cannot escape the requirement by omission.
func isBattleExperiment(domain, question, record string) bool {
	if domain != "" {
		return strings.EqualFold(domain, "battle")
	}
	return battleWordRe.MatchString(question) || battleWordRe.MatchString(record)
}

var expRecordIDRe = regexp.MustCompile(`^EXP-\d{4}`)

// CheckExperimentRecordsInManifest verifies every record file under
// docs/experiments/EXP-*.md has a manifest entry — the reverse of
// CheckManifests' record-exists check, catching records committed
// without manifest registration.
func CheckExperimentRecordsInManifest(root string) ([]Finding, error) {
	em, err := loadExperiments(root)
	if err != nil {
		return []Finding{{"manifests", err.Error()}}, nil
	}
	known := make(map[string]bool, len(em.Experiments))
	for _, e := range em.Experiments {
		known[e.ID] = true
	}
	paths, err := filepath.Glob(filepath.Join(root, "docs", "experiments", "EXP-*.md"))
	if err != nil {
		return nil, err
	}
	var fs []Finding
	for _, p := range paths {
		id := expRecordIDRe.FindString(filepath.Base(p))
		if id != "" && !known[id] {
			fs = append(fs, Finding{"manifests", fmt.Sprintf("record %s has no experiments.json entry", rel(root, p))})
		}
	}
	return fs, nil
}

// CheckExperimentIndexSync verifies indexes/EXPERIMENTS.md lists every
// experiment in the manifest (by ID).
func CheckExperimentIndexSync(root string) ([]Finding, error) {
	var fs []Finding
	em, err := loadExperiments(root)
	if err != nil {
		return []Finding{{"index-sync", err.Error()}}, nil
	}
	idx, err := os.ReadFile(filepath.Join(root, "indexes", "EXPERIMENTS.md"))
	if err != nil {
		return []Finding{{"index-sync", "indexes/EXPERIMENTS.md unreadable"}}, nil
	}
	for _, e := range em.Experiments {
		if !strings.Contains(string(idx), e.ID) {
			fs = append(fs, Finding{"index-sync", "indexes/EXPERIMENTS.md missing " + e.ID})
		}
	}
	return fs, nil
}

// CheckCensus wraps the content-census and ROM-ledger validation
// (internal/census) into audit findings, and skips cleanly when the
// census manifests do not exist yet (pre-census fixture trees).
func CheckCensus(root string) ([]Finding, error) {
	if _, err := os.Stat(filepath.Join(root, "manifests", "content-census.json")); err != nil {
		return nil, nil
	}
	issues, err := census.Validate(root)
	if err != nil {
		return nil, err
	}
	fs := make([]Finding, 0, len(issues))
	for _, i := range issues {
		fs = append(fs, Finding{i.Check, i.Message})
	}
	return fs, nil
}

var linkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)\s#][^)#]*)\)`)

// CheckMarkdownLinks reports relative links whose targets do not exist.
// The historical lost-source citation in SESSION_001 is a documented
// permanent gap and is skipped.
func CheckMarkdownLinks(root string) ([]Finding, error) {
	var fs []Finding
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".md") {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		norm := regexp.MustCompile(`\s+`).ReplaceAllString(string(data), " ")
		for _, m := range linkRe.FindAllStringSubmatch(norm, -1) {
			u := strings.TrimSpace(m[1])
			if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") ||
				strings.HasPrefix(u, "mailto:") || strings.Contains(u, " ") {
				continue
			}
			if strings.HasSuffix(u, "FF6_Decompilation_Session_01_Summary.md") {
				continue // documented permanent provenance gap
			}
			target := filepath.Join(filepath.Dir(p), u)
			if _, err := os.Stat(target); err != nil {
				fs = append(fs, Finding{"md-links", fmt.Sprintf("%s -> %s", rel(root, p), u)})
			}
		}
		return nil
	})
	return fs, err
}

var restrictedExt = regexp.MustCompile(`(?i)\.(sfc|smc|fig|spc|srm|sav|state[0-9]*|mss|wav|flac|mp3|ogg)$`)

// renderedExt covers formats that carry decoded commercial graphics or audio.
// The demo renders ROM-derived frames, so these became reachable the moment
// `ff6demo -capture-last` existed; nothing in the repository tracks one
// today, which is exactly when a guard is cheap to add.
//
// Frame hashes and glyph hashes are tracked instead — they are verifiable
// without reproducing the bytes (docs/legal/ASSET_POLICY.md).
var renderedExt = regexp.MustCompile(`(?i)\.(png|bmp|gif|jpe?g|tiff?|pcm|brr)$`)

// CheckRestrictedTracked scans the git index for restricted file types.
func CheckRestrictedTracked(root string) ([]Finding, error) {
	files, err := gitLsFiles(root)
	if err != nil {
		return []Finding{{"restricted", "git unavailable: " + err.Error()}}, nil
	}
	var fs []Finding
	for _, f := range files {
		switch {
		case restrictedExt.MatchString(f):
			fs = append(fs, Finding{"restricted", "tracked restricted file: " + f})
		case renderedExt.MatchString(f):
			fs = append(fs, Finding{"restricted", "tracked rendered asset: " + f +
				" — decoded graphics and audio stay in local_artifacts/; track hashes instead (docs/legal/ASSET_POLICY.md)"})
		}
	}
	return fs, nil
}

// RequiredTracked are sources whose absence from the index breaks a
// clean clone (the .gitignore collision class of defect).
var RequiredTracked = []string{
	"cmd/ff6lab/main.go",
	"cmd/ff6demo/main.go",
	"local_artifacts/README.md",
	"go.mod",
	"CLAUDE.md",
}

// CheckRequiredTracked verifies the index contains the required files.
func CheckRequiredTracked(root string) ([]Finding, error) {
	files, err := gitLsFiles(root)
	if err != nil {
		return []Finding{{"required", "git unavailable: " + err.Error()}}, nil
	}
	set := make(map[string]bool, len(files))
	for _, f := range files {
		set[f] = true
	}
	var fs []Finding
	for _, want := range RequiredTracked {
		if !set[want] {
			fs = append(fs, Finding{"required", "not tracked: " + want})
		}
	}
	return fs, nil
}

func gitLsFiles(root string) ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(strings.TrimSpace(string(out))), nil
}

// GenerateExperimentIndex renders indexes/EXPERIMENTS.md from the
// manifest, preserving the project's table shape.
func GenerateExperimentIndex(root string) (string, error) {
	em, err := loadExperiments(root)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("# Experiments\n\n")
	b.WriteString("<!-- generated by `ff6lab indexes generate` from manifests/experiments.json -->\n\n")
	b.WriteString("| ID | Title | Status | Confidence | Record |\n|---|---|---|---|---|\n")
	for _, e := range em.Experiments {
		record := e.Record
		title := record
		if record != "" {
			base := filepath.Base(record)
			base = strings.TrimSuffix(base, ".md")
			title = strings.TrimPrefix(base, e.ID+"-")
			record = fmt.Sprintf("[%s](../%s)", base, e.Record)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | see record | %s |\n", e.ID, title, e.Status, record)
	}
	return b.String(), nil
}

func rel(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return r
}
