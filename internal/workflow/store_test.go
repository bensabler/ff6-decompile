package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "README.md")
	runGit(t, root, "-c", "user.name=Workflow Test", "-c", "user.email=workflow@example.invalid",
		"commit", "-q", "-m", "fixture")
	return root
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmdArgs := append([]string{"-C", root}, args...)
	if out, err := exec.Command("git", cmdArgs...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func startedRun(t *testing.T, root string, reqs ...Requirement) (*Store, *Contract) {
	t.Helper()
	s := NewStore(root)
	c := frozen(t, reqs...)
	if _, err := s.Start(c, "2026-08-02T00:00:00Z"); err != nil {
		t.Fatalf("start: %v", err)
	}
	return s, c
}

func nextEvent(t *testing.T, s *Store, id string, kind EventKind, selector string,
	source SourceKind, trust TrustBasis) RunEvent {
	t.Helper()
	st, err := s.LoadState(id)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := s.LoadRunIdentity(id)
	if err != nil {
		t.Fatal(err)
	}
	eventID, err := newEventID()
	if err != nil {
		t.Fatal(err)
	}
	created, err := time.Parse(time.RFC3339Nano, identity.CreatedAt)
	if err != nil {
		t.Fatal(err)
	}
	return RunEvent{
		SchemaVersion: RunEventSchemaVersion, Sequence: st.LedgerSequence + 1,
		EventID: eventID, WorkflowID: identity.WorkflowID, RunID: identity.RunID,
		ContractHash: identity.ContractHash,
		ObservedAt:   created.Add(time.Duration(st.LedgerSequence+1) * time.Second).Format(time.RFC3339Nano),
		Provider:     "fixture-provider", SourceKind: source, CollectorID: "fixture-collector",
		TrustBasis: trust, SessionID: "session-1", TurnID: "turn-1",
		RepositoryIdentity: identity.RepositoryIdentity, CWD: identity.RepositoryRoot,
		Branch: identity.StartingBranch, Head: identity.StartingHead,
		EventKind: kind, Selector: selector, ToolUseID: "tool-1", EvidenceRef: "",
	}
}

func appendEvent(t *testing.T, s *Store, id string, event RunEvent) {
	t.Helper()
	if err := s.appendEvent(id, event); err != nil {
		t.Fatalf("append event: %v", err)
	}
}

func TestStartCreatesImmutableRunIdentityAndRestrictiveLedger(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if !runIDPattern.MatchString(st.RunID) || st.RunID == c.WorkflowID {
		t.Fatalf("generated run_id = %q", st.RunID)
	}
	identity, err := s.LoadRunIdentity(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := canonicalPath(root)
	if err != nil {
		t.Fatal(err)
	}
	if identity.ContractHash != c.FrozenHash || identity.WorkflowID != c.WorkflowID ||
		identity.RepositoryRoot != canonicalRoot || identity.StartingBranch != "main" {
		t.Errorf("identity = %+v", identity)
	}
	if !strings.HasPrefix(identity.RepositoryIdentity, "local-git-common-dir:") {
		t.Errorf("repository without remote must be identified honestly, got %q", identity.RepositoryIdentity)
	}
	if identity.RepositoryIdentity == filepath.Base(root) {
		t.Error("repository identity must not be only the directory name")
	}
	for _, path := range []string{s.IdentityPath(c.WorkflowID, st.RunID), s.EventsPath(c.WorkflowID, st.RunID), s.StatePath(c.WorkflowID)} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("%s mode = %o, want 600", path, got)
		}
	}
	if got := s.EventsPath(c.WorkflowID, st.RunID); !strings.Contains(got,
		filepath.Join("local_artifacts", "workflows", c.WorkflowID, st.RunID, "events.jsonl")) {
		t.Errorf("ledger path = %s", got)
	}
}

func TestRepositoryRemoteIdentityStripsCredentials(t *testing.T) {
	root := fixtureRoot(t)
	runGit(t, root, "remote", "add", "origin", "https://secret@example.com/owner/repo.git?token=private")
	ctx, err := resolveRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(ctx.Identity, "secret") || strings.Contains(ctx.Identity, "private") {
		t.Fatalf("repository identity leaked remote credentials: %q", ctx.Identity)
	}
	if ctx.Identity != "remote:https://example.com/owner/repo" {
		t.Errorf("repository identity = %q", ctx.Identity)
	}
}

