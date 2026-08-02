// Package workflow implements the frozen execution contract behind the
// operator-facing workflow surface (AUDIT-0002 remediation item R12).
//
// A contract declares, before any work begins, which agents, skills, backend
// operations, outputs and validations a workflow requires, and what evidence
// would prove each one ran. It is displayed to the operator, approved, then
// frozen and hashed. After freezing, a requirement may not be weakened,
// removed, or reclassified as optional or not-applicable — a skipped
// requirement stays skipped, and the workflow cannot report completion.
//
// The rules encoded here come from two audits that found the opposite
// happening: work performed by hand being credited to a named command, and a
// completion status written before the phase that would have contradicted it.
// See docs/workflows/AUDIT-0002-remediation-plan-v2.md.
//
// This package decides whether a contract is well formed and whether its
// requirements were satisfied by a given set of observations. It does not
// gather observations and does not own durable state; that is R14.
package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SchemaVersion is the contract format this package reads and writes.
const SchemaVersion = "1.0"

// ExecutionMode says how a requirement is discharged, and therefore what
// evidence could prove it. The mode determines which observations count:
// an artifact never proves an invocation, and a general-purpose agent never
// satisfies a requirement for a named specialist.
type ExecutionMode string

const (
	// ModeExplicitAgentCall requires a subagent invoked by its exact
	// subagent_type. Omitting the name yields the general-purpose agent, which
	// does not satisfy this mode.
	ModeExplicitAgentCall ExecutionMode = "explicit_agent_call"
	// ModeExplicitSkillCall requires a skill invoked as a skill. A skill that
	// was available, mentioned in prose, or loaded as context earns nothing.
	ModeExplicitSkillCall ExecutionMode = "explicit_skill_call"
	// ModeDeterministicBackend requires a command whose exit status was captured.
	ModeDeterministicBackend ExecutionMode = "deterministic_backend"
	// ModeOperatorAction requires a recorded operator decision, such as approval.
	ModeOperatorAction ExecutionMode = "operator_action"
	// ModeContextOnly requires only that the resource informed the work. It is
	// satisfied without an invocation and must never be used for a resource the
	// workflow genuinely needs to run.
	ModeContextOnly ExecutionMode = "context_only"
	// ModeNotApplicable marks a conditional requirement whose applicability rule
	// did not hold. It is legal only on a conditional requirement carrying both
	// an applicability rule and a recorded reason.
	ModeNotApplicable ExecutionMode = "not_applicable"
)

var executionModes = map[ExecutionMode]bool{
	ModeExplicitAgentCall: true, ModeExplicitSkillCall: true,
	ModeDeterministicBackend: true, ModeOperatorAction: true,
	ModeContextOnly: true, ModeNotApplicable: true,
}

// Necessity distinguishes a requirement that must hold from one whose
// applicability is decided at execution time.
type Necessity string

const (
	Required    Necessity = "required"
	Conditional Necessity = "conditional"
)

// FailurePolicy says what an unsatisfied requirement does to the verdict.
type FailurePolicy string

const (
	// BlockCompletion prevents Complete. The default for anything required.
	BlockCompletion FailurePolicy = "block_completion"
	// DegradeToPartial permits Partial but never Complete.
	DegradeToPartial FailurePolicy = "degrade_to_partial"
	// WarnOnly records the gap without affecting the verdict. Legal only on a
	// conditional requirement — a required step that can be skipped silently is
	// not a requirement.
	WarnOnly FailurePolicy = "warn_only"
)

var failurePolicies = map[FailurePolicy]bool{
	BlockCompletion: true, DegradeToPartial: true, WarnOnly: true,
}

// State is the contract lifecycle position. Work may only begin at Frozen.
type State string

const (
	Draft      State = "draft"
	Displayed  State = "displayed"
	Approved   State = "approved"
	Frozen     State = "frozen"
	Executing  State = "executing"
	Reconciled State = "reconciled"
)

var lifecycle = map[State]State{
	Draft: Displayed, Displayed: Approved, Approved: Frozen,
	Frozen: Executing, Executing: Reconciled,
}

