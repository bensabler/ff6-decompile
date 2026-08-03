package workflow

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	deterministicBackendShell       = "/bin/sh"
	deterministicBackendProvider    = "local-process"
	deterministicBackendCollectorID = "internal/workflow.deterministic-backend"
)

type deterministicBackendProcess struct {
	started    bool
	terminated bool
	exitStatus int
}

type deterministicBackendRunner interface {
	run(context.Context, string, string, io.Writer, io.Writer) (deterministicBackendProcess, error)
}

type shellDeterministicBackendRunner struct{}

// run passes the exact frozen command as one /bin/sh -c argument. Start errors
// are launch failures. Once ProcessState exists, its ExitCode is the observed
// outcome; shell output and the error text returned by Wait are not evidence.
func (shellDeterministicBackendRunner) run(ctx context.Context, root, command string,
	stdout, stderr io.Writer) (deterministicBackendProcess, error) {
	cmd := exec.CommandContext(ctx, deterministicBackendShell, "-c", command)
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return deterministicBackendProcess{}, err
	}
	process := deterministicBackendProcess{started: true}
	waitErr := cmd.Wait()
	if cmd.ProcessState == nil {
		if waitErr == nil {
			waitErr = fmt.Errorf("process wait returned without process state")
		}
		return process, waitErr
	}
	process.terminated = true
	process.exitStatus = cmd.ProcessState.ExitCode()
	return process, nil
}

type deterministicBackendPlan struct {
	workflowID     string
	runID          string
	contractHash   string
	command        string
	repositoryRoot string
}

type deterministicBackendResult struct {
	ExitStatus int
}

// executeDeterministicBackend executes only an exact deterministic-backend
// requirement selected from the active frozen contract. It is deliberately
// package-private: callers can select a requirement and route process output,
// but cannot assert a command, exit status, provenance, event, or verdict.
func (s *Store) executeDeterministicBackend(ctx context.Context, workflowID, requirementID string,
	stdout, stderr io.Writer) (deterministicBackendResult, error) {
	if ctx == nil {
		return deterministicBackendResult{}, fmt.Errorf("deterministic backend context is nil")
	}
	if !workflowIDPattern.MatchString(workflowID) {
		return deterministicBackendResult{}, fmt.Errorf("workflow_id %q must match WF-NNNN", workflowID)
	}
	plan, err := s.prepareDeterministicBackend(workflowID, requirementID)
	if err != nil {
		return deterministicBackendResult{}, err
	}
	if s.backendRunner == nil {
		return deterministicBackendResult{}, fmt.Errorf("deterministic backend runner is unavailable")
	}
	process, err := s.backendRunner.run(ctx, plan.repositoryRoot, plan.command, stdout, stderr)
	if err != nil {
		if process.started {
			return deterministicBackendResult{}, fmt.Errorf(
				"deterministic backend command %q started but no terminal process status was observed; no ledger event recorded: %w",
				plan.command, err)
		}
		return deterministicBackendResult{}, fmt.Errorf(
			"launch deterministic backend command %q: %w", plan.command, err)
	}
	if !process.started || !process.terminated {
		return deterministicBackendResult{}, fmt.Errorf(
			"deterministic backend command %q returned without an observed terminal process status; no ledger event recorded",
			plan.command)
	}

	result := deterministicBackendResult{ExitStatus: process.exitStatus}
	observedAt := time.Now().UTC().Format(time.RFC3339Nano)
	if err := s.recordDeterministicBackend(plan, observedAt, process); err != nil {
		return result, fmt.Errorf(
			"deterministic backend command %q may have run with exit status %d, but no trustworthy ledger credit was recorded: %w",
			plan.command, process.exitStatus, err)
	}
	return result, nil
}

func (s *Store) prepareDeterministicBackend(workflowID, requirementID string) (deterministicBackendPlan, error) {
	st, err := s.LoadState(workflowID)
	if err != nil {
		return deterministicBackendPlan{}, fmt.Errorf("load workflow %s run state: %w", workflowID, err)
	}
	expectedRunID := st.RunID
	if !runIDPattern.MatchString(expectedRunID) {
		return deterministicBackendPlan{}, fmt.Errorf("workflow %s state has invalid run_id %q",
			workflowID, expectedRunID)
	}
	release, err := s.acquireLedgerLock(workflowID, expectedRunID)
	if err != nil {
		return deterministicBackendPlan{}, err
	}
	defer release()

	st, err = s.LoadState(workflowID)
	if err != nil {
		return deterministicBackendPlan{}, fmt.Errorf("reload workflow %s run state: %w", workflowID, err)
	}
	if st.RunID != expectedRunID {
		return deterministicBackendPlan{}, fmt.Errorf("workflow %s run changed from %s to %s before execution",
			workflowID, expectedRunID, st.RunID)
	}
	if st.Phase != PhaseExecuting {
		return deterministicBackendPlan{}, fmt.Errorf("workflow %s run %s is not executable in phase %q",
			workflowID, st.RunID, st.Phase)
	}
	identity, err := s.loadVerifiedRunIdentityForState(workflowID, st)
	if err != nil {
		return deterministicBackendPlan{}, err
	}
	verification := verifyLedger(*identity, s.EventsPath(workflowID, st.RunID), "",
		st.LedgerSequence, st.LedgerHash)
	if !verification.Valid {
		return deterministicBackendPlan{}, fmt.Errorf("refusing to execute against invalid ledger: %s",
			strings.Join(verification.Problems, "; "))
	}
	requirement, err := s.loadDeterministicBackendRequirement(workflowID, requirementID, st)
	if err != nil {
		return deterministicBackendPlan{}, err
	}
	repository, err := resolveRepository(identity.RepositoryRoot)
	if err != nil {
		return deterministicBackendPlan{}, err
	}
	if repository.Root != identity.RepositoryRoot || repository.Identity != identity.RepositoryIdentity {
		return deterministicBackendPlan{}, fmt.Errorf(
			"current repository root or identity does not match immutable run identity")
	}
	return deterministicBackendPlan{
		workflowID: workflowID, runID: st.RunID, contractHash: st.ContractHash,
		command: requirement.ResourceID, repositoryRoot: identity.RepositoryRoot,
	}, nil
}