func TestLinkedWorktreesShareRepositoryIdentity(t *testing.T) {
	root := fixtureRoot(t)
	linked := filepath.Join(filepath.Dir(root), "linked-worktree")
	runGit(t, root, "worktree", "add", "-q", "-b", "linked-test", linked)
	rootContext, err := resolveRepository(root)
	if err != nil {
		t.Fatal(err)
	}
	linkedContext, err := resolveRepository(linked)
	if err != nil {
		t.Fatal(err)
	}
	if rootContext.Identity != linkedContext.Identity {
		t.Errorf("linked worktree identity %q differs from repository identity %q",
			linkedContext.Identity, rootContext.Identity)
	}
	if rootContext.Root == linkedContext.Root {
		t.Error("worktree roots should remain distinct even when Git identity is shared")
	}
}

func TestStartRefusesUnfrozenContract(t *testing.T) {
	s := NewStore(fixtureRoot(t))
	c := validContract()
	if _, err := s.Start(c, "2026-08-02T00:00:00Z"); err == nil {
		t.Fatal("a run must not begin from an unfrozen contract")
	}
}

func TestStartRefusesToOverwriteAnExistingRun(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	if _, err := s.Start(c, "2026-08-02T00:00:01Z"); err == nil {
		t.Error("starting an existing run must fail")
	}
}

func TestStartRejectsMalformedWorkflowID(t *testing.T) {
	s := NewStore(fixtureRoot(t))
	c := frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
	c.WorkflowID = "extract-run-1"
	if _, err := s.Start(c, "2026-08-02T00:00:00Z"); err == nil {
		t.Error("workflow_id must match WF-NNNN")
	}
}

func TestStartRollsBackWhenLedgerCreationFailsAfterIdentity(t *testing.T) {
	root := fixtureRoot(t)
	s := NewStore(root)
	s.generateRunID = func() (string, error) { return "run-00000000000000000000000000000001", nil }
	identityExisted := false
	s.createLedger = func(path string) error {
		identityPath := filepath.Join(filepath.Dir(path), identityFile)
		_, statErr := os.Stat(identityPath)
		identityExisted = statErr == nil
		return errors.New("injected ledger creation failure")
	}
	c := frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
	if _, err := s.Start(c, "2026-08-02T00:00:00Z"); err == nil {
		t.Fatal("want start failure")
	}
	if !identityExisted {
		t.Error("test did not reach the identity-created boundary")
	}
	if _, err := os.Stat(filepath.Join(root, runStateDir, c.WorkflowID)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("failed start left runtime state: %v", err)
	}
	if _, err := os.Stat(s.ContractPath(c.WorkflowID)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("failed start published a contract: %v", err)
	}
}