// Requirement is one thing a workflow must do, and the evidence that proves it.
type Requirement struct {
	ResourceID   string        `json:"resource_id"`
	ResourceType string        `json:"resource_type"`
	Necessity    Necessity     `json:"requirement"`
	Mode         ExecutionMode `json:"execution_mode"`
	// EvidenceRule states in prose what would prove this requirement ran. It is
	// recorded for the operator and the receipt; the machine check is Mode.
	EvidenceRule string        `json:"evidence_rule"`
	Policy       FailurePolicy `json:"failure_policy"`
	// Applicability is required on a conditional requirement: the rule under
	// which it applies. Without it, not_applicable cannot be justified.
	Applicability string `json:"applicability_rule,omitempty"`
	// NotApplicableReason is recorded when a conditional requirement is
	// discharged as not applicable. It is preserved, never discarded.
	NotApplicableReason string `json:"not_applicable_reason,omitempty"`
}

// Amendment records a post-approval change. The original contract is preserved;
// an amendment never edits it in place.
type Amendment struct {
	Reason         string `json:"reason"`
	ApprovedBy     string `json:"approved_by"`
	SupersedesHash string `json:"supersedes_hash"`
	FrozenHash     string `json:"frozen_hash"`
	CreatedAt      string `json:"created_at"`
}

// Contract is the frozen plan for one workflow run.
type Contract struct {
	SchemaVersion      string        `json:"schema_version"`
	WorkflowID         string        `json:"workflow_id"`
	Workflow           string        `json:"workflow"`
	Scope              string        `json:"scope"`
	State              State         `json:"state"`
	Requirements       []Requirement `json:"requirements"`
	Outputs            []string      `json:"required_outputs"`
	StoppingConditions []string      `json:"stopping_conditions"`
	ApprovalBoundaries []string      `json:"approval_boundaries"`
	FrozenAt           string        `json:"frozen_at,omitempty"`
	FrozenHash         string        `json:"frozen_hash,omitempty"`
	Amendments         []Amendment   `json:"amendments,omitempty"`
}

// Validate reports every structural problem with a contract, so a malformed
// contract fails once with a full list rather than one error at a time.
func (c *Contract) Validate() []error {
	var errs []error
	add := func(format string, a ...any) { errs = append(errs, fmt.Errorf(format, a...)) }

	if c.SchemaVersion != SchemaVersion {
		add("schema_version %q, want %q", c.SchemaVersion, SchemaVersion)
	}
	if strings.TrimSpace(c.WorkflowID) == "" {
		add("workflow_id is empty")
	}
	if strings.TrimSpace(c.Scope) == "" {
		add("scope is empty: a workflow with no bounded scope cannot stop")
	}
	if _, ok := lifecycle[c.State]; !ok && c.State != Reconciled {
		add("state %q is not a lifecycle state", c.State)
	}
	if len(c.Requirements) == 0 {
		add("contract declares no requirements")
	}
	if len(c.StoppingConditions) == 0 {
		add("contract declares no stopping conditions")
	}

	seen := make(map[string]bool, len(c.Requirements))
	for i, r := range c.Requirements {
		where := fmt.Sprintf("requirement %d (%s)", i, r.ResourceID)
		if strings.TrimSpace(r.ResourceID) == "" {
			add("%s: resource_id is empty", where)
		}
		if seen[r.ResourceID] {
			add("%s: duplicate resource_id", where)
		}
		seen[r.ResourceID] = true

		if r.Necessity != Required && r.Necessity != Conditional {
			add("%s: requirement %q must be required or conditional", where, r.Necessity)
		}
		if !executionModes[r.Mode] {
			add("%s: execution_mode %q is not a known mode", where, r.Mode)
		}
		if !failurePolicies[r.Policy] {
			add("%s: failure_policy %q is not a known policy", where, r.Policy)
		}
		if strings.TrimSpace(r.EvidenceRule) == "" {
			add("%s: evidence_rule is empty: nothing states what would prove this ran", where)
		}
		// A required step that can be skipped without consequence is not a
		// requirement. This is the rule AUDIT-0001 broke by reclassifying
		// skipped work rather than reporting it.
		if r.Necessity == Required && r.Policy == WarnOnly {
			add("%s: required requirements may not use warn_only", where)
		}
		if r.Necessity == Conditional && strings.TrimSpace(r.Applicability) == "" {
			add("%s: conditional requirements need an applicability_rule", where)
		}
		if r.Mode == ModeNotApplicable {
			if r.Necessity != Conditional {
				add("%s: not_applicable is legal only on a conditional requirement", where)
			}
			if strings.TrimSpace(r.NotApplicableReason) == "" {
				add("%s: not_applicable needs a preserved reason", where)
			}
		}
	}
	return errs
}