func (s *Store) loadDeterministicBackendRequirement(workflowID, requirementID string,
	st *RunState) (Requirement, error) {
	c, err := s.LoadContract(workflowID)
	if err != nil {
		return Requirement{}, fmt.Errorf("load workflow %s contract: %w", workflowID, err)
	}
	if c.State != Frozen {
		return Requirement{}, fmt.Errorf("workflow %s contract is not frozen: state is %q", workflowID, c.State)
	}
	if err := c.VerifyFrozen(); err != nil {
		return Requirement{}, fmt.Errorf("workflow %s contract is not frozen and intact: %w", workflowID, err)
	}
	if c.WorkflowID != workflowID {
		return Requirement{}, fmt.Errorf("stored contract workflow_id %q does not match %q", c.WorkflowID, workflowID)
	}
	if c.FrozenHash != st.ContractHash {
		return Requirement{}, fmt.Errorf(
			"contract for %s changed since the run started: state records %s, contract carries %s",
			workflowID, st.ContractHash, c.FrozenHash)
	}
	var matches []Requirement
	for _, requirement := range c.Requirements {
		if requirement.ResourceID == requirementID {
			matches = append(matches, requirement)
		}
	}
	if len(matches) == 0 {
		return Requirement{}, fmt.Errorf("frozen contract has no requirement with exact resource_id %q", requirementID)
	}
	if len(matches) != 1 {
		return Requirement{}, fmt.Errorf("frozen contract requirement resource_id %q is ambiguous: found %d matches",
			requirementID, len(matches))
	}
	if errs := c.Validate(); len(errs) > 0 {
		return Requirement{}, fmt.Errorf("frozen contract is structurally invalid: %w", errs[0])
	}
	if matches[0].Mode != ModeDeterministicBackend {
		return Requirement{}, fmt.Errorf("frozen requirement %q uses execution_mode %q, not %q",
			requirementID, matches[0].Mode, ModeDeterministicBackend)
	}
	return matches[0], nil
}

func (s *Store) recordDeterministicBackend(plan deterministicBackendPlan, observedAt string,
	process deterministicBackendProcess) error {
	if !process.started || !process.terminated {
		return fmt.Errorf("deterministic backend process has no observed terminal status")
	}
	return s.appendEventAllocatingSequence(plan.workflowID, plan.runID,
		func(context lockedAppendContext) (RunEvent, error) {
			if context.state.Phase != PhaseExecuting {
				return RunEvent{}, fmt.Errorf("workflow %s run %s is no longer executable in phase %q",
					plan.workflowID, context.state.RunID, context.state.Phase)
			}
			if context.state.ContractHash != plan.contractHash {
				return RunEvent{}, fmt.Errorf("workflow %s contract binding changed during backend execution",
					plan.workflowID)
			}
			if context.identity.RepositoryRoot != plan.repositoryRoot {
				return RunEvent{}, fmt.Errorf("workflow %s repository root changed during backend execution",
					plan.workflowID)
			}
			requirement, err := s.loadDeterministicBackendRequirement(
				plan.workflowID, plan.command, context.state)
			if err != nil {
				return RunEvent{}, err
			}
			if requirement.ResourceID != plan.command {
				return RunEvent{}, fmt.Errorf("frozen deterministic backend changed during execution")
			}
			repository, err := resolveRepository(context.identity.RepositoryRoot)
			if err != nil {
				return RunEvent{}, err
			}
			if repository.Root != context.identity.RepositoryRoot ||
				repository.Identity != context.identity.RepositoryIdentity {
				return RunEvent{}, fmt.Errorf(
					"current repository root or identity does not match immutable run identity")
			}
			eventID, err := newEventID()
			if err != nil {
				return RunEvent{}, err
			}
			status := process.exitStatus
			return RunEvent{
				SchemaVersion: RunEventSchemaVersion, EventID: eventID,
				WorkflowID: context.identity.WorkflowID, RunID: context.identity.RunID,
				ContractHash: context.identity.ContractHash, ObservedAt: observedAt,
				Provider: deterministicBackendProvider, SourceKind: SourceDeterministicBackend,
				CollectorID: deterministicBackendCollectorID, TrustBasis: TrustBackendExitStatus,
				SessionID: "", TurnID: "", RepositoryIdentity: repository.Identity,
				CWD: context.identity.RepositoryRoot, Branch: repository.Branch, Head: repository.Head,
				EventKind: EventBackendFinished, Selector: requirement.ResourceID,
				ToolUseID: "", ExitStatus: &status, EvidenceRef: "",
			}, nil
		})
}