func TestRestartSafeContinuation(t *testing.T) {
	root := fixtureRoot(t)
	reqs := []Requirement{
		agentReq("graphics-researcher", Required, BlockCompletion),
		agentReq("asset-librarian", Required, BlockCompletion),
	}
	_, c := startedRun(t, root, reqs...)
	reloaded := NewStore(root)
	st, err := reloaded.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Phase != PhaseExecuting || st.RunID == "" {
		t.Errorf("state = %+v", st)
	}
	got, err := reloaded.LoadContract(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if err := got.VerifyFrozen(); err != nil {
		t.Errorf("round-tripped contract: %v", err)
	}
	if rem := st.Remaining(got); len(rem) != 2 {
		t.Errorf("remaining = %v", rem)
	}
	ids, err := reloaded.List()
	if err != nil || len(ids) != 1 || ids[0] != c.WorkflowID {
		t.Errorf("List() = %v, %v", ids, err)
	}
}

func TestRunIdentityReplacementBreaksVerification(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := s.LoadRunIdentity(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	identity.StartingBranch = "silently-replaced"
	if err := WriteJSON(s.IdentityPath(c.WorkflowID, st.RunID), identity); err != nil {
		t.Fatal(err)
	}
	verification := s.VerifyLedger(c.WorkflowID)
	if verification.Valid || !strings.Contains(strings.Join(verification.Problems, "; "), "identity hash") {
		t.Errorf("replaced identity verification = %+v", verification)
	}
}

func TestAppendRejectsReplacedIdentityBeforeLedgerOrStateMutation(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := s.LoadRunIdentity(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	event := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	ledgerPath := s.EventsPath(c.WorkflowID, st.RunID)
	ledgerBefore, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	stateBefore, err := os.ReadFile(s.StatePath(c.WorkflowID))
	if err != nil {
		t.Fatal(err)
	}

	identity.WorkflowID = "WF-9999"
	identity.RunID = "run-ffffffffffffffffffffffffffffffff"
	identity.ContractHash = strings.Repeat("f", 64)
	if problems := identity.validate(); len(problems) != 0 {
		t.Fatalf("replacement identity must remain structurally valid: %v", problems)
	}
	event.WorkflowID = identity.WorkflowID
	event.RunID = identity.RunID
	event.ContractHash = identity.ContractHash
	if err := WriteJSON(s.IdentityPath(c.WorkflowID, st.RunID), identity); err != nil {
		t.Fatal(err)
	}

	err = s.appendEvent(c.WorkflowID, event)
	if err == nil {
		t.Fatal("append accepted a replaced immutable identity")
	}
	if got := err.Error(); !strings.Contains(got, "identity hash") ||
		!strings.Contains(got, "run identity does not match durable state") {
		t.Fatalf("append error = %v, want complete identity/state verification", err)
	}
	ledgerAfter, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatal(err)
	}
	stateAfter, err := os.ReadFile(s.StatePath(c.WorkflowID))
	if err != nil {
		t.Fatal(err)
	}
	if string(ledgerAfter) != string(ledgerBefore) {
		t.Error("identity rejection modified ledger bytes")
	}
	if string(stateAfter) != string(stateBefore) {
		t.Error("identity rejection modified durable tail state bytes")
	}
}

func TestAppendRejectsCrossBoundaryAndHistoricalEvents(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*RunEvent)
		want   string
	}{
		{name: "another workflow", mutate: func(e *RunEvent) { e.WorkflowID = "WF-9999" }, want: "workflow_id"},
		{name: "another run", mutate: func(e *RunEvent) { e.RunID = "run-ffffffffffffffffffffffffffffffff" }, want: "run_id"},
		{name: "another contract", mutate: func(e *RunEvent) { e.ContractHash = strings.Repeat("f", 64) }, want: "contract_hash"},
		{name: "another repository", mutate: func(e *RunEvent) { e.RepositoryIdentity = "remote:https://example.invalid/other" }, want: "repository_identity"},
		{name: "before run start", mutate: func(e *RunEvent) { e.ObservedAt = "2026-08-01T23:59:59Z" }, want: "predates"},
		{name: "unsupported schema", mutate: func(e *RunEvent) { e.SchemaVersion = "99" }, want: "unsupported event schema"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := fixtureRoot(t)
			s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
			event := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
				SourceProviderHook, TrustCollectorObserved)
			tt.mutate(&event)
			if err := s.appendEvent(c.WorkflowID, event); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("append error = %v, want %q", err, tt.want)
			}
			rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
			if err != nil {
				t.Fatal(err)
			}
			if rec.Verdict == Complete {
				t.Error("rejected cross-boundary event satisfied the run")
			}
		})
	}
}

