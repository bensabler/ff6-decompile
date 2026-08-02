package workflow

import (
	"fmt"
	"sort"
)

// ObservationKind is what was actually seen. The kinds are deliberately narrow:
// an artifact appearing is its own kind and can never stand in for an
// invocation. AUDIT-0001 credited two commands from artifact existence alone
// and was right about one of them by luck, which is exactly why the two are
// separated here.
type ObservationKind string

const (
	ObsAgentCall       ObservationKind = "agent_call"
	ObsSkillCall       ObservationKind = "skill_call"
	ObsBackendRun      ObservationKind = "backend_run"
	ObsOperatorAction  ObservationKind = "operator_action"
	ObsArtifactPresent ObservationKind = "artifact_present"
	// ObsMention records that a resource was named in prose or loaded as
	// context. It satisfies nothing except ModeContextOnly.
	ObsMention ObservationKind = "mention"
)

// GeneralPurposeAgent is the agent the harness uses when no subagent_type is
// given. AUDIT-0002 probes confirmed that omitting the name never selects a
// project specialist, so this value can never satisfy a named requirement.
const GeneralPurposeAgent = "general-purpose"

// Observation is one thing that was seen during execution.
type Observation struct {
	Kind ObservationKind `json:"kind"`
	// Selector is the exact identifier observed: a subagent_type, a skill name,
	// a command, or an output path.
	Selector string `json:"selector"`
	// ExitStatus is set for ObsBackendRun and nil when no status was captured.
	// A nil status is not a pass; it makes the requirement unverifiable.
	ExitStatus  *int   `json:"exit_status,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	EvidenceRef string `json:"evidence_ref,omitempty"`
}

// Evidence is the observation set a reconciliation runs against.
type Evidence struct {
	Observations []Observation `json:"observations"`
	// Incomplete names observation kinds whose record is known to be missing or
	// partial. Requirements discharged through those kinds become Unverifiable
	// rather than Unsatisfied: absence of evidence is not evidence of absence.
	Incomplete map[ObservationKind]string `json:"incomplete,omitempty"`
	// IntegrityErrors records whole-ledger failures. No observation from an
	// invalid ledger is admitted, and the workflow can never complete.
	IntegrityErrors []string `json:"integrity_errors,omitempty"`
	// Limitations explains structurally valid events that were not eligible to
	// satisfy an invocation because of provenance, trust, or event kind.
	Limitations []string `json:"limitations,omitempty"`
}

// Outcome is the per-requirement result.
type Outcome string

const (
	OutcomeSatisfied     Outcome = "satisfied"
	OutcomeUnsatisfied   Outcome = "unsatisfied"
	OutcomeUnverifiable  Outcome = "unverifiable"
	OutcomeNotApplicable Outcome = "not_applicable"
)

// Result explains one requirement's outcome in terms the receipt can print.
type Result struct {
	ResourceID string        `json:"resource_id"`
	Necessity  Necessity     `json:"requirement"`
	Mode       ExecutionMode `json:"execution_mode"`
	Outcome    Outcome       `json:"outcome"`
	Reason     string        `json:"reason"`
	Evidence   []string      `json:"evidence_refs,omitempty"`
}

// Verdict is the workflow status. It is computed, never asserted.
type Verdict string

const (
	// Complete: every required requirement satisfied, nothing unverifiable.
	Complete Verdict = "complete"
	// Partial: nothing failed outright, but something is missing or unproven.
	Partial Verdict = "partial"
	// Failed: a required backend exited non-zero, or a blocking requirement was
	// observed to be unsatisfied.
	Failed Verdict = "failed"
	// Unverifiable: the evidence record is too incomplete to judge.
	Unverifiable Verdict = "unverifiable"
	// Rejected: the contract was not frozen, or changed after freezing.
	Rejected Verdict = "rejected"
)

// Reconciliation is the full result of checking a contract against evidence.
type Reconciliation struct {
	WorkflowID string   `json:"workflow_id"`
	Verdict    Verdict  `json:"verdict"`
	Results    []Result `json:"results"`
	Notes      []string `json:"notes,omitempty"`
}

// Reconcile checks a frozen contract against observed evidence and computes the
// verdict. It never returns Complete on incomplete evidence, and it rejects a
// contract that was edited after freezing.
func Reconcile(c *Contract, ev Evidence) Reconciliation {
	rec := Reconciliation{WorkflowID: c.WorkflowID}

	if err := c.VerifyFrozen(); err != nil {
		rec.Verdict = Rejected
		rec.Notes = append(rec.Notes, err.Error())
		return rec
	}

	for _, r := range c.Requirements {
		rec.Results = append(rec.Results, check(r, ev))
	}
	for _, limitation := range ev.Limitations {
		rec.Notes = append(rec.Notes, limitation)
	}
	for _, integrityError := range ev.IntegrityErrors {
		rec.Notes = append(rec.Notes, "invalid evidence ledger: "+integrityError)
	}
	sort.Slice(rec.Results, func(i, j int) bool {
		return rec.Results[i].ResourceID < rec.Results[j].ResourceID
	})

	var blocking, degrading, unverifiable, failed int
	for _, res := range rec.Results {
		req := requirementFor(c, res.ResourceID)
		switch res.Outcome {
		case OutcomeUnverifiable:
			if req.Necessity == Required {
				unverifiable++
			}
		case OutcomeUnsatisfied:
			switch req.Policy {
			case BlockCompletion:
				if res.Reason == reasonBackendFailed {
					failed++
				} else {
					blocking++
				}
			case DegradeToPartial:
				degrading++
			case WarnOnly:
				rec.Notes = append(rec.Notes,
					fmt.Sprintf("%s unsatisfied (warn_only): %s", res.ResourceID, res.Reason))
			}
		}
	}

	switch {
	case len(ev.IntegrityErrors) > 0:
		rec.Verdict = Unverifiable
	case failed > 0:
		rec.Verdict = Failed
	case unverifiable > 0 && blocking == 0 && degrading == 0:
		rec.Verdict = Unverifiable
	case blocking > 0 || degrading > 0 || unverifiable > 0:
		rec.Verdict = Partial
	default:
		rec.Verdict = Complete
	}
	return rec
}

const reasonBackendFailed = "required backend exited non-zero"

func requirementFor(c *Contract, id string) Requirement {
	for _, r := range c.Requirements {
		if r.ResourceID == id {
			return r
		}
	}
	return Requirement{}
}

// check applies the evidence rule for one requirement's execution mode.
func check(r Requirement, ev Evidence) Result {
	res := Result{ResourceID: r.ResourceID, Necessity: r.Necessity, Mode: r.Mode}

	switch r.Mode {
	case ModeNotApplicable:
		// Legal only where the contract carries an applicability rule and a
		// preserved reason; Validate enforces both. The reason travels into the
		// receipt so a not-applicable decision stays auditable.
		res.Outcome, res.Reason = OutcomeNotApplicable,
			"conditional requirement not applicable: "+r.NotApplicableReason
		return res

	case ModeContextOnly:
		res.Outcome, res.Reason = OutcomeSatisfied, "context_only requires no invocation"
		return res

	case ModeExplicitAgentCall:
		if why, ok := ev.Incomplete[ObsAgentCall]; ok {
			res.Outcome, res.Reason = OutcomeUnverifiable, "agent-call record incomplete: "+why
			return res
		}
		for _, o := range ev.Observations {
			if o.Kind != ObsAgentCall {
				continue
			}
			if o.Selector == r.ResourceID {
				res.Outcome, res.Reason = OutcomeSatisfied, "invoked by exact subagent_type"
				res.Evidence = append(res.Evidence, o.EvidenceRef)
				return res
			}
		}
		// A general-purpose call is reported specifically, because substituting
		// it for a named specialist is a distinct and likely mistake.
		for _, o := range ev.Observations {
			if o.Kind == ObsAgentCall && o.Selector == GeneralPurposeAgent {
				res.Outcome, res.Reason = OutcomeUnsatisfied,
					"general-purpose was invoked; it does not satisfy a requirement for "+r.ResourceID
				return res
			}
		}
		res.Outcome, res.Reason = OutcomeUnsatisfied, "no agent call with subagent_type "+r.ResourceID
		return res

	case ModeExplicitSkillCall:
		if why, ok := ev.Incomplete[ObsSkillCall]; ok {
			res.Outcome, res.Reason = OutcomeUnverifiable, "skill-call record incomplete: "+why
			return res
		}
		for _, o := range ev.Observations {
			if o.Kind == ObsSkillCall && o.Selector == r.ResourceID {
				res.Outcome, res.Reason = OutcomeSatisfied, "skill invoked"
				res.Evidence = append(res.Evidence, o.EvidenceRef)
				return res
			}
		}
		res.Outcome, res.Reason = OutcomeUnsatisfied,
			"no skill invocation for "+r.ResourceID+"; being mentioned or loaded as context earns no credit"
		return res

	case ModeDeterministicBackend:
		if why, ok := ev.Incomplete[ObsBackendRun]; ok {
			res.Outcome, res.Reason = OutcomeUnverifiable, "backend record incomplete: "+why
			return res
		}
		for _, o := range ev.Observations {
			if o.Kind != ObsBackendRun || o.Selector != r.ResourceID {
				continue
			}
			if o.ExitStatus == nil {
				res.Outcome, res.Reason = OutcomeUnverifiable,
					"backend ran but no exit status was captured; a blank status is not a pass"
				return res
			}
			res.Evidence = append(res.Evidence, o.EvidenceRef)
			if *o.ExitStatus != 0 {
				res.Outcome, res.Reason = OutcomeUnsatisfied, reasonBackendFailed
				return res
			}
			res.Outcome, res.Reason = OutcomeSatisfied, "backend exited 0"
			return res
		}
		res.Outcome, res.Reason = OutcomeUnsatisfied, "backend "+r.ResourceID+" did not run"
		return res

	case ModeOperatorAction:
		if why, ok := ev.Incomplete[ObsOperatorAction]; ok {
			res.Outcome, res.Reason = OutcomeUnverifiable, "operator-action record incomplete: "+why
			return res
		}
		for _, o := range ev.Observations {
			if o.Kind == ObsOperatorAction && o.Selector == r.ResourceID {
				res.Outcome, res.Reason = OutcomeSatisfied, "operator action recorded"
				res.Evidence = append(res.Evidence, o.EvidenceRef)
				return res
			}
		}
		res.Outcome, res.Reason = OutcomeUnsatisfied, "no recorded operator action for "+r.ResourceID
		return res
	}

	res.Outcome, res.Reason = OutcomeUnverifiable, "unknown execution mode "+string(r.Mode)
	return res
}
