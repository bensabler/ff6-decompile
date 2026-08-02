package workflow

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureRoot(t *testing.T) string { t.Helper(); return t.TempDir() }

func startedRun(t *testing.T, root string, reqs ...Requirement) (*Store, *Contract) {
	t.Helper()
	s := NewStore(root)
	c := frozen(t, reqs...)
	if _, err := s.Start(c, "2026-08-02T00:00:00Z"); err != nil {
		t.Fatalf("start: %v", err)
	}
	return s, c
}

func TestStartRefusesUnfrozenContract(t *testing.T) {
	s := NewStore(fixtureRoot(t))
	c := validContract() // Draft
	if _, err := s.Start(c, "t0"); err == nil {
		t.Fatal("a run must not begin from an unfrozen contract")
	}
}

func TestStartRefusesToOverwriteAnExistingRun(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	if _, err := s.Start(c, "t1"); err == nil {
		t.Error("starting an existing run must fail: amend it rather than overwriting")
	}
}

func TestStartRejectsMalformedWorkflowID(t *testing.T) {
	s := NewStore(fixtureRoot(t))
	c := frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
	c.WorkflowID = "extract-run-1"
	if _, err := s.Start(c, "t0"); err == nil {
		t.Error("workflow_id must match WF-NNNN")
	}
}

// State survives a process restart: reloading is enough to answer what is open
// and what remains, without replaying anything.
func TestRestartSafeContinuation(t *testing.T) {
	root := fixtureRoot(t)
	reqs := []Requirement{
		agentReq("graphics-researcher", Required, BlockCompletion),
		agentReq("asset-librarian", Required, BlockCompletion),
	}
	_, c := startedRun(t, root, reqs...)

	// A fresh Store, as a later session would build.
	reloaded := NewStore(root)
	st, err := reloaded.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatalf("reload state: %v", err)
	}
	if st.Phase != PhaseExecuting {
		t.Errorf("phase = %q, want %q", st.Phase, PhaseExecuting)
	}
	got, err := reloaded.LoadContract(c.WorkflowID)
	if err != nil {
		t.Fatalf("reload contract: %v", err)
	}
	if err := got.VerifyFrozen(); err != nil {
		t.Errorf("a round-tripped contract must still verify: %v", err)
	}
	if rem := st.Remaining(got); len(rem) != 2 {
		t.Errorf("remaining = %v, want both requirements before any evidence", rem)
	}

	ids, err := reloaded.List()
	if err != nil || len(ids) != 1 || ids[0] != c.WorkflowID {
		t.Errorf("List() = %v, %v", ids, err)
	}
}

func TestCloseDetectsContractSwappedAfterStart(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))

	// Replace the stored contract with a different frozen one, as an operator or
	// a process might by editing the tracked file.
	other := frozen(t, agentReq("asset-librarian", Required, BlockCompletion))
	if err := WriteJSON(s.ContractPath(c.WorkflowID), other); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Close(c.WorkflowID, Evidence{}, "t1"); err == nil {
		t.Error("closing must fail when the contract changed since the run started")
	}
}

func TestCloseRecordsComputedVerdictAndRemaining(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root,
		agentReq("graphics-researcher", Required, BlockCompletion),
		agentReq("asset-librarian", Required, BlockCompletion))

	ev := Evidence{Observations: []Observation{
		{Kind: ObsAgentCall, Selector: "graphics-researcher", EvidenceRef: "t.jsonl"}}}
	rec, err := s.Close(c.WorkflowID, ev, "t1")
	if err != nil {
		t.Fatalf("close: %v", err)
	}
	if rec.Verdict != Partial {
		t.Errorf("verdict = %q, want partial: one required agent was never invoked", rec.Verdict)
	}
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Phase != PhaseClosed || st.Reconciliation == nil {
		t.Fatalf("phase=%q reconciliation=%v", st.Phase, st.Reconciliation)
	}
	rem := st.Remaining(c)
	if len(rem) != 1 || rem[0] != "asset-librarian" {
		t.Errorf("remaining = %v, want [asset-librarian]", rem)
	}
}

// N6: a required invocation that cannot be verified yields unverifiable or
// partial, never complete.
func TestN6MissingTranscriptEvidenceIsUnverifiable(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))

	ev, err := TranscriptObservations(filepath.Join(root, "no-such-transcripts"))
	if err != nil {
		t.Fatalf("transcript observations: %v", err)
	}
	if _, ok := ev.Incomplete[ObsAgentCall]; !ok {
		t.Fatal("a missing transcript directory must be reported as incomplete evidence")
	}
	rec, err := s.Close(c.WorkflowID, ev, "t1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict == Complete {
		t.Error("a run whose invocations cannot be verified must never be complete")
	}
	if rec.Verdict != Unverifiable {
		t.Errorf("verdict = %q, want unverifiable", rec.Verdict)
	}
}

// N9: a receipt claiming a verdict the verifier did not compute fails
// validation, and the computed verdict stands.
func TestN9ReceiptDisagreesWithVerifier(t *testing.T) {
	tests := []struct {
		name     string
		receipt  string
		computed Verdict
		wantErr  bool
	}{
		{name: "receipt agrees", receipt: "## Result\n\nVerdict: complete\n",
			computed: Complete},
		{name: "receipt claims complete, verifier says partial",
			receipt: "**Verdict:** complete\n", computed: Partial, wantErr: true},
		{name: "receipt claims complete, verifier says failed",
			receipt: "- verdict = `complete`\n", computed: Failed, wantErr: true},
		{name: "receipt states no verdict at all",
			receipt: "Everything went fine.\n", computed: Complete, wantErr: true},
		{name: "bulleted and bolded form is recognised",
			receipt: "- **Verdict**: `partial`\n", computed: Partial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReceipt(tt.receipt, tt.computed)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want a validation failure")
				}
				if !errors.Is(err, ErrReceiptDisagrees) {
					t.Errorf("want ErrReceiptDisagrees, got %v", err)
				}
				return
			}
			if err != nil {
				t.Errorf("want agreement, got %v", err)
			}
		})
	}
}