func TestProvenanceEligibilityIsSeparateFromChainValidity(t *testing.T) {
	agent := agentReq("graphics-researcher", Required, BlockCompletion)
	skill := Requirement{ResourceID: "dma-tracer", ResourceType: "skill", Necessity: Required,
		Mode: ModeExplicitSkillCall, EvidenceRule: "ledger skill event", Policy: BlockCompletion}
	backend := Requirement{ResourceID: "go test ./...", ResourceType: "backend", Necessity: Required,
		Mode: ModeDeterministicBackend, EvidenceRule: "ledger backend event", Policy: BlockCompletion}
	tests := []struct {
		name        string
		requirement Requirement
		kind        EventKind
		source      SourceKind
		trust       TrustBasis
		mutate      func(*RunEvent)
	}{
		{name: "manual-import agent", requirement: agent, kind: EventAgentStarted,
			source: SourceManualImport, trust: TrustSelfReported},
		{name: "manual-import skill", requirement: skill, kind: EventSkillSelected,
			source: SourceManualImport, trust: TrustSelfReported},
		{name: "unknown-provenance backend", requirement: backend, kind: EventBackendFinished,
			source: SourceUnknown, trust: TrustUnsupported, mutate: func(e *RunEvent) { e.ExitStatus = ptr(0) }},
		{name: "artifact event with agent selector", requirement: agent, kind: EventOutputObserved,
			source: SourceManualImport, trust: TrustSelfReported},
		{name: "artifact event with backend selector", requirement: backend, kind: EventOutputObserved,
			source: SourceManualImport, trust: TrustSelfReported},
		{name: "self-reported provider hook", requirement: agent, kind: EventAgentStarted,
			source: SourceProviderHook, trust: TrustSelfReported},
		{name: "selector alone on tool start", requirement: backend, kind: EventToolStarted,
			source: SourceProviderHook, trust: TrustCollectorObserved},
		{name: "provider event missing session identity", requirement: agent, kind: EventAgentStarted,
			source: SourceProviderHook, trust: TrustCollectorObserved,
			mutate: func(e *RunEvent) { e.SessionID = "" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := fixtureRoot(t)
			s, c := startedRun(t, root, tt.requirement)
			event := nextEvent(t, s, c.WorkflowID, tt.kind, tt.requirement.ResourceID, tt.source, tt.trust)
			if tt.mutate != nil {
				tt.mutate(&event)
			}
			appendEvent(t, s, c.WorkflowID, event)
			if verification := s.VerifyLedger(c.WorkflowID); !verification.Valid {
				t.Fatalf("event must remain structurally valid and hash-chained: %v", verification.Problems)
			}
			rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
			if err != nil {
				t.Fatal(err)
			}
			if rec.Verdict == Complete {
				t.Errorf("valid but ineligible event satisfied %s requirement: %+v",
					tt.requirement.ResourceType, rec)
			}
			if len(rec.Notes) == 0 {
				t.Error("eligibility limitation must remain visible")
			}
		})
	}
}

func TestAgentInvocationEligibilityRequiresStart(t *testing.T) {
	const requiredAgent = "graphics-researcher"
	type lifecycleEvent struct {
		kind     EventKind
		selector string
		turnID   string
	}
	tests := []struct {
		name                  string
		events                []lifecycleEvent
		wantVerdict           Verdict
		wantOutcome           Outcome
		wantAgentCalls        int
		wantMatchingCalls     int
		wantFinishLimitations int
	}{
		{
			name: "start-only event satisfies invocation",
			events: []lifecycleEvent{
				{kind: EventAgentStarted, selector: requiredAgent, turnID: "turn-1"},
			},
			wantVerdict: Complete, wantOutcome: OutcomeSatisfied,
			wantAgentCalls: 1, wantMatchingCalls: 1,
		},
		{
			name: "matching start and finish satisfy through start",
			events: []lifecycleEvent{
				{kind: EventAgentStarted, selector: requiredAgent, turnID: "turn-1"},
				{kind: EventAgentFinished, selector: requiredAgent, turnID: "turn-1"},
			},
			wantVerdict: Complete, wantOutcome: OutcomeSatisfied,
			wantAgentCalls: 1, wantMatchingCalls: 1, wantFinishLimitations: 1,
		},
		{
			name: "finish-only event cannot satisfy invocation",
			events: []lifecycleEvent{
				{kind: EventAgentFinished, selector: requiredAgent, turnID: "turn-1"},
			},
			wantVerdict: Partial, wantOutcome: OutcomeUnsatisfied,
			wantAgentCalls: 0, wantMatchingCalls: 0, wantFinishLimitations: 1,
		},
		{
			name: "unmatched finish cannot satisfy invocation",
			events: []lifecycleEvent{
				{kind: EventAgentStarted, selector: "asset-librarian", turnID: "turn-1"},
				{kind: EventAgentFinished, selector: requiredAgent, turnID: "turn-2"},
			},
			wantVerdict: Partial, wantOutcome: OutcomeUnsatisfied,
			wantAgentCalls: 1, wantMatchingCalls: 0, wantFinishLimitations: 1,
		},
		{
			name: "repeated finishes create no additional invocation credit",
			events: []lifecycleEvent{
				{kind: EventAgentStarted, selector: requiredAgent, turnID: "turn-1"},
				{kind: EventAgentFinished, selector: requiredAgent, turnID: "turn-1"},
				{kind: EventAgentFinished, selector: requiredAgent, turnID: "turn-2"},
			},
			wantVerdict: Complete, wantOutcome: OutcomeSatisfied,
			wantAgentCalls: 1, wantMatchingCalls: 1, wantFinishLimitations: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := fixtureRoot(t)
			s, c := startedRun(t, root, agentReq(requiredAgent, Required, BlockCompletion))
			for _, spec := range tt.events {
				event := nextEvent(t, s, c.WorkflowID, spec.kind, spec.selector,
					SourceProviderHook, TrustCollectorObserved)
				event.TurnID = spec.turnID
				appendEvent(t, s, c.WorkflowID, event)
			}

			verification := s.VerifyLedger(c.WorkflowID)
			if !verification.Valid {
				t.Fatalf("lifecycle events must remain structurally valid: %v", verification.Problems)
			}
			evidence := evidenceFromLedger(verification)
			var agentCalls, matchingCalls int
			for _, observation := range evidence.Observations {
				if observation.Kind != ObsAgentCall {
					continue
				}
				agentCalls++
				if observation.Selector == requiredAgent {
					matchingCalls++
				}
			}
			if agentCalls != tt.wantAgentCalls || matchingCalls != tt.wantMatchingCalls {
				t.Errorf("agent observations = %d total, %d matching; want %d total, %d matching",
					agentCalls, matchingCalls, tt.wantAgentCalls, tt.wantMatchingCalls)
			}
			if got := len(evidence.Limitations); got != tt.wantFinishLimitations {
				t.Errorf("limitations = %d, want %d: %v", got, tt.wantFinishLimitations, evidence.Limitations)
			}
			for _, limitation := range evidence.Limitations {
				if !strings.Contains(limitation, "agent_finished is lifecycle evidence only") {
					t.Errorf("finish limitation = %q", limitation)
				}
			}

			rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
			if err != nil {
				t.Fatal(err)
			}
			if rec.Verdict != tt.wantVerdict {
				t.Fatalf("verdict = %q, want %q: %+v", rec.Verdict, tt.wantVerdict, rec)
			}
			if len(rec.Results) != 1 || rec.Results[0].Outcome != tt.wantOutcome {
				t.Fatalf("results = %+v, want outcome %q", rec.Results, tt.wantOutcome)
			}
		})
	}
}

