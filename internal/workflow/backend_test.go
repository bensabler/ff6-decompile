package workflow

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var _ func(*Store, context.Context, string, string, io.Writer, io.Writer) (deterministicBackendResult, error) = (*Store).executeDeterministicBackend

type deterministicBackendRunnerFunc func(context.Context, string, string, io.Writer, io.Writer) (deterministicBackendProcess, error)

func (f deterministicBackendRunnerFunc) run(ctx context.Context, root, command string,
	stdout, stderr io.Writer) (deterministicBackendProcess, error) {
	return f(ctx, root, command, stdout, stderr)
}

type backendInvocation struct {
	root    string
	command string
}

type recordingBackendRunner struct {
	invocations []backendInvocation
	process     deterministicBackendProcess
	err         error
}

func (r *recordingBackendRunner) run(_ context.Context, root, command string,
	_, _ io.Writer) (deterministicBackendProcess, error) {
	r.invocations = append(r.invocations, backendInvocation{root: root, command: command})
	return r.process, r.err
}

type backendExecution struct {
	result deterministicBackendResult
	err    error
}

func executeBackendAsync(s *Store, workflowID, requirementID string) <-chan backendExecution {
	done := make(chan backendExecution, 1)
	go func() {
		result, err := s.executeDeterministicBackend(
			context.Background(), workflowID, requirementID, io.Discard, io.Discard)
		done <- backendExecution{result: result, err: err}
	}()
	return done
}

func waitForBackendInvocation(t *testing.T, started <-chan backendInvocation) backendInvocation {
	t.Helper()
	select {
	case invocation := <-started:
		return invocation
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for deterministic backend runner")
		return backendInvocation{}
	}
}

func waitForBackendExecution(t *testing.T, done <-chan backendExecution) backendExecution {
	t.Helper()
	select {
	case execution := <-done:
		return execution
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for deterministic backend execution")
		return backendExecution{}
	}
}

func futureBackendTime() string {
	return time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano)
}

func requireRealBackendShell(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(deterministicBackendShell); err != nil {
		t.Skipf("real deterministic backend shell %s is unavailable: %v", deterministicBackendShell, err)
	}
}

func verifiedBackendEvents(t *testing.T, s *Store, workflowID string) []RunEvent {
	t.Helper()
	verification := s.VerifyLedger(workflowID)
	if !verification.Valid {
		t.Fatalf("ledger verification failed: %v", verification.Problems)
	}
	return verification.Events
}

func assertBackendFinishedEvent(t *testing.T, event RunEvent, sequence uint64,
	selector string, exitStatus int) {
	t.Helper()
	if event.Sequence != sequence || event.EventKind != EventBackendFinished ||
		event.SourceKind != SourceDeterministicBackend || event.TrustBasis != TrustBackendExitStatus ||
		event.Provider != deterministicBackendProvider || event.CollectorID != deterministicBackendCollectorID ||
		event.Selector != selector || event.ExitStatus == nil || *event.ExitStatus != exitStatus {
		t.Fatalf("backend event = %+v, want sequence %d selector %q exit status %d",
			event, sequence, selector, exitStatus)
	}
	if event.SessionID != "" || event.TurnID != "" || event.ToolUseID != "" || event.EvidenceRef != "" {
		t.Errorf("backend event manufactured provider provenance: %+v", event)
	}
	if eligibility := evaluateEventEligibility(event); !eligibility.Eligible || eligibility.ObservationKind != ObsBackendRun {
		t.Errorf("backend event eligibility = %+v", eligibility)
	}
}

func readBackendDurableFiles(t *testing.T, s *Store, workflowID string) (state, ledger []byte) {
	t.Helper()
	st, err := s.LoadState(workflowID)
	if err != nil {
		t.Fatal(err)
	}
	state, err = os.ReadFile(s.StatePath(workflowID))
	if err != nil {
		t.Fatal(err)
	}
	ledger, err = os.ReadFile(s.EventsPath(workflowID, st.RunID))
	if err != nil {
		t.Fatal(err)
	}
	return state, ledger
}