func TestValidateReceiptFileMissingReceiptIsNotAnError(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	if err := s.ValidateReceiptFile(c.WorkflowID, Complete); err != nil {
		t.Errorf("a run may be reconciled before a receipt exists: %v", err)
	}

	if err := os.WriteFile(s.ReceiptPath(c.WorkflowID),
		[]byte("Verdict: complete\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.ValidateReceiptFile(c.WorkflowID, Failed); err == nil {
		t.Error("a receipt disagreeing with the computed verdict must fail validation")
	}
}

func TestGateLogObservations(t *testing.T) {
	root := fixtureRoot(t)
	dir := filepath.Join(root, "gate-logs", "run")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := strings.Join([]string{
		"label\trc\treq\tstart\tend\tcwd\thead\tout_sha256",
		"gofmt\t0\trequired\t2026-08-02T00:00:00Z\tx\ty\tz\tw",
		"go-test\t2\trequired\t2026-08-02T00:00:01Z\tx\ty\tz\tw",
		"census\t\trequired\t2026-08-02T00:00:02Z\tx\ty\tz\tw", // blank status
		"",
	}, "\n")
	path := filepath.Join(dir, "status.tsv")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	ev, err := GateLogObservations(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ev.Observations) != 3 {
		t.Fatalf("got %d observations, want 3", len(ev.Observations))
	}
	byID := map[string]Observation{}
	for _, o := range ev.Observations {
		byID[o.Selector] = o
	}
	if got := byID["gofmt"].ExitStatus; got == nil || *got != 0 {
		t.Errorf("gofmt exit = %v, want 0", got)
	}
	if got := byID["go-test"].ExitStatus; got == nil || *got != 2 {
		t.Errorf("go-test exit = %v, want 2", got)
	}
	// The defect AUDIT-0001 shipped twice: a blank status is not a pass.
	if byID["census"].ExitStatus != nil {
		t.Error("a blank exit status must stay nil, never default to 0")
	}
}

func TestGateLogMissingIsIncompleteNotEmpty(t *testing.T) {
	ev, err := GateLogObservations(filepath.Join(t.TempDir(), "absent.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ev.Incomplete[ObsBackendRun]; !ok {
		t.Error("a missing gate log must be reported as incomplete, not as zero backends run")
	}
}

func TestTranscriptObservationsReadsOnlyToolMetadata(t *testing.T) {
	dir := t.TempDir()
	// A realistic record carrying conversation text that must never be surfaced.
	line := `{"timestamp":"2026-08-02T15:04:48Z","message":{"content":[` +
		`{"type":"text","text":"SECRET CONVERSATION CONTENT"},` +
		`{"type":"tool_use","name":"Skill","input":{"skill":"run-quality-gates"}},` +
		`{"type":"tool_use","name":"Agent","input":{"subagent_type":"graphics-researcher"}}` +
		`]}}`
	if err := os.WriteFile(filepath.Join(dir, "s.jsonl"),
		[]byte(line+"\nnot json\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ev, err := TranscriptObservations(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ev.Incomplete) != 0 {
		t.Errorf("a readable transcript must not be marked incomplete: %v", ev.Incomplete)
	}
	if len(ev.Observations) != 2 {
		t.Fatalf("got %d observations, want 2", len(ev.Observations))
	}
	for _, o := range ev.Observations {
		if strings.Contains(o.Selector, "SECRET") || strings.Contains(o.EvidenceRef, "SECRET") {
			t.Fatal("conversation content must never reach an observation")
		}
	}
	if ev.Observations[0].Kind != ObsSkillCall || ev.Observations[0].Selector != "run-quality-gates" {
		t.Errorf("skill observation = %+v", ev.Observations[0])
	}
	if ev.Observations[1].Kind != ObsAgentCall || ev.Observations[1].Selector != "graphics-researcher" {
		t.Errorf("agent observation = %+v", ev.Observations[1])
	}
}

func TestMergePreservesEveryGap(t *testing.T) {
	a := Evidence{Observations: []Observation{{Kind: ObsAgentCall, Selector: "x"}}}
	b := Evidence{Incomplete: map[ObservationKind]string{ObsBackendRun: "no gate log"}}
	got := Merge(a, b)
	if len(got.Observations) != 1 {
		t.Errorf("observations = %d, want 1", len(got.Observations))
	}
	if got.Incomplete[ObsBackendRun] == "" {
		t.Error("one source seeing clearly must not erase another source's blindness")
	}
}

func TestOutputObservationsHashPresentFilesOnly(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "out"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "out", "there.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	obs, err := OutputObservations(root, []string{"out/there.txt", "out/missing.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if len(obs) != 1 {
		t.Fatalf("got %d observations, want only the present file", len(obs))
	}
	if obs[0].Kind != ObsArtifactPresent {
		t.Errorf("kind = %q; an output's presence is an artifact, never an invocation", obs[0].Kind)
	}
	if !strings.HasPrefix(obs[0].EvidenceRef, "sha256:") {
		t.Errorf("evidence ref = %q, want a content hash", obs[0].EvidenceRef)
	}
}