func TestBranchAndHeadMayAdvanceWithinSameImmutableRun(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	event := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	event.Branch = "codex/work-advanced"
	event.Head = strings.Repeat("b", 40)
	appendEvent(t, s, c.WorkflowID, event)
	rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != Complete {
		t.Errorf("legitimate branch/HEAD advance invalidated run: %+v", rec)
	}
}

func TestBackendEventRequiresEligibleCompletionAndExitStatus(t *testing.T) {
	root := fixtureRoot(t)
	req := Requirement{ResourceID: "go test ./...", ResourceType: "backend", Necessity: Required,
		Mode: ModeDeterministicBackend, EvidenceRule: "ledger backend completion", Policy: BlockCompletion}
	s, c := startedRun(t, root, req)
	event := nextEvent(t, s, c.WorkflowID, EventBackendFinished, req.ResourceID,
		SourceDeterministicBackend, TrustBackendExitStatus)
	event.ExitStatus = ptr(0)
	appendEvent(t, s, c.WorkflowID, event)
	rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != Complete {
		t.Errorf("verdict = %q, want complete", rec.Verdict)
	}
}

func TestProviderToolFinishedEligibility(t *testing.T) {
	const selector = "go test ./..."
	tests := []struct {
		name        string
		mutate      func(*RunEvent)
		wantOutcome Outcome
		wantVerdict Verdict
		wantReason  string
		wantNote    string
	}{
		{
			name: "tool_finished without tool_use_id cannot satisfy a backend",
			mutate: func(e *RunEvent) {
				e.ToolUseID = ""
			},
			wantOutcome: OutcomeUnsatisfied, wantVerdict: Partial,
			wantReason: "did not run", wantNote: "lacks a tool-use binding",
		},
		{
			name: "tool_finished with empty tool_use_id cannot satisfy a backend",
			mutate: func(e *RunEvent) {
				e.ToolUseID = "   "
			},
			wantOutcome: OutcomeUnsatisfied, wantVerdict: Partial,
			wantReason: "did not run", wantNote: "lacks a tool-use binding",
		},
		{
			name: "bound tool_finished with null exit status is unverifiable, never pass",
			mutate: func(e *RunEvent) {
				e.ToolUseID = "provider-tool-use-1"
				e.ExitStatus = nil
			},
			wantOutcome: OutcomeUnverifiable, wantVerdict: Unverifiable,
			wantReason: "no exit status was captured",
		},
		{
			name: "selector session and turn alone are insufficient",
			mutate: func(e *RunEvent) {
				e.ToolUseID = ""
				e.ExitStatus = nil
			},
			wantOutcome: OutcomeUnsatisfied, wantVerdict: Partial,
			wantReason: "did not run", wantNote: "lacks a tool-use binding",
		},
		{
			name: "properly bound provider tool completion satisfies a backend",
			mutate: func(e *RunEvent) {
				e.ToolUseID = "provider-tool-use-1"
			},
			wantOutcome: OutcomeSatisfied, wantVerdict: Complete,
			wantReason: "backend exited 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := fixtureRoot(t)
			s, c := startedRun(t, root, backendReq(selector))
			event := nextEvent(t, s, c.WorkflowID, EventToolFinished, selector,
				SourceProviderHook, TrustCollectorObserved)
			event.ExitStatus = ptr(0)
			tt.mutate(&event)
			appendEvent(t, s, c.WorkflowID, event)
			if verification := s.VerifyLedger(c.WorkflowID); !verification.Valid {
				t.Fatalf("provider event must remain structurally valid: %v", verification.Problems)
			}
			rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
			if err != nil {
				t.Fatal(err)
			}
			if rec.Verdict != tt.wantVerdict {
				t.Fatalf("verdict = %q, want %q: %+v", rec.Verdict, tt.wantVerdict, rec)
			}
			if len(rec.Results) != 1 || rec.Results[0].Outcome != tt.wantOutcome {
				t.Fatalf("results = %+v, want outcome %q", rec.Results, tt.wantOutcome)
			}
			if !strings.Contains(rec.Results[0].Reason, tt.wantReason) {
				t.Errorf("reason = %q, want substring %q", rec.Results[0].Reason, tt.wantReason)
			}
			if tt.wantNote != "" && !strings.Contains(strings.Join(rec.Notes, "; "), tt.wantNote) {
				t.Errorf("notes = %v, want substring %q", rec.Notes, tt.wantNote)
			}
		})
	}
}

