package workflow

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Durable, restart-safe workflow state for R14.
//
// Approved contracts remain tracked. Immutable run identity, mutable state,
// and raw event evidence remain under the repository's ignored run-state
// boundary. One human workflow ID is bound to one generated run ID.

const (
	runsDir        = "docs/workflows/runs"
	runStateDir    = "local_artifacts/workflows"
	contractFile   = "contract.json"
	receiptFile    = "receipt.md"
	stateFile      = "state.json"
	identityFile   = "identity.json"
	eventsFile     = "events.jsonl"
	ledgerLockFile = "events.lock"
)

var workflowIDPattern = regexp.MustCompile(`^WF-[0-9]{4}$`)

// Store reads and writes workflow runs beneath a repository root passed in
// explicitly, so tests run against fixture trees.
type Store struct {
	root          string
	generateRunID func() (string, error)
	createLedger  func(string) error
}

// NewStore returns a Store rooted at the given repository path.
func NewStore(root string) *Store {
	return &Store{root: root, generateRunID: newRunID, createLedger: createLedgerFile}
}

// ContractPath is the tracked location of a run's approved contract.
func (s *Store) ContractPath(id string) string {
	return filepath.Join(s.root, runsDir, id, contractFile)
}

// ReceiptPath is the tracked location of a run's receipt. R13 writes it; R14
// only validates what it claims.
func (s *Store) ReceiptPath(id string) string {
	return filepath.Join(s.root, runsDir, id, receiptFile)
}

// StatePath is the ignored location of a workflow's terminal-aware run state.
func (s *Store) StatePath(id string) string {
	return filepath.Join(s.root, runStateDir, id, stateFile)
}

// RunDir is the ignored directory holding one generated run's immutable
// identity and append-only event ledger.
func (s *Store) RunDir(id, runID string) string {
	return filepath.Join(s.root, runStateDir, id, runID)
}

// IdentityPath is the immutable identity record for one generated run.
func (s *Store) IdentityPath(id, runID string) string {
	return filepath.Join(s.RunDir(id, runID), identityFile)
}

// EventsPath is the append-only JSONL ledger for one generated run.
func (s *Store) EventsPath(id, runID string) string {
	return filepath.Join(s.RunDir(id, runID), eventsFile)
}

// Phase is where a run has reached. It is derived from recorded facts, never
// asserted: a run is Closed only once a reconciliation has been stored.
type Phase string

const (
	PhasePlanned   Phase = "planned"
	PhaseExecuting Phase = "executing"
	PhaseClosed    Phase = "closed"
)

// RunState is the restart-safe record of one workflow run. LedgerSequence and
// LedgerHash anchor the last synced append, including before closure, so
// deleting the final complete event is detectable.
type RunState struct {
	WorkflowID     string          `json:"workflow_id"`
	RunID          string          `json:"run_id"`
	Phase          Phase           `json:"phase"`
	StartedAt      string          `json:"started_at"`
	UpdatedAt      string          `json:"updated_at,omitempty"`
	ContractHash   string          `json:"contract_hash"`
	IdentityHash   string          `json:"identity_hash"`
	LedgerSequence uint64          `json:"ledger_sequence"`
	LedgerHash     string          `json:"ledger_hash"`
	Evidence       Evidence        `json:"evidence"`
	Reconciliation *Reconciliation `json:"reconciliation,omitempty"`
}

// Remaining lists the requirements not yet satisfied, so a resumed session is
// told what is left rather than having to work it out.
func (st *RunState) Remaining(c *Contract) []string {
	if st.Reconciliation == nil {
		ids := make([]string, 0, len(c.Requirements))
		for _, r := range c.Requirements {
			ids = append(ids, r.ResourceID)
		}
		return ids
	}
	var out []string
	for _, res := range st.Reconciliation.Results {
		if res.Outcome != OutcomeSatisfied && res.Outcome != OutcomeNotApplicable {
			out = append(out, res.ResourceID)
		}
	}
	return out
}