// canonical serialises the contract for hashing with the volatile freeze fields
// cleared, so a hash covers the plan rather than the act of recording it.
func (c *Contract) canonical() ([]byte, error) {
	clone := *c
	clone.FrozenAt = ""
	clone.FrozenHash = ""
	clone.State = Approved
	clone.Requirements = append([]Requirement(nil), c.Requirements...)
	sort.Slice(clone.Requirements, func(i, j int) bool {
		return clone.Requirements[i].ResourceID < clone.Requirements[j].ResourceID
	})
	b, err := json.Marshal(clone)
	if err != nil {
		return nil, fmt.Errorf("canonicalise contract: %w", err)
	}
	return b, nil
}

// Hash returns the contract's content hash, independent of when it was frozen.
func (c *Contract) Hash() (string, error) {
	b, err := c.canonical()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Advance moves the contract to the next lifecycle state. It refuses to skip a
// state, so a contract cannot reach Executing without having been approved and
// frozen.
func (c *Contract) Advance() error {
	next, ok := lifecycle[c.State]
	if !ok {
		return fmt.Errorf("state %q cannot advance", c.State)
	}
	c.State = next
	return nil
}

// Freeze records the approved contract's hash and timestamp. It is the only
// point after which execution may begin.
func (c *Contract) Freeze(at string) error {
	if c.State != Approved {
		return fmt.Errorf("freeze requires state %q, have %q", Approved, c.State)
	}
	if errs := c.Validate(); len(errs) > 0 {
		return fmt.Errorf("refusing to freeze an invalid contract: %w", errs[0])
	}
	h, err := c.Hash()
	if err != nil {
		return err
	}
	c.FrozenHash, c.FrozenAt, c.State = h, at, Frozen
	return nil
}

// VerifyFrozen reports whether a contract's content still matches the hash
// recorded when it was frozen. A mismatch means the contract was edited after
// approval without an amendment, and reconciliation must be rejected.
func (c *Contract) VerifyFrozen() error {
	if c.FrozenHash == "" {
		return fmt.Errorf("contract is not frozen")
	}
	h, err := c.Hash()
	if err != nil {
		return err
	}
	if h != c.FrozenHash {
		return fmt.Errorf("contract changed after freezing: content hash %s, frozen hash %s", h, c.FrozenHash)
	}
	return nil
}

// Amend produces a successor contract carrying an explicit amendment record.
// The receiver is not modified: the original contract is preserved so that what
// was originally approved stays readable. The successor must be re-approved and
// re-frozen before execution continues.
func (c *Contract) Amend(next *Contract, reason, approvedBy, at string) (*Contract, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, fmt.Errorf("an amendment needs a reason")
	}
	if strings.TrimSpace(approvedBy) == "" {
		return nil, fmt.Errorf("an amendment needs renewed approval")
	}
	if err := c.VerifyFrozen(); err != nil {
		return nil, fmt.Errorf("cannot amend: %w", err)
	}
	if err := weakeningCheck(c, next); err != nil {
		return nil, err
	}
	out := *next
	out.State = Draft
	out.FrozenAt, out.FrozenHash = "", ""
	out.Amendments = append(append([]Amendment(nil), c.Amendments...), Amendment{
		Reason: reason, ApprovedBy: approvedBy,
		SupersedesHash: c.FrozenHash, CreatedAt: at,
	})
	return &out, nil
}

// weakeningCheck refuses an amendment that drops a required resource or
// downgrades it to conditional, context_only or not_applicable. A requirement
// that was skipped may not be retroactively defined away.
func weakeningCheck(from, to *Contract) error {
	after := make(map[string]Requirement, len(to.Requirements))
	for _, r := range to.Requirements {
		after[r.ResourceID] = r
	}
	for _, before := range from.Requirements {
		if before.Necessity != Required {
			continue
		}
		now, ok := after[before.ResourceID]
		if !ok {
			return fmt.Errorf("amendment removes required resource %q", before.ResourceID)
		}
		if now.Necessity != Required {
			return fmt.Errorf("amendment downgrades required resource %q to %q", before.ResourceID, now.Necessity)
		}
		if now.Mode == ModeNotApplicable || (before.Mode != ModeContextOnly && now.Mode == ModeContextOnly) {
			return fmt.Errorf("amendment weakens required resource %q from %q to %q",
				before.ResourceID, before.Mode, now.Mode)
		}
	}
	return nil
}