func TestLedgerSequenceSurvivesConversionAndLatestBackendCompletionGoverns(t *testing.T) {
	const selector = "go test ./..."
	root := fixtureRoot(t)
	s, c := startedRun(t, root, backendReq(selector))
	pass := nextEvent(t, s, c.WorkflowID, EventBackendFinished, selector,
		SourceDeterministicBackend, TrustBackendExitStatus)
	pass.ExitStatus = ptr(0)
	appendEvent(t, s, c.WorkflowID, pass)
	fail := nextEvent(t, s, c.WorkflowID, EventBackendFinished, selector,
		SourceDeterministicBackend, TrustBackendExitStatus)
	fail.ExitStatus = ptr(2)
	appendEvent(t, s, c.WorkflowID, fail)

	rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != Failed || len(rec.Results) != 1 || rec.Results[0].Outcome != OutcomeUnsatisfied {
		t.Fatalf("latest failing completion did not govern: %+v", rec)
	}
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Evidence.Observations) != 2 || st.Evidence.Observations[0].Sequence != 1 ||
		st.Evidence.Observations[1].Sequence != 2 {
		t.Fatalf("converted observation sequences = %+v, want [1 2]", st.Evidence.Observations)
	}
	if len(rec.Results[0].Evidence) != 1 ||
		rec.Results[0].Evidence[0] != st.Evidence.Observations[1].EvidenceRef {
		t.Errorf("result evidence = %v, want final completion %q",
			rec.Results[0].Evidence, st.Evidence.Observations[1].EvidenceRef)
	}
}

