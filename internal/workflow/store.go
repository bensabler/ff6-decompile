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
)

// Durable, restart-safe workflow state for R14.
//
// The approved contract and the receipt are tracked, because they are the
// record of what was agreed and what happened. Observations and reconciliation
// detail are bulky and evidence-shaped, so they stay ignored under
// local_artifacts. The split is the same one the project already applies to
// experiment evidence.

const (
	runsDir      = "docs/workflows/runs"
	runStateDir  = "local_artifacts/workflow-runs"
	contractFile = "contract.json"
	receiptFile  = "receipt.md"
	stateFile    = "state.json"
)

var workflowIDPattern = regexp.MustCompile(`^WF-[0-9]{4}$`)

// Store reads and writes workflow runs beneath a repository root passed in
// explicitly, so tests run against fixture trees.
type Store struct{ root string }

// NewStore returns a Store rooted at the given repository path.
func NewStore(root string) *Store { return &Store{root: root} }

// ContractPath is the tracked location of a run's approved contract.
func (s *Store) ContractPath(id string) string {
	return filepath.Join(s.root, runsDir, id, contractFile)
}

// ReceiptPath is the tracked location of a run's receipt. R13 writes it; R14
// only validates what it claims.
func (s *Store) ReceiptPath(id string) string {
	return filepath.Join(s.root, runsDir, id, receiptFile)
}

// StatePath is the ignored location of a run's mutable execution state.
func (s *Store) StatePath(id string) string {
	return filepath.Join(s.root, runStateDir, id, stateFile)
}

// Phase is where a run has reached. It is derived from recorded facts, never
// asserted: a run is Closed only once a reconciliation has been stored.
type Phase string

const (
	PhasePlanned   Phase = "planned"
	PhaseExecuting Phase = "executing"
	PhaseClosed    Phase = "closed"
)

// RunState is the restart-safe record of one workflow run. Reloading it is
// enough to answer "what is open, what is done, what remains" without replaying
// anything.
type RunState struct {
	WorkflowID     string          `json:"workflow_id"`
	Phase          Phase           `json:"phase"`
	StartedAt      string          `json:"started_at"`
	UpdatedAt      string          `json:"updated_at,omitempty"`
	ContractHash   string          `json:"contract_hash"`
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

// Start records a frozen contract as an open run. It refuses an unfrozen
// contract, so no run can begin before the operator has approved and frozen it.
func (s *Store) Start(c *Contract, at string) (*RunState, error) {
	if !workflowIDPattern.MatchString(c.WorkflowID) {
		return nil, fmt.Errorf("workflow_id %q must match WF-NNNN", c.WorkflowID)
	}
	if err := c.VerifyFrozen(); err != nil {
		return nil, fmt.Errorf("cannot start a run from an unfrozen contract: %w", err)
	}
	if _, err := os.Stat(s.ContractPath(c.WorkflowID)); err == nil {
		return nil, fmt.Errorf("run %s already exists; amend it rather than overwriting", c.WorkflowID)
	}
	if err := WriteJSON(s.ContractPath(c.WorkflowID), c); err != nil {
		return nil, err
	}
	st := &RunState{WorkflowID: c.WorkflowID, Phase: PhaseExecuting, StartedAt: at,
		ContractHash: c.FrozenHash}
	if err := WriteJSON(s.StatePath(c.WorkflowID), st); err != nil {
		return nil, err
	}
	return st, nil
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

// SaveState writes a run's execution state.
func (s *Store) SaveState(st *RunState) error {
	return WriteJSON(s.StatePath(st.WorkflowID), st)
}

// List returns every run id with a stored contract, sorted.
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

// Close reconciles a run against evidence and records the computed verdict.
//
// This is the authoritative step. The verdict is computed from the frozen
// contract and the observations; nothing in this path accepts a verdict as
// input, so no caller can assert one.
func (s *Store) Close(id string, ev Evidence, at string) (*Reconciliation, error) {
	c, err := s.LoadContract(id)
	if err != nil {
		return nil, err
	}
	st, err := s.LoadState(id)
	if err != nil {
		return nil, err
	}
	if st.ContractHash != c.FrozenHash {
		return nil, fmt.Errorf("contract for %s changed since the run started: state records %s, contract carries %s",
			id, st.ContractHash, c.FrozenHash)
	}

	rec := Reconcile(c, ev)
	st.Evidence, st.Reconciliation, st.UpdatedAt = ev, &rec, at
	st.Phase = PhaseClosed
	if err := s.SaveState(st); err != nil {
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