// Start records a frozen contract and creates an immutable run identity plus
// an empty restrictive ledger. The tracked contract is the final publication
// step. Any ordinary creation failure removes the new runtime bundle, so a
// failed start cannot be mistaken for an open run.
func (s *Store) Start(c *Contract, at string) (*RunState, error) {
	if !workflowIDPattern.MatchString(c.WorkflowID) {
		return nil, fmt.Errorf("workflow_id %q must match WF-NNNN", c.WorkflowID)
	}
	if err := c.VerifyFrozen(); err != nil {
		return nil, fmt.Errorf("cannot start a run from an unfrozen contract: %w", err)
	}
	if _, err := time.Parse(time.RFC3339Nano, at); err != nil {
		return nil, fmt.Errorf("invalid run creation time: %w", err)
	}
	if err := requireAbsent(s.ContractPath(c.WorkflowID), "tracked contract"); err != nil {
		return nil, err
	}
	if err := requireAbsent(s.StatePath(c.WorkflowID), "run state"); err != nil {
		return nil, err
	}
	runtimeRoot := filepath.Join(s.root, runStateDir, c.WorkflowID)
	if err := requireAbsent(runtimeRoot, "workflow runtime"); err != nil {
		return nil, err
	}

	repository, err := resolveRepository(s.root)
	if err != nil {
		return nil, err
	}
	runID, err := s.generateRunID()
	if err != nil {
		return nil, err
	}
	identity := RunIdentity{
		SchemaVersion: RunIdentitySchemaVersion, WorkflowID: c.WorkflowID,
		RunID: runID, ContractHash: c.FrozenHash, CreatedAt: at,
		RepositoryRoot: repository.Root, RepositoryIdentity: repository.Identity,
		StartingBranch: repository.Branch, StartingHead: repository.Head,
	}
	if problems := identity.validate(); len(problems) > 0 {
		return nil, fmt.Errorf("invalid generated run identity: %s", strings.Join(problems, "; "))
	}
	identityHash, err := runIdentityHash(identity)
	if err != nil {
		return nil, err
	}

	runDir := s.RunDir(c.WorkflowID, runID)
	if err := os.MkdirAll(runDir, 0o700); err != nil {
		return nil, fmt.Errorf("create run directory: %w", err)
	}
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.RemoveAll(runtimeRoot)
			_ = os.Remove(s.ContractPath(c.WorkflowID))
		}
	}()
	if err := writeJSONExclusive(s.IdentityPath(c.WorkflowID, runID), identity, 0o700, 0o600); err != nil {
		return nil, fmt.Errorf("create run identity: %w", err)
	}
	if err := s.createLedger(s.EventsPath(c.WorkflowID, runID)); err != nil {
		return nil, fmt.Errorf("create run ledger after identity: %w", err)
	}
	st := &RunState{
		WorkflowID: c.WorkflowID, RunID: runID, Phase: PhaseExecuting,
		StartedAt: at, ContractHash: c.FrozenHash, IdentityHash: identityHash,
		LedgerHash: genesisEventHash,
	}
	if err := writeJSONExclusive(s.StatePath(c.WorkflowID), st, 0o700, 0o600); err != nil {
		return nil, fmt.Errorf("create run state: %w", err)
	}
	if err := syncDirectory(runtimeRoot); err != nil {
		return nil, err
	}
	if err := writeJSONExclusive(s.ContractPath(c.WorkflowID), c, 0o755, 0o644); err != nil {
		return nil, fmt.Errorf("publish frozen contract: %w", err)
	}
	if err := syncDirectory(filepath.Dir(s.ContractPath(c.WorkflowID))); err != nil {
		return nil, err
	}
	cleanup = false
	return st, nil
}

func requireAbsent(path, label string) error {
	_, err := os.Lstat(path)
	if err == nil {
		return fmt.Errorf("%s already exists at %s; refusing to overwrite", label, path)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("inspect %s at %s: %w", label, path, err)
	}
	return nil
}

func createLedgerFile(path string) error {
	if err := writeFileExclusive(path, nil, 0o600); err != nil {
		return err
	}
	return syncDirectory(filepath.Dir(path))
}