func rebindContractForBackendTest(t *testing.T, s *Store, contract *Contract) {
	t.Helper()
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := s.LoadRunIdentity(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	identity.ContractHash = contract.FrozenHash
	identityHash, err := runIdentityHash(*identity)
	if err != nil {
		t.Fatal(err)
	}
	st.ContractHash = contract.FrozenHash
	st.IdentityHash = identityHash
	if err := WriteJSON(s.IdentityPath(contract.WorkflowID, st.RunID), identity); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(s.StatePath(contract.WorkflowID), st); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(s.ContractPath(contract.WorkflowID), contract); err != nil {
		t.Fatal(err)
	}
}

func TestDeterministicBackendAPITrustBoundary(t *testing.T) {
	storeType := reflect.TypeOf((*Store)(nil))
	runEventType := reflect.TypeOf(RunEvent{})
	for i := 0; i < storeType.NumMethod(); i++ {
		method := storeType.Method(i)
		name := strings.ToLower(method.Name)
		mutatesEvidence := strings.Contains(name, "append") || strings.Contains(name, "record") ||
			strings.Contains(name, "import") || strings.Contains(name, "assert")
		if strings.Contains(name, "backend") ||
			(mutatesEvidence && strings.Contains(name, "event")) {
			t.Errorf("Store exposes forbidden production assertion method %s", method.Name)
		}
		for argument := 1; argument < method.Type.NumIn(); argument++ {
			argumentType := method.Type.In(argument)
			if argumentType == runEventType ||
				(argumentType.Kind() == reflect.Pointer && argumentType.Elem() == runEventType) {
				t.Errorf("Store method %s publicly accepts a RunEvent", method.Name)
			}
		}
	}
}

func TestDeterministicBackendSelectsExactFrozenRequirement(t *testing.T) {
	root := fixtureRoot(t)
	const command = `printf '%s\n' "the exact frozen command"`
	requirement := backendReq(command)
	// Eligibility is determined only by execution mode, not these descriptive fields.
	requirement.ResourceType = "documentation"
	requirement.EvidenceRule = "descriptive text that does not imply a backend"
	s, contract := startedRun(t, root, requirement)
	runner := &recordingBackendRunner{
		process: deterministicBackendProcess{started: true, terminated: true, exitStatus: 0},
	}
	s.backendRunner = runner

	result, err := s.executeDeterministicBackend(
		context.Background(), contract.WorkflowID, command, io.Discard, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if result.ExitStatus != 0 {
		t.Fatalf("exit status = %d, want 0", result.ExitStatus)
	}
	if len(runner.invocations) != 1 {
		t.Fatalf("runner invocation count = %d, want 1", len(runner.invocations))
	}
	canonicalRoot, err := canonicalPath(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := runner.invocations[0]; got.root != canonicalRoot || got.command != command {
		t.Fatalf("runner invocation = %+v, want root %q and exact command %q", got, canonicalRoot, command)
	}
	events := verifiedBackendEvents(t, s, contract.WorkflowID)
	if len(events) != 1 {
		t.Fatalf("ledger event count = %d, want 1", len(events))
	}
	assertBackendFinishedEvent(t, events[0], 1, command, 0)
}

func TestDeterministicBackendRejectsInvalidPreconditionsBeforeLaunch(t *testing.T) {
	const command = "exit 0"
	tests := []struct {
		name        string
		requirement Requirement
		selector    string
		mutate      func(*testing.T, *Store, *Contract)
		wantError   string
	}{
		{
			name: "unknown requirement", selector: "exit 1", wantError: "no requirement with exact resource_id",
		},
		{
			name: "non-backend requirement",
			requirement: Requirement{ResourceID: command, ResourceType: "backend", Necessity: Required,
				Mode: ModeExplicitAgentCall, EvidenceRule: "command-like prose", Policy: BlockCompletion},
			selector: command, wantError: "not \"deterministic_backend\"",
		},
		{
			name: "ambiguous requirement", selector: command, wantError: "ambiguous",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				stored, err := s.LoadContract(contract.WorkflowID)
				if err != nil {
					t.Fatal(err)
				}
				stored.Requirements = append(stored.Requirements, stored.Requirements[0])
				stored.FrozenHash, err = stored.Hash()
				if err != nil {
					t.Fatal(err)
				}
				rebindContractForBackendTest(t, s, stored)
			},
		},
		{
			name: "unfrozen contract", selector: command, wantError: "contract is not frozen",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				contract.State = Draft
				if err := WriteJSON(s.ContractPath(contract.WorkflowID), contract); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "contract hash mismatch", selector: command, wantError: "changed since the run started",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				replacement := frozen(t, backendReq(command), backendReq("exit 9"))
				if err := WriteJSON(s.ContractPath(contract.WorkflowID), replacement); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "closed run", selector: command, wantError: "not executable",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				if _, err := s.Close(contract.WorkflowID, futureBackendTime()); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "non-executing run", selector: command, wantError: "not executable",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				st, err := s.LoadState(contract.WorkflowID)
				if err != nil {
					t.Fatal(err)
				}
				st.Phase = PhasePlanned
				if err := WriteJSON(s.StatePath(contract.WorkflowID), st); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "replaced immutable identity", selector: command, wantError: "immutable run identity hash",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				st, err := s.LoadState(contract.WorkflowID)
				if err != nil {
					t.Fatal(err)
				}
				identity, err := s.LoadRunIdentity(contract.WorkflowID)
				if err != nil {
					t.Fatal(err)
				}
				identity.StartingBranch = "replaced-identity"
				if err := WriteJSON(s.IdentityPath(contract.WorkflowID, st.RunID), identity); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "invalid ledger", selector: command, wantError: "invalid ledger",
			mutate: func(t *testing.T, s *Store, contract *Contract) {
				st, err := s.LoadState(contract.WorkflowID)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(s.EventsPath(contract.WorkflowID, st.RunID), []byte("{not-json}\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := fixtureRoot(t)
			requirement := test.requirement
			if requirement.ResourceID == "" {
				requirement = backendReq(command)
			}
			s, contract := startedRun(t, root, requirement)
			if test.mutate != nil {
				test.mutate(t, s, contract)
			}
			beforeState, beforeLedger := readBackendDurableFiles(t, s, contract.WorkflowID)
			runner := &recordingBackendRunner{
				process: deterministicBackendProcess{started: true, terminated: true, exitStatus: 0},
			}
			s.backendRunner = runner

			_, err := s.executeDeterministicBackend(
				context.Background(), contract.WorkflowID, test.selector, io.Discard, io.Discard)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("execution error = %v, want substring %q", err, test.wantError)
			}
			if len(runner.invocations) != 0 {
				t.Fatalf("invalid precondition launched runner %d times", len(runner.invocations))
			}
			afterState, afterLedger := readBackendDurableFiles(t, s, contract.WorkflowID)
			if !bytes.Equal(afterState, beforeState) || !bytes.Equal(afterLedger, beforeLedger) {
				t.Error("rejected precondition mutated durable state or ledger")
			}
		})
	}
}

func TestDeterministicBackendRealProcessExitStatusAndReconciliation(t *testing.T) {
	requireRealBackendShell(t)
	const (
		stdoutMarker = "BACKEND_STDOUT_MARKER: PASS; exit 0; success"
		stderrMarker = "BACKEND_STDERR_MARKER: PASS; exit 0; success"
	)
	tests := []struct {
		name        string
		command     string
		wantStatus  int
		wantVerdict Verdict
	}{
		{name: "exit zero", command: "exit 0", wantStatus: 0, wantVerdict: Complete},
		{
			name: "exit seven despite misleading output",
			command: "printf '%s%s\\n' 'BACKEND_STDOUT_' 'MARKER: PASS; exit 0; success'; " +
				"printf '%s%s\\n' 'BACKEND_STDERR_' 'MARKER: PASS; exit 0; success' >&2; exit 7",
			wantStatus: 7, wantVerdict: Failed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := fixtureRoot(t)
			s, contract := startedRun(t, root, backendReq(test.command))
			var stdout, stderr bytes.Buffer
			result, err := s.executeDeterministicBackend(
				context.Background(), contract.WorkflowID, test.command, &stdout, &stderr)
			if err != nil {
				t.Fatal(err)
			}
			if result.ExitStatus != test.wantStatus {
				t.Fatalf("real process exit status = %d, want %d", result.ExitStatus, test.wantStatus)
			}
			events := verifiedBackendEvents(t, s, contract.WorkflowID)
			if len(events) != 1 {
				t.Fatalf("ledger event count = %d, want 1", len(events))
			}
			assertBackendFinishedEvent(t, events[0], 1, test.command, test.wantStatus)

			if test.wantStatus == 7 {
				if !strings.Contains(stdout.String(), stdoutMarker) || !strings.Contains(stderr.String(), stderrMarker) {
					t.Fatalf("misleading process output was not captured: stdout=%q stderr=%q", stdout.String(), stderr.String())
				}
				state, ledger := readBackendDurableFiles(t, s, contract.WorkflowID)
				for _, marker := range []string{stdoutMarker, stderrMarker} {
					if bytes.Contains(state, []byte(marker)) || bytes.Contains(ledger, []byte(marker)) {
						t.Errorf("process output marker %q entered durable state or ledger", marker)
					}
				}
			}

			reconciliation, err := s.Close(contract.WorkflowID, futureBackendTime())
			if err != nil {
				t.Fatal(err)
			}
			if reconciliation.Verdict != test.wantVerdict {
				t.Fatalf("reconciliation verdict = %q, want %q: %+v",
					reconciliation.Verdict, test.wantVerdict, reconciliation)
			}
			if test.wantStatus == 7 {
				state, ledger := readBackendDurableFiles(t, s, contract.WorkflowID)
				for _, marker := range []string{stdoutMarker, stderrMarker} {
					if bytes.Contains(state, []byte(marker)) || bytes.Contains(ledger, []byte(marker)) {
						t.Errorf("process output marker %q entered reconciled state or ledger", marker)
					}
				}
			}
		})
	}
}

func TestDeterministicBackendCanceledContextLaunchFailureRecordsNoEvent(t *testing.T) {
	requireRealBackendShell(t)
	root := fixtureRoot(t)
	const command = "exit 0"
	s, contract := startedRun(t, root, backendReq(command))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := s.executeDeterministicBackend(ctx, contract.WorkflowID, command, io.Discard, io.Discard); err == nil {
		t.Fatal("canceled context unexpectedly launched and credited a backend")
	}
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := os.ReadFile(s.EventsPath(contract.WorkflowID, st.RunID))
	if err != nil {
		t.Fatal(err)
	}
	if st.LedgerSequence != 0 || st.LedgerHash != genesisEventHash || len(ledger) != 0 {
		t.Fatalf("launch failure mutated ledger state: sequence=%d hash=%s ledger=%q",
			st.LedgerSequence, st.LedgerHash, ledger)
	}
}

func TestDeterministicBackendAllocatesSequenceAfterBlockedRunner(t *testing.T) {
	root := fixtureRoot(t)
	const command = "exit 0"
	s, contract := startedRun(t, root, backendReq(command))
	started := make(chan backendInvocation, 1)
	release := make(chan struct{})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	s.backendRunner = deterministicBackendRunnerFunc(func(_ context.Context, root, command string,
		_, _ io.Writer) (deterministicBackendProcess, error) {
		started <- backendInvocation{root: root, command: command}
		<-release
		return deterministicBackendProcess{started: true, terminated: true, exitStatus: 0}, nil
	})

	done := executeBackendAsync(s, contract.WorkflowID, command)
	waitForBackendInvocation(t, started)
	unrelated := nextEvent(t, s, contract.WorkflowID, EventAgentStarted, "unrelated-agent",
		SourceProviderHook, TrustCollectorObserved)
	appendEvent(t, s, contract.WorkflowID, unrelated)
	close(release)
	execution := waitForBackendExecution(t, done)
	if execution.err != nil || execution.result.ExitStatus != 0 {
		t.Fatalf("backend execution = %+v", execution)
	}

	events := verifiedBackendEvents(t, s, contract.WorkflowID)
	if len(events) != 2 {
		t.Fatalf("ledger event count = %d, want 2", len(events))
	}
	if events[0].Sequence != 1 || events[0].EventKind != EventAgentStarted {
		t.Fatalf("unrelated event = %+v, want sequence 1", events[0])
	}
	assertBackendFinishedEvent(t, events[1], 2, command, 0)
}

func TestDeterministicBackendCloseWhileBlockedRecordsNoCredit(t *testing.T) {
	root := fixtureRoot(t)
	const command = "exit 0"
	s, contract := startedRun(t, root, backendReq(command))
	started := make(chan backendInvocation, 1)
	release := make(chan struct{})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	s.backendRunner = deterministicBackendRunnerFunc(func(_ context.Context, root, command string,
		_, _ io.Writer) (deterministicBackendProcess, error) {
		started <- backendInvocation{root: root, command: command}
		<-release
		return deterministicBackendProcess{started: true, terminated: true, exitStatus: 0}, nil
	})

	done := executeBackendAsync(s, contract.WorkflowID, command)
	waitForBackendInvocation(t, started)
	if _, err := s.Close(contract.WorkflowID, futureBackendTime()); err != nil {
		t.Fatal(err)
	}
	close(release)
	execution := waitForBackendExecution(t, done)
	if execution.err == nil || !strings.Contains(execution.err.Error(), "no trustworthy ledger credit") {
		t.Fatalf("post-close backend error = %v", execution.err)
	}
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := os.ReadFile(s.EventsPath(contract.WorkflowID, st.RunID))
	if err != nil {
		t.Fatal(err)
	}
	if st.Phase != PhaseClosed || st.LedgerSequence != 0 || len(ledger) != 0 {
		t.Fatalf("closed run received backend credit: state=%+v ledger=%q", st, ledger)
	}
}

func TestDeterministicBackendContractReplacementWhileBlockedRecordsNoCredit(t *testing.T) {
	root := fixtureRoot(t)
	const command = "exit 0"
	s, contract := startedRun(t, root, backendReq(command))
	started := make(chan backendInvocation, 1)
	release := make(chan struct{})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	s.backendRunner = deterministicBackendRunnerFunc(func(_ context.Context, root, command string,
		_, _ io.Writer) (deterministicBackendProcess, error) {
		started <- backendInvocation{root: root, command: command}
		<-release
		return deterministicBackendProcess{started: true, terminated: true, exitStatus: 0}, nil
	})

	done := executeBackendAsync(s, contract.WorkflowID, command)
	waitForBackendInvocation(t, started)
	replacement := frozen(t, backendReq(command), backendReq("exit 9"))
	if err := WriteJSON(s.ContractPath(contract.WorkflowID), replacement); err != nil {
		t.Fatal(err)
	}
	close(release)
	execution := waitForBackendExecution(t, done)
	if execution.err == nil || !strings.Contains(execution.err.Error(), "no trustworthy ledger credit") {
		t.Fatalf("post-replacement backend error = %v", execution.err)
	}
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := os.ReadFile(s.EventsPath(contract.WorkflowID, st.RunID))
	if err != nil {
		t.Fatal(err)
	}
	if st.LedgerSequence != 0 || st.LedgerHash != genesisEventHash || len(ledger) != 0 {
		t.Fatalf("replaced contract received backend credit: state=%+v ledger=%q", st, ledger)
	}
}

func TestDeterministicBackendPersistenceFailureFailsClosed(t *testing.T) {
	root := fixtureRoot(t)
	const command = "exit 0"
	s, contract := startedRun(t, root, backendReq(command))
	s.backendRunner = &recordingBackendRunner{
		process: deterministicBackendProcess{started: true, terminated: true, exitStatus: 0},
	}
	beforeState, _ := readBackendDurableFiles(t, s, contract.WorkflowID)
	s.saveStateOverride = func(*RunState) error { return errors.New("injected state persistence failure") }

	result, err := s.executeDeterministicBackend(
		context.Background(), contract.WorkflowID, command, io.Discard, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "no trustworthy ledger credit") ||
		!strings.Contains(err.Error(), "injected state persistence failure") {
		t.Fatalf("persistence failure error = %v", err)
	}
	if result.ExitStatus != 0 {
		t.Fatalf("observed process status = %d, want 0", result.ExitStatus)
	}
	afterState, ledger := readBackendDurableFiles(t, s, contract.WorkflowID)
	if !bytes.Equal(afterState, beforeState) {
		t.Error("failed state persistence unexpectedly advanced durable state")
	}
	if len(ledger) == 0 {
		t.Fatal("test did not reach the synced-ledger persistence boundary")
	}
	verification := s.VerifyLedger(contract.WorkflowID)
	if verification.Valid {
		t.Fatal("unanchored ledger event was treated as valid credit")
	}

	s.saveStateOverride = nil
	reconciliation, err := s.Close(contract.WorkflowID, futureBackendTime())
	if err != nil {
		t.Fatal(err)
	}
	if reconciliation.Verdict != Unverifiable {
		t.Fatalf("persistence-damaged run verdict = %q, want %q", reconciliation.Verdict, Unverifiable)
	}
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Evidence.Observations) != 0 || len(st.Evidence.IntegrityErrors) == 0 {
		t.Fatalf("persistence-damaged run received credit: %+v", st.Evidence)
	}
}

func TestDeterministicBackendSecondExecutionLatestSequenceGoverns(t *testing.T) {
	requireRealBackendShell(t)
	root := fixtureRoot(t)
	const command = "if [ -f .backend-second-execution ]; then exit 7; fi; " +
		": > .backend-second-execution; exit 0"
	s, contract := startedRun(t, root, backendReq(command))

	first, err := s.executeDeterministicBackend(
		context.Background(), contract.WorkflowID, command, io.Discard, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.executeDeterministicBackend(
		context.Background(), contract.WorkflowID, command, io.Discard, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if first.ExitStatus != 0 || second.ExitStatus != 7 {
		t.Fatalf("real execution statuses = [%d %d], want [0 7]", first.ExitStatus, second.ExitStatus)
	}
	events := verifiedBackendEvents(t, s, contract.WorkflowID)
	if len(events) != 2 {
		t.Fatalf("ledger event count = %d, want 2", len(events))
	}
	assertBackendFinishedEvent(t, events[0], 1, command, 0)
	assertBackendFinishedEvent(t, events[1], 2, command, 7)

	reconciliation, err := s.Close(contract.WorkflowID, futureBackendTime())
	if err != nil {
		t.Fatal(err)
	}
	if reconciliation.Verdict != Failed || len(reconciliation.Results) != 1 ||
		reconciliation.Results[0].Outcome != OutcomeUnsatisfied {
		t.Fatalf("latest nonzero execution did not govern: %+v", reconciliation)
	}
	st, err := s.LoadState(contract.WorkflowID)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Evidence.Observations) != 2 || st.Evidence.Observations[0].Sequence != 1 ||
		st.Evidence.Observations[1].Sequence != 2 || len(reconciliation.Results[0].Evidence) != 1 ||
		reconciliation.Results[0].Evidence[0] != st.Evidence.Observations[1].EvidenceRef {
		t.Fatalf("latest-sequence evidence binding is wrong: state=%+v reconciliation=%+v", st, reconciliation)
	}
}
