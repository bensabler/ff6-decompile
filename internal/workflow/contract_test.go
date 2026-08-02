package workflow

import (
	"strings"
	"testing"
)

func validContract() *Contract {
	return &Contract{
		SchemaVersion: SchemaVersion, WorkflowID: "WF-0001", Workflow: "extract",
		Scope: "Narshe field tiles", State: Draft,
		StoppingConditions: []string{"scope boundary reached"},
		Requirements: []Requirement{{
			ResourceID: "graphics-researcher", ResourceType: "agent", Necessity: Required,
			Mode: ModeExplicitAgentCall, EvidenceRule: "exact subagent_type",
			Policy: BlockCompletion}},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Contract)
		wantErr string // substring; empty means the contract must validate
	}{
		{name: "valid", mutate: func(*Contract) {}},
		{name: "wrong schema version",
			mutate:  func(c *Contract) { c.SchemaVersion = "0.9" },
			wantErr: "schema_version"},
		{name: "empty scope cannot stop",
			mutate:  func(c *Contract) { c.Scope = "  " },
			wantErr: "scope is empty"},
		{name: "no requirements",
			mutate:  func(c *Contract) { c.Requirements = nil },
			wantErr: "no requirements"},
		{name: "no stopping conditions",
			mutate:  func(c *Contract) { c.StoppingConditions = nil },
			wantErr: "no stopping conditions"},
		{name: "duplicate resource id",
			mutate:  func(c *Contract) { c.Requirements = append(c.Requirements, c.Requirements[0]) },
			wantErr: "duplicate resource_id"},
		{name: "unknown execution mode",
			mutate:  func(c *Contract) { c.Requirements[0].Mode = "vibes" },
			wantErr: "not a known mode"},
		{name: "unknown failure policy",
			mutate:  func(c *Contract) { c.Requirements[0].Policy = "maybe" },
			wantErr: "not a known policy"},
		{name: "empty evidence rule",
			mutate:  func(c *Contract) { c.Requirements[0].EvidenceRule = "" },
			wantErr: "evidence_rule is empty"},
		{name: "required plus warn_only is not a requirement",
			mutate:  func(c *Contract) { c.Requirements[0].Policy = WarnOnly },
			wantErr: "may not use warn_only"},
		{name: "conditional without applicability rule",
			mutate:  func(c *Contract) { c.Requirements[0].Necessity = Conditional },
			wantErr: "need an applicability_rule"},
		{name: "not_applicable on a required requirement",
			mutate: func(c *Contract) {
				c.Requirements[0].Mode = ModeNotApplicable
				c.Requirements[0].NotApplicableReason = "because"
			},
			wantErr: "legal only on a conditional"},
		{name: "not_applicable without a preserved reason",
			mutate: func(c *Contract) {
				c.Requirements[0].Necessity = Conditional
				c.Requirements[0].Applicability = "when audio is in scope"
				c.Requirements[0].Mode = ModeNotApplicable
			},
			wantErr: "needs a preserved reason"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validContract()
			tt.mutate(c)
			errs := c.Validate()
			if tt.wantErr == "" {
				if len(errs) != 0 {
					t.Fatalf("want valid, got %v", errs)
				}
				return
			}
			for _, e := range errs {
				if strings.Contains(e.Error(), tt.wantErr) {
					return
				}
			}
			t.Errorf("want an error containing %q, got %v", tt.wantErr, errs)
		})
	}
}