// LoadContract reads a run's approved contract.
func (s *Store) LoadContract(id string) (*Contract, error) {
	var c Contract
	if err := ReadJSON(s.ContractPath(id), &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// LoadState reads a run's execution state.
func (s *Store) LoadState(id string) (*RunState, error) {
	var st RunState
	if err := ReadJSON(s.StatePath(id), &st); err != nil {
		return nil, err
	}
	return &st, nil
}

// LoadRunIdentity reads and validates the immutable identity for a state.
func (s *Store) LoadRunIdentity(id string) (*RunIdentity, error) {
	st, err := s.LoadState(id)
	if err != nil {
		return nil, err
	}
	var identity RunIdentity
	if err := ReadJSON(s.IdentityPath(id, st.RunID), &identity); err != nil {
		return nil, err
	}
	if problems := identity.validate(); len(problems) > 0 {
		return nil, fmt.Errorf("invalid run identity: %s", strings.Join(problems, "; "))
	}
	return &identity, nil
}

func (s *Store) saveState(st *RunState) error {
	return writeJSONAtomic(s.StatePath(st.WorkflowID), st, 0o600)
}

// List returns every workflow ID with a stored contract, sorted.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(s.root, runsDir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() && workflowIDPattern.MatchString(e.Name()) {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// appendEvent appends one internally collected event. It is package-private so
// this atomic unit cannot become a generic production assertion surface before
// a reviewed adapter defines a narrower capability. Identity, time, sequence,
// and duplicate checks happen under the same per-run lock as closure.
func (s *Store) appendEvent(id string, event RunEvent) error {
	st, err := s.LoadState(id)
	if err != nil {
		return err
	}
	release, err := s.acquireLedgerLock(id, st.RunID)
	if err != nil {
		return err
	}
	defer release()

	st, err = s.LoadState(id)
	if err != nil {
		return err
	}
	if st.Phase == PhaseClosed {
		return fmt.Errorf("workflow %s run %s is terminally closed; event append rejected", id, st.RunID)
	}
	identity, err := s.LoadRunIdentity(id)
	if err != nil {
		return err
	}
	verification := verifyLedger(*identity, s.EventsPath(id, st.RunID), "", st.LedgerSequence, st.LedgerHash)
	if !verification.Valid {
		return fmt.Errorf("refusing to append to invalid ledger: %s", strings.Join(verification.Problems, "; "))
	}
	expectedSequence := verification.TailSequence + 1
	if event.Sequence != expectedSequence {
		return fmt.Errorf("event sequence %d, want next sequence %d", event.Sequence, expectedSequence)
	}
	for _, existing := range verification.Events {
		if event.EventID == existing.EventID {
			return fmt.Errorf("duplicate event_id %q", event.EventID)
		}
	}
	event.PreviousHash = verification.TailHash
	event.EventHash = ""
	hash, err := eventHash(event)
	if err != nil {
		return err
	}
	event.EventHash = hash
	created, _ := time.Parse(time.RFC3339Nano, identity.CreatedAt)
	if problems := validateEvent(event, *identity, created, time.Time{}); len(problems) > 0 {
		return fmt.Errorf("reject event: %s", strings.Join(problems, "; "))
	}
	line, err := marshalJSONLine(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	if err := appendAndSync(s.EventsPath(id, st.RunID), line); err != nil {
		return err
	}
	st.LedgerSequence, st.LedgerHash = event.Sequence, event.EventHash
	if err := s.saveState(st); err != nil {
		return fmt.Errorf("event was synced but ledger anchor update failed; run is invalid until externally reviewed: %w", err)
	}
	return nil
}

func appendAndSync(path string, line []byte) (err error) {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect ledger: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("ledger %s is not a regular file", path)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return fmt.Errorf("open ledger for append: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close appended ledger: %w", closeErr)
		}
	}()
	n, err := f.Write(line)
	if err != nil {
		return fmt.Errorf("append ledger: %w", err)
	}
	if n != len(line) {
		return fmt.Errorf("append ledger: wrote %d of %d bytes", n, len(line))
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync appended ledger: %w", err)
	}
	return nil
}

func (s *Store) acquireLedgerLock(id, runID string) (func(), error) {
	path := filepath.Join(s.RunDir(id, runID), ledgerLockFile)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, fmt.Errorf("acquire run ledger lock: %w", err)
	}
	return func() {
		_ = f.Close()
		_ = os.Remove(path)
	}, nil
}

// VerifyLedger validates the current run's immutable identity, complete JSONL
// history, hash chain, time bounds, and durable tail anchor.
func (s *Store) VerifyLedger(id string) LedgerVerification {
	st, err := s.LoadState(id)
	if err != nil {
		return LedgerVerification{Problems: []string{err.Error()}}
	}
	release, err := s.acquireLedgerLock(id, st.RunID)
	if err != nil {
		return LedgerVerification{Problems: []string{err.Error()}}
	}
	defer release()
	return s.verifyLedgerForState(id, st)
}

