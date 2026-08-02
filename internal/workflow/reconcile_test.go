package workflow

import "testing"

func ptr(i int) *int { return &i }

// frozen builds a valid, frozen contract with the given requirements.
func frozen(t *testing.T, reqs ...Requirement) *Contract {
	t.Helper()
	c := &Contract{
		SchemaVersion: SchemaVersion, WorkflowID: "WF-0001", Workflow: "extract",
		Scope: "Narshe field tiles, milestone 02 to 04", State: Approved,
		Requirements: reqs, StoppingConditions: []string{"scope boundary reached"},
	}
	if err := c.Freeze("2026-08-02T00:00:00Z"); err != nil {
		t.Fatalf("freeze: %v", err)
	}
	return c
}

func agentReq(id string, n Necessity, p FailurePolicy) Requirement {
	r := Requirement{ResourceID: id, ResourceType: "agent", Necessity: n,
		Mode: ModeExplicitAgentCall, EvidenceRule: "transcript Agent call with exact subagent_type",
		Policy: p}
	if n == Conditional {
		r.Applicability = "only when the scope includes graphics"
	}
	return r
}

// The nine negative acceptance tests from AUDIT-0002 remediation plan v2.
// Each reproduces a failure mode actually observed in AUDIT-0001 or AUDIT-0002.
// N6 and N9 belong to R14, which owns the transcript reader and the receipt.
func TestNegativeAcceptance(t *testing.T) {
	tests := []struct {
		name     string
		contract func(*testing.T) *Contract
		evidence Evidence
		want     Verdict
		wantNot  Verdict
	}{
		{
			name: "N1 required specialist omitted cannot complete",
			contract: func(t *testing.T) *Contract {
				return frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
			},
			evidence: Evidence{},
			want:     Partial, wantNot: Complete,
		},
		{
			name: "N2 general-purpose does not substitute for a named specialist",
			contract: func(t *testing.T) *Contract {
				return frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsAgentCall, Selector: GeneralPurposeAgent, EvidenceRef: "EV-1"},
			}},
			want: Partial, wantNot: Complete,
		},
		{
			name: "N3 skill mentioned but not invoked earns no credit",
			contract: func(t *testing.T) *Contract {
				return frozen(t, Requirement{
					ResourceID: "dma-tracer", ResourceType: "skill", Necessity: Required,
					Mode: ModeExplicitSkillCall, EvidenceRule: "skill invocation record",
					Policy: BlockCompletion})
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsMention, Selector: "dma-tracer", EvidenceRef: "EV-2"},
			}},
			want: Partial, wantNot: Complete,
		},
		{
			name: "N4 conditional not applicable is accepted with a preserved reason",
			contract: func(t *testing.T) *Contract {
				r := agentReq("audio-researcher", Conditional, DegradeToPartial)
				r.Mode = ModeNotApplicable
				r.NotApplicableReason = "scope contains no audio data"
				return frozen(t, r)
			},
			evidence: Evidence{},
			want:     Complete,
		},
		{
			name: "N5 required backend exiting non-zero fails the workflow",
			contract: func(t *testing.T) *Contract {
				return frozen(t, Requirement{
					ResourceID: "ff6lab extract graphics", ResourceType: "backend",
					Necessity: Required, Mode: ModeDeterministicBackend,
					EvidenceRule: "captured exit status", Policy: BlockCompletion})
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsBackendRun, Selector: "ff6lab extract graphics",
					ExitStatus: ptr(2), EvidenceRef: "EV-3"},
			}},
			want: Failed, wantNot: Complete,
		},
		{
			name: "N8 artifact without invocation evidence is not credited",
			contract: func(t *testing.T) *Contract {
				return frozen(t, agentReq("asset-librarian", Required, BlockCompletion))
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsArtifactPresent, Selector: "asset-librarian", EvidenceRef: "EV-4"},
				{Kind: ObsArtifactPresent, Selector: "local_artifacts/archive/tiles.png"},
			}},
			want: Partial, wantNot: Complete,
		},
		{
			name: "missing exit status is unverifiable, never a pass",
			contract: func(t *testing.T) *Contract {
				return frozen(t, Requirement{
					ResourceID: "go test ./...", ResourceType: "backend", Necessity: Required,
					Mode: ModeDeterministicBackend, EvidenceRule: "captured exit status",
					Policy: BlockCompletion})
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsBackendRun, Selector: "go test ./...", ExitStatus: nil},
			}},
			want: Unverifiable, wantNot: Complete,
		},
		{
			name: "incomplete evidence record yields unverifiable, not unsatisfied",
			contract: func(t *testing.T) *Contract {
				return frozen(t, agentReq("verification-engineer", Required, BlockCompletion))
			},
			evidence: Evidence{Incomplete: map[ObservationKind]string{
				ObsAgentCall: "no transcript preserved for this session"}},
			want: Unverifiable, wantNot: Complete,
		},
		{
			name: "fully satisfied contract completes",
			contract: func(t *testing.T) *Contract {
				return frozen(t,
					agentReq("graphics-researcher", Required, BlockCompletion),
					Requirement{ResourceID: "ff6lab extract graphics", ResourceType: "backend",
						Necessity: Required, Mode: ModeDeterministicBackend,
						EvidenceRule: "captured exit status", Policy: BlockCompletion})
			},
			evidence: Evidence{Observations: []Observation{
				{Kind: ObsAgentCall, Selector: "graphics-researcher", EvidenceRef: "EV-5"},
				{Kind: ObsBackendRun, Selector: "ff6lab extract graphics",
					ExitStatus: ptr(0), EvidenceRef: "EV-6"},
			}},
			want: Complete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reconcile(tt.contract(t), tt.evidence)
			if got.Verdict != tt.want {
				t.Errorf("verdict = %q, want %q\nresults: %+v", got.Verdict, tt.want, got.Results)
			}
			if tt.wantNot != "" && got.Verdict == tt.wantNot {
				t.Errorf("verdict must never be %q here", tt.wantNot)
			}
		})
	}
}