func TestFreezeRequiresApproval(t *testing.T) {
	c := validContract() // Draft
	if err := c.Freeze("t0"); err == nil {
		t.Fatal("freezing a draft must fail: the operator has not approved it")
	}
	for c.State != Approved {
		if err := c.Advance(); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.Freeze("t0"); err != nil {
		t.Fatalf("freeze approved contract: %v", err)
	}
	if c.State != Frozen || c.FrozenHash == "" {
		t.Fatalf("state=%q hash=%q", c.State, c.FrozenHash)
	}
	if err := c.VerifyFrozen(); err != nil {
		t.Errorf("freshly frozen contract must verify: %v", err)
	}
}

func TestFreezeRefusesInvalidContract(t *testing.T) {
	c := validContract()
	c.StoppingConditions = nil
	c.State = Approved
	if err := c.Freeze("t0"); err == nil {
		t.Error("an invalid contract must not be freezable")
	}
}

// The hash covers the plan, not the act of recording it: requirement order and
// freeze timestamp must not change it, or a reordered contract would look edited.
func TestHashIsOrderAndTimestampStable(t *testing.T) {
	a := validContract()
	a.Requirements = append(a.Requirements, Requirement{
		ResourceID: "asset-librarian", ResourceType: "agent", Necessity: Required,
		Mode: ModeExplicitAgentCall, EvidenceRule: "exact subagent_type", Policy: BlockCompletion})
	b := validContract()
	b.Requirements = []Requirement{a.Requirements[1], a.Requirements[0]}

	ha, err := a.Hash()
	if err != nil {
		t.Fatal(err)
	}
	hb, err := b.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if ha != hb {
		t.Errorf("hash must not depend on requirement order:\n%s\n%s", ha, hb)
	}

	a.FrozenAt = "later"
	h2, err := a.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if h2 != ha {
		t.Error("hash must not depend on the freeze timestamp")
	}
}

func TestVerifyFrozenDetectsEdits(t *testing.T) {
	c := frozen(t, agentReq("graphics-researcher", Required, BlockCompletion))
	if err := c.VerifyFrozen(); err != nil {
		t.Fatalf("baseline: %v", err)
	}
	c.Scope = "quietly widened"
	if err := c.VerifyFrozen(); err == nil {
		t.Error("an edit after freezing must be detected")
	}
}

func TestAmend(t *testing.T) {
	base := frozen(t,
		agentReq("graphics-researcher", Required, BlockCompletion),
		agentReq("asset-librarian", Required, BlockCompletion))

	t.Run("preserves the original and records the amendment", func(t *testing.T) {
		next := *base
		next.Requirements = append([]Requirement(nil), base.Requirements...)
		next.Requirements = append(next.Requirements,
			agentReq("dma-researcher", Required, BlockCompletion))

		out, err := base.Amend(&next, "scope widened to the mines", "operator", "t1")
		if err != nil {
			t.Fatalf("amend: %v", err)
		}
		if out.State != Draft {
			t.Errorf("an amended contract must be re-approved, state=%q", out.State)
		}
		if out.FrozenHash != "" {
			t.Error("an amended contract must be re-frozen")
		}
		if len(out.Amendments) != 1 || out.Amendments[0].SupersedesHash != base.FrozenHash {
			t.Errorf("amendment must cite the contract it supersedes: %+v", out.Amendments)
		}
		if len(base.Requirements) != 2 {
			t.Error("the original contract must not be modified")
		}
	})

	weakenings := []struct {
		name    string
		mutate  func(*Contract)
		wantErr string
	}{
		{name: "removing a required resource",
			mutate:  func(c *Contract) { c.Requirements = c.Requirements[:1] },
			wantErr: "removes required resource"},
		{name: "downgrading required to conditional",
			mutate: func(c *Contract) {
				c.Requirements[1].Necessity = Conditional
				c.Requirements[1].Applicability = "whenever convenient"
			},
			wantErr: "downgrades required resource"},
		{name: "downgrading an invocation to context_only",
			mutate:  func(c *Contract) { c.Requirements[1].Mode = ModeContextOnly },
			wantErr: "weakens required resource"},
	}
	for _, tt := range weakenings {
		t.Run("refuses "+tt.name, func(t *testing.T) {
			next := *base
			next.Requirements = append([]Requirement(nil), base.Requirements...)
			tt.mutate(&next)
			_, err := base.Amend(&next, "reason", "operator", "t1")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("want an error containing %q, got %v", tt.wantErr, err)
			}
		})
	}

	t.Run("refuses an amendment without a reason or approval", func(t *testing.T) {
		next := *base
		if _, err := base.Amend(&next, "", "operator", "t1"); err == nil {
			t.Error("an amendment needs a reason")
		}
		if _, err := base.Amend(&next, "reason", "", "t1"); err == nil {
			t.Error("an amendment needs renewed approval")
		}
	})
}

func TestAdvanceRefusesToSkipStates(t *testing.T) {
	c := validContract()
	want := []State{Displayed, Approved, Frozen, Executing, Reconciled}
	for _, w := range want {
		if err := c.Advance(); err != nil {
			t.Fatalf("advance to %q: %v", w, err)
		}
		if c.State != w {
			t.Fatalf("state = %q, want %q", c.State, w)
		}
	}
	if err := c.Advance(); err == nil {
		t.Error("advancing past the final state must fail")
	}
}