func (s *Store) verifyLedgerForState(id string, st *RunState) LedgerVerification {
	identity, err := s.LoadRunIdentity(id)
	if err != nil {
		return LedgerVerification{Problems: []string{err.Error()}}
	}
	var closedAt string
	if st.Phase == PhaseClosed {
		closedAt = st.UpdatedAt
	}
	verification := verifyLedger(*identity, s.EventsPath(id, st.RunID), closedAt,
		st.LedgerSequence, st.LedgerHash)
	identityHash, hashErr := runIdentityHash(*identity)
	if hashErr != nil {
		verification.Valid = false
		verification.Events = nil
		verification.Problems = append(verification.Problems, hashErr.Error())
	} else if identityHash != st.IdentityHash {
		verification.Valid = false
		verification.Events = nil
		verification.Problems = append(verification.Problems,
			"immutable run identity hash does not match durable state")
	}
	if identity.WorkflowID != id || identity.RunID != st.RunID || identity.ContractHash != st.ContractHash {
		verification.Valid = false
		verification.Events = nil
		verification.Problems = append(verification.Problems, "run identity does not match durable state")
	}
	return verification
}

// Close verifies the run ledger, converts only eligible verified events into
// observations, computes the verdict, and records a terminal state. A second
// close is rejected before mutation.
func (s *Store) Close(id, at string) (*Reconciliation, error) {
	closedAt, err := time.Parse(time.RFC3339Nano, at)
	if err != nil {
		return nil, fmt.Errorf("invalid closure time: %w", err)
	}
	c, err := s.LoadContract(id)
	if err != nil {
		return nil, err
	}
	st, err := s.LoadState(id)
	if err != nil {
		return nil, err
	}
	release, err := s.acquireLedgerLock(id, st.RunID)
	if err != nil {
		return nil, err
	}
	defer release()

	st, err = s.LoadState(id)
	if err != nil {
		return nil, err
	}
	if st.Phase == PhaseClosed {
		return nil, fmt.Errorf("workflow %s run %s is already terminally closed", id, st.RunID)
	}
	startedAt, err := time.Parse(time.RFC3339Nano, st.StartedAt)
	if err != nil || closedAt.Before(startedAt) {
		return nil, fmt.Errorf("closure time precedes or cannot be compared with run creation")
	}
	if st.ContractHash != c.FrozenHash {
		return nil, fmt.Errorf("contract for %s changed since the run started: state records %s, contract carries %s",
			id, st.ContractHash, c.FrozenHash)
	}

	verification := s.verifyLedgerForState(id, st)
	if verification.Valid {
		// Recheck event upper bounds against the proposed terminal timestamp.
		identity, identityErr := s.LoadRunIdentity(id)
		if identityErr != nil {
			verification = LedgerVerification{Problems: []string{identityErr.Error()}}
		} else {
			verification = verifyLedger(*identity, s.EventsPath(id, st.RunID), at,
				st.LedgerSequence, st.LedgerHash)
		}
	}
	ev := evidenceFromLedger(verification)
	rec := Reconcile(c, ev)
	st.Evidence, st.Reconciliation, st.UpdatedAt = ev, &rec, at
	st.Phase = PhaseClosed
	if err := s.saveState(st); err != nil {
		return nil, err
	}
	return &rec, nil
}

// ErrReceiptDisagrees reports a receipt claiming a verdict the verifier did not
// compute.
var ErrReceiptDisagrees = errors.New("receipt disagrees with the reconciled verdict")

// verdictLine finds the verdict a receipt claims. Receipts are human-readable,
// so the claim is carried on an explicit machine-readable line.
var verdictLine = regexp.MustCompile(`(?m)^\s*[-*]?\s*(?:\*\*)?[Vv]erdict(?:\*\*)?\s*[:=]\s*` +
	"`?(complete|partial|failed|unverifiable|rejected)`?")

// ValidateReceipt checks that a receipt's claimed verdict matches the computed
// one. This is acceptance test N9: a receipt saying "complete" while
// reconciliation says partial or failed must fail validation, and the
// verifier's verdict stands.
//
// A receipt that claims no verdict at all also fails: silence is not agreement.
func ValidateReceipt(receipt string, computed Verdict) error {
	m := verdictLine.FindStringSubmatch(receipt)
	if m == nil {
		return fmt.Errorf("%w: receipt states no verdict, but the reconciled verdict is %q",
			ErrReceiptDisagrees, computed)
	}
	claimed := Verdict(strings.ToLower(m[1]))
	if claimed != computed {
		return fmt.Errorf("%w: receipt claims %q, reconciliation computed %q",
			ErrReceiptDisagrees, claimed, computed)
	}
	return nil
}

// ValidateReceiptFile validates the receipt stored for a run against a computed
// verdict. A missing receipt is not an error here: R13 owns receipt generation,
// and a run may legitimately be reconciled before one exists.
func (s *Store) ValidateReceiptFile(id string, computed Verdict) error {
	b, err := os.ReadFile(s.ReceiptPath(id))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read receipt for %s: %w", id, err)
	}
	return ValidateReceipt(string(b), computed)
}