// N7: a contract edited after freezing is rejected outright, before any
// requirement is examined.
func TestN7ContractChangedAfterFreezing(t *testing.T) {
	c := frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
	ev := Evidence{Observations: []Observation{
		{Kind: ObsAgentCall, Selector: "graphics-researcher", EvidenceRef: "EV-1"}}}

	if got := Reconcile(c, ev); got.Verdict != Complete {
		t.Fatalf("baseline verdict = %q, want complete", got.Verdict)
	}

	// Weaken the frozen contract in place, the way a skipped requirement would
	// be defined away.
	c.Requirements[0].Necessity = Conditional

	got := Reconcile(c, ev)
	if got.Verdict != Rejected {
		t.Errorf("verdict = %q, want %q for a post-freeze edit", got.Verdict, Rejected)
	}
	if len(got.Notes) == 0 {
		t.Error("rejection must explain itself")
	}
}

func TestWarnOnlyDoesNotBlock(t *testing.T) {
	r := agentReq("documentation-reviewer", Conditional, WarnOnly)
	got := Reconcile(frozen(t, r), Evidence{})
	if got.Verdict != Complete {
		t.Errorf("verdict = %q, want complete: warn_only must not block", got.Verdict)
	}
	if len(got.Notes) == 0 {
		t.Error("an unsatisfied warn_only requirement must still be reported")
	}
}

func TestUnfrozenContractIsRejected(t *testing.T) {
	c := &Contract{SchemaVersion: SchemaVersion, WorkflowID: "WF-0002", Scope: "s",
		State: Draft, StoppingConditions: []string{"x"},
		Requirements: []Requirement{agentReq("quality-reviewer", Required, BlockCompletion)}}
	if got := Reconcile(c, Evidence{}); got.Verdict != Rejected {
		t.Errorf("verdict = %q, want %q: execution may not begin before freezing", got.Verdict, Rejected)
	}
}