func TestWriterRejectsDuplicateEventIDAndSequenceWithoutMutation(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	first := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	appendEvent(t, s, c.WorkflowID, first)
	st, _ := s.LoadState(c.WorkflowID)
	path := s.EventsPath(c.WorkflowID, st.RunID)
	before, _ := os.ReadFile(path)
	beforeSequence, beforeHash := st.LedgerSequence, st.LedgerHash

	duplicateID := nextEvent(t, s, c.WorkflowID, EventAgentFinished, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	duplicateID.EventID = first.EventID
	if err := s.appendEvent(c.WorkflowID, duplicateID); err == nil || !strings.Contains(err.Error(), "duplicate event_id") {
		t.Fatalf("duplicate event id error = %v", err)
	}
	duplicateSequence := nextEvent(t, s, c.WorkflowID, EventAgentFinished, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	duplicateSequence.Sequence = 1
	if err := s.appendEvent(c.WorkflowID, duplicateSequence); err == nil || !strings.Contains(err.Error(), "sequence") {
		t.Fatalf("duplicate sequence error = %v", err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Error("rejected append modified the ledger")
	}
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if st.LedgerSequence != beforeSequence || st.LedgerHash != beforeHash {
		t.Errorf("rejected append advanced durable sequence state: sequence=%d hash=%s", st.LedgerSequence, st.LedgerHash)
	}
}

func seededLedger(t *testing.T) (*Store, *Contract, string, [][]byte) {
	t.Helper()
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	for i := 0; i < 3; i++ {
		event := nextEvent(t, s, c.WorkflowID, EventAgentStarted,
			fmt.Sprintf("graphics-researcher-%d", i), SourceProviderHook, TrustCollectorObserved)
		appendEvent(t, s, c.WorkflowID, event)
	}
	st, _ := s.LoadState(c.WorkflowID)
	path := s.EventsPath(c.WorkflowID, st.RunID)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	trimmed := strings.TrimSuffix(string(b), "\n")
	parts := strings.Split(trimmed, "\n")
	lines := make([][]byte, len(parts))
	for i := range parts {
		lines[i] = []byte(parts[i])
	}
	return s, c, path, lines
}

func writeLedgerLines(t *testing.T, path string, lines [][]byte, finalNewline bool) {
	t.Helper()
	body := bytesJoin(lines, []byte{'\n'})
	if finalNewline {
		body = append(body, '\n')
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
}

func bytesJoin(parts [][]byte, separator []byte) []byte {
	var out []byte
	for i, part := range parts {
		if i != 0 {
			out = append(out, separator...)
		}
		out = append(out, part...)
	}
	return out
}

func mutateLedgerEvent(t *testing.T, line []byte, mutate func(*RunEvent), rehash bool) []byte {
	t.Helper()
	var event RunEvent
	if err := json.Unmarshal(line, &event); err != nil {
		t.Fatal(err)
	}
	mutate(&event)
	if rehash {
		hash, err := eventHash(event)
		if err != nil {
			t.Fatal(err)
		}
		event.EventHash = hash
	}
	b, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestVerifierDetectsTamperedLedgerHistories(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string, [][]byte)
		want   string
	}{
		{name: "modified event", mutate: func(t *testing.T, path string, lines [][]byte) {
			lines[0] = mutateLedgerEvent(t, lines[0], func(e *RunEvent) { e.Selector = "changed" }, false)
			writeLedgerLines(t, path, lines, true)
		}, want: "canonical hash"},
		{name: "deleted middle event", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, [][]byte{lines[0], lines[2]}, true)
		}, want: "sequence"},
		{name: "deleted final event", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, lines[:2], true)
		}, want: "ledger tail"},
		{name: "reordered events", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, [][]byte{lines[1], lines[0], lines[2]}, true)
		}, want: "sequence"},
		{name: "inserted event", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, [][]byte{lines[0], lines[0], lines[1], lines[2]}, true)
		}, want: "duplicate event_id"},
		{name: "malformed JSON", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, [][]byte{lines[0], []byte("not-json"), lines[1], lines[2]}, true)
		}, want: "malformed JSON"},
		{name: "truncated final event", mutate: func(t *testing.T, path string, lines [][]byte) {
			writeLedgerLines(t, path, [][]byte{lines[0], lines[1], []byte(`{"schema_version":"1.0"`)}, false)
		}, want: "truncated final record"},
		{name: "unsupported schema version", mutate: func(t *testing.T, path string, lines [][]byte) {
			lines[0] = mutateLedgerEvent(t, lines[0], func(e *RunEvent) { e.SchemaVersion = "99" }, true)
			writeLedgerLines(t, path, lines, true)
		}, want: "unsupported event schema"},
		{name: "broken previous hash", mutate: func(t *testing.T, path string, lines [][]byte) {
			lines[1] = mutateLedgerEvent(t, lines[1], func(e *RunEvent) { e.PreviousHash = strings.Repeat("f", 64) }, true)
			writeLedgerLines(t, path, lines, true)
		}, want: "previous_hash"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, c, path, lines := seededLedger(t)
			tt.mutate(t, path, lines)
			verification := s.VerifyLedger(c.WorkflowID)
			if verification.Valid {
				t.Fatal("tampered ledger verified")
			}
			if !strings.Contains(strings.Join(verification.Problems, "; "), tt.want) {
				t.Errorf("problems = %v, want %q", verification.Problems, tt.want)
			}
			if len(verification.Events) != 0 {
				t.Error("invalid ledger events must not be released as observations")
			}
		})
	}
}

