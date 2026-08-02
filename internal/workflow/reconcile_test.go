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

func backendReq(id string) Requirement {
	return Requirement{ResourceID: id, ResourceType: "backend", Necessity: Required,
		Mode: ModeDeterministicBackend, EvidenceRule: "captured exit status",
		Policy: BlockCompletion}
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

func TestDeterministicBackendLatestSequenceGoverns(t *testing.T) {
	const selector = "go test ./..."
	observation := func(sequence uint64, observedSelector string, status *int, ref string) Observation {
		return Observation{Kind: ObsBackendRun, Sequence: sequence, Selector: observedSelector,
			ExitStatus: status, EvidenceRef: ref}
	}
	tests := []struct {
		name         string
		observations []Observation
		wantOutcome  Outcome
		wantVerdict  Verdict
		wantEvidence string
	}{
		{
			name: "pass then fail",
			observations: []Observation{
				observation(1, selector, ptr(0), "EV-1"),
				observation(2, selector, ptr(2), "EV-2"),
			},
			wantOutcome: OutcomeUnsatisfied, wantVerdict: Failed, wantEvidence: "EV-2",
		},
		{
			name: "fail then pass",
			observations: []Observation{
				observation(1, selector, ptr(2), "EV-1"),
				observation(2, selector, ptr(0), "EV-2"),
			},
			wantOutcome: OutcomeSatisfied, wantVerdict: Complete, wantEvidence: "EV-2",
		},
		{
			name: "pass then unverifiable",
			observations: []Observation{
				observation(1, selector, ptr(0), "EV-1"),
				observation(2, selector, nil, "EV-2"),
			},
			wantOutcome: OutcomeUnverifiable, wantVerdict: Unverifiable, wantEvidence: "EV-2",
		},
		{
			name: "unverifiable then pass",
			observations: []Observation{
				observation(1, selector, nil, "EV-1"),
				observation(2, selector, ptr(0), "EV-2"),
			},
			wantOutcome: OutcomeSatisfied, wantVerdict: Complete, wantEvidence: "EV-2",
		},
		{
			name: "no matching run",
			observations: []Observation{
				observation(1, "go vet ./...", ptr(0), "EV-other"),
			},
			wantOutcome: OutcomeUnsatisfied, wantVerdict: Partial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := Reconcile(frozen(t, backendReq(selector)), Evidence{Observations: tt.observations})
			if rec.Verdict != tt.wantVerdict {
				t.Fatalf("verdict = %q, want %q: %+v", rec.Verdict, tt.wantVerdict, rec)
			}
			if len(rec.Results) != 1 || rec.Results[0].Outcome != tt.wantOutcome {
				t.Fatalf("results = %+v, want outcome %q", rec.Results, tt.wantOutcome)
			}
			if tt.wantEvidence == "" {
				if len(rec.Results[0].Evidence) != 0 {
					t.Errorf("evidence = %v, want none", rec.Results[0].Evidence)
				}
				return
			}
			if len(rec.Results[0].Evidence) != 1 || rec.Results[0].Evidence[0] != tt.wantEvidence {
				t.Errorf("evidence = %v, want governing ref %q", rec.Results[0].Evidence, tt.wantEvidence)
			}
		})
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
