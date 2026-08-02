package workflow

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Evidence collection for R14. Each source answers for exactly one observation
// kind, and each reports its own incompleteness rather than returning an empty
// set that would read as "observed absent". A source that cannot see is not a
// source that saw nothing.

// TranscriptObservations extracts agent and skill invocations from harness
// session transcripts in dir.
//
// Only tool-call metadata is read: the tool name and its selector. Prompts,
// responses and conversation content are never parsed, copied or returned.
//
// A missing or empty directory marks agent and skill observations Incomplete,
// so requirements discharged through them become Unverifiable rather than
// Unsatisfied. This is the N6 path: an invocation that cannot be verified must
// never read as an invocation that did not happen.
func TranscriptObservations(dir string) (Evidence, error) {
	ev := Evidence{Incomplete: map[ObservationKind]string{}}

	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		markTranscriptGap(&ev, "no transcript directory at "+dir)
		return ev, nil
	}
	if err != nil {
		return ev, fmt.Errorf("read transcript dir %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	if len(files) == 0 {
		markTranscriptGap(&ev, "no .jsonl transcripts in "+dir)
		return ev, nil
	}

	for _, f := range files {
		obs, err := transcriptFile(f)
		if err != nil {
			return ev, err
		}
		ev.Observations = append(ev.Observations, obs...)
	}
	return ev, nil
}

func markTranscriptGap(ev *Evidence, why string) {
	ev.Incomplete[ObsAgentCall] = why
	ev.Incomplete[ObsSkillCall] = why
}

// transcriptLine is the narrow slice of a transcript record this package reads.
// Everything else in the record, including all message text, is ignored.
type transcriptLine struct {
	Timestamp string `json:"timestamp"`
	Message   struct {
		Content []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Input struct {
				Skill        string `json:"skill"`
				SubagentType string `json:"subagent_type"`
			} `json:"input"`
		} `json:"content"`
	} `json:"message"`
}

func transcriptFile(path string) ([]Observation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open transcript %s: %w", path, err)
	}
	defer f.Close()

	ref := filepath.Base(path)
	var obs []Observation
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		var rec transcriptLine
		// A record this package cannot parse is skipped rather than fatal: a
		// transcript is an external artifact and its shape may drift.
		if err := json.Unmarshal(sc.Bytes(), &rec); err != nil {
			continue
		}
		for _, b := range rec.Message.Content {
			if b.Type != "tool_use" {
				continue
			}
			switch {
			case b.Name == "Skill" && b.Input.Skill != "":
				obs = append(obs, Observation{Kind: ObsSkillCall, Selector: b.Input.Skill,
					Timestamp: rec.Timestamp, EvidenceRef: ref})
			case b.Name == "Agent" && b.Input.SubagentType != "":
				obs = append(obs, Observation{Kind: ObsAgentCall, Selector: b.Input.SubagentType,
					Timestamp: rec.Timestamp, EvidenceRef: ref})
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan transcript %s: %w", path, err)
	}
	return obs, nil
}

// GateLogObservations reads backend runs from a gate harness status.tsv.
//
// The harness records one row per gate with its real captured exit status. A
// row whose status column is empty yields an observation with a nil ExitStatus,
// which reconciliation treats as Unverifiable — never as a pass. That is the
// defect AUDIT-0001 shipped twice.
func GateLogObservations(path string) (Evidence, error) {
	ev := Evidence{Incomplete: map[ObservationKind]string{}}

	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		ev.Incomplete[ObsBackendRun] = "no gate log at " + path
		return ev, nil
	}
	if err != nil {
		return ev, fmt.Errorf("open gate log %s: %w", path, err)
	}
	defer f.Close()

	ref := filepath.Base(filepath.Dir(path))
	sc := bufio.NewScanner(f)
	for n := 0; sc.Scan(); n++ {
		line := sc.Text()
		if n == 0 && strings.HasPrefix(line, "label\t") {
			continue // header
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		col := strings.Split(line, "\t")
		if len(col) < 2 {
			continue
		}
		o := Observation{Kind: ObsBackendRun, Selector: col[0], EvidenceRef: ref}
		if rc, err := strconv.Atoi(strings.TrimSpace(col[1])); err == nil {
			o.ExitStatus = &rc
		}
		if len(col) > 3 {
			o.Timestamp = col[3]
		}
		ev.Observations = append(ev.Observations, o)
	}
	if err := sc.Err(); err != nil {
		return ev, fmt.Errorf("scan gate log %s: %w", path, err)
	}
	return ev, nil
}

// OutputObservations records which of a contract's required outputs exist under
// root, with their content hashes. An output's presence is an artifact
// observation and never satisfies an invocation requirement.
func OutputObservations(root string, outputs []string) ([]Observation, error) {
	var obs []Observation
	for _, rel := range outputs {
		p := filepath.Join(root, filepath.FromSlash(rel))
		b, err := os.ReadFile(p)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read declared output %s: %w", rel, err)
		}
		sum := sha256.Sum256(b)
		obs = append(obs, Observation{Kind: ObsArtifactPresent, Selector: rel,
			EvidenceRef: "sha256:" + hex.EncodeToString(sum[:])})
	}
	return obs, nil
}

// Merge combines evidence sets, preserving every incompleteness marker. A gap
// reported by any source stays reported: one source seeing clearly does not
// make another source's blindness disappear.
func Merge(sets ...Evidence) Evidence {
	out := Evidence{Incomplete: map[ObservationKind]string{}}
	for _, s := range sets {
		out.Observations = append(out.Observations, s.Observations...)
		for k, v := range s.Incomplete {
			if _, exists := out.Incomplete[k]; !exists {
				out.Incomplete[k] = v
			}
		}
	}
	if len(out.Incomplete) == 0 {
		out.Incomplete = nil
	}
	return out
}

// WriteJSON writes v as indented JSON, creating parent directories.
func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dir for %s: %w", path, err)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ReadJSON reads indented JSON into v.
func ReadJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}