func TestDamagedLedgerClosesAsIncompleteNeverConfirmedAbsence(t *testing.T) {
	s, c, path, lines := seededLedger(t)
	writeLedgerLines(t, path, [][]byte{lines[0], []byte("{"), lines[2]}, false)
	rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != Unverifiable {
		t.Errorf("verdict = %q, want unverifiable", rec.Verdict)
	}
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Phase != PhaseClosed || len(st.Evidence.IntegrityErrors) == 0 {
		t.Errorf("damaged ledger state = %+v", st)
	}
}

func TestTerminalClosureRejectsEventsAndSecondCloseWithoutMutation(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	event := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	appendEvent(t, s, c.WorkflowID, event)
	if _, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z"); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(s.StatePath(c.WorkflowID))
	if err != nil {
		t.Fatal(err)
	}
	afterClose := event
	afterClose.Sequence++
	newID, _ := newEventID()
	afterClose.EventID = newID
	afterClose.ObservedAt = "2026-08-02T00:01:01Z"
	if err := s.appendEvent(c.WorkflowID, afterClose); err == nil || !strings.Contains(err.Error(), "terminally closed") {
		t.Errorf("post-close append error = %v", err)
	}
	if _, err := s.Close(c.WorkflowID, "2026-08-02T00:02:00Z"); err == nil || !strings.Contains(err.Error(), "already terminally closed") {
		t.Errorf("second close error = %v", err)
	}
	after, err := os.ReadFile(s.StatePath(c.WorkflowID))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("a rejected append or second close mutated terminal state")
	}
}

func TestCloseDetectsContractSwappedAfterStart(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root, agentReq("graphics-researcher", Required, BlockCompletion))
	other := frozen(t, agentReq("asset-librarian", Required, BlockCompletion))
	if err := WriteJSON(s.ContractPath(c.WorkflowID), other); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z"); err == nil {
		t.Error("closing must fail when the contract changed since start")
	}
}

func TestCloseRecordsComputedVerdictAndRemaining(t *testing.T) {
	root := fixtureRoot(t)
	s, c := startedRun(t, root,
		agentReq("graphics-researcher", Required, BlockCompletion),
		agentReq("asset-librarian", Required, BlockCompletion))
	event := nextEvent(t, s, c.WorkflowID, EventAgentStarted, "graphics-researcher",
		SourceProviderHook, TrustCollectorObserved)
	appendEvent(t, s, c.WorkflowID, event)
	rec, err := s.Close(c.WorkflowID, "2026-08-02T00:01:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != Partial {
		t.Errorf("verdict = %q, want partial", rec.Verdict)
	}
	st, err := s.LoadState(c.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Phase != PhaseClosed || st.Reconciliation == nil {
		t.Fatalf("terminal state = %+v", st)
	}
	if remaining := st.Remaining(c); len(remaining) != 1 || remaining[0] != "asset-librarian" {
		t.Errorf("remaining = %v, want [asset-librarian]", remaining)
	}
}

func TestN9ReceiptDisagreesWithVerifier(t *testing.T) {
	tests := []struct {
		name     string
		receipt  string
		computed Verdict
		wantErr  bool
	}{
		{name: "receipt agrees", receipt: "## Result\n\nVerdict: complete\n", computed: Complete},
		{name: "receipt claims complete verifier says partial", receipt: "**Verdict:** complete\n", computed: Partial, wantErr: true},
		{name: "receipt claims complete verifier says failed", receipt: "- verdict = `complete`\n", computed: Failed, wantErr: true},
		{name: "receipt states no verdict", receipt: "Everything went fine.\n", computed: Complete, wantErr: true},
		{name: "bulleted bolded form", receipt: "- **Verdict**: `partial`\n", computed: Partial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReceipt(tt.receipt, tt.computed)
			if tt.wantErr {
				if err == nil || !errors.Is(err, ErrReceiptDisagrees) {
					t.Fatalf("want ErrReceiptDisagrees, got %v", err)
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
		t.Errorf("missing receipt: %v", err)
	}
	if err := os.WriteFile(s.ReceiptPath(c.WorkflowID), []byte("Verdict: complete\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.ValidateReceiptFile(c.WorkflowID, Failed); err == nil {
		t.Error("disagreeing receipt must fail")
	}
}
