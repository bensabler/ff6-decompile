package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bensabler/ff6-decompile/internal/workflow"
)

// `ff6lab workflow` is the deterministic authority for workflow closure
// (AUDIT-0002 remediation R14).
//
// Claude may propose and execute work. Claude may not decide in prose that a
// workflow is complete. Every verdict printed here is computed from a frozen
// contract and observed evidence; no subcommand accepts a verdict as input, so
// there is no path by which a caller can assert one.
func runWorkflow(root string, args []string, out io.Writer) error {
	if len(args) < 2 {
		return workflowUsage()
	}
	s := workflow.NewStore(root)

	switch args[1] {
	case "list":
		ids, err := s.List()
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			fmt.Fprintln(out, "no workflow runs")
			return nil
		}
		for _, id := range ids {
			st, err := s.LoadState(id)
			if err != nil {
				fmt.Fprintf(out, "%s  (state unreadable: %v)\n", id, err)
				continue
			}
			verdict := "not reconciled"
			if st.Reconciliation != nil {
				verdict = string(st.Reconciliation.Verdict)
			}
			fmt.Fprintf(out, "%s  phase=%-9s verdict=%s\n", id, st.Phase, verdict)
		}
		return nil

	case "validate":
		if len(args) < 3 {
			return fmt.Errorf("usage: ff6lab workflow validate <contract.json>")
		}
		var c workflow.Contract
		if err := workflow.ReadJSON(args[2], &c); err != nil {
			return err
		}
		errs := c.Validate()
		for _, e := range errs {
			fmt.Fprintln(out, "contract:", e)
		}
		if len(errs) > 0 {
			return fmt.Errorf("contract: %d problem(s)", len(errs))
		}
		fmt.Fprintln(out, "contract: valid")
		return nil

	case "start":
		if len(args) < 3 {
			return fmt.Errorf("usage: ff6lab workflow start <contract.json>")
		}
		return workflowStart(s, args[2], out)

	case "status":
		if len(args) < 3 {
			return fmt.Errorf("usage: ff6lab workflow status <WF-NNNN>")
		}
		return workflowStatus(s, args[2], out)

	case "close":
		if len(args) < 3 {
			return fmt.Errorf("usage: ff6lab workflow close <WF-NNNN>")
		}
		return workflowClose(s, root, args[2], out)
	}
	return workflowUsage()
}

func workflowUsage() error {
	return fmt.Errorf("usage:\n" +
		"  ff6lab workflow list                    open and closed runs with their verdicts\n" +
		"  ff6lab workflow validate <contract>     check a contract before approval\n" +
		"  ff6lab workflow start <contract>        freeze an approved contract and open the run\n" +
		"  ff6lab workflow status <WF-NNNN>        planned vs observed, and what remains\n" +
		"  ff6lab workflow close <WF-NNNN>         reconcile against evidence and compute the verdict")
}

// workflowStart freezes an approved contract and opens the run. It refuses a
// contract that has not reached the approved state, so no run can begin before
// the operator has seen and accepted the plan.
func workflowStart(s *workflow.Store, path string, out io.Writer) error {
	var c workflow.Contract
	if err := workflow.ReadJSON(path, &c); err != nil {
		return err
	}
	if errs := c.Validate(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(out, "contract:", e)
		}
		return fmt.Errorf("contract: %d problem(s); refusing to start", len(errs))
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if c.State == workflow.Approved {
		if err := c.Freeze(now); err != nil {
			return err
		}
	}
	st, err := s.Start(&c, now)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s started\n  scope: %s\n  frozen: %s\n  requirements: %d\n",
		st.WorkflowID, c.Scope, c.FrozenHash[:16], len(c.Requirements))
	return nil
}

func workflowStatus(s *workflow.Store, id string, out io.Writer) error {
	c, err := s.LoadContract(id)
	if err != nil {
		return err
	}
	st, err := s.LoadState(id)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s  %s\n  scope: %s\n  phase: %s\n", id, c.Workflow, c.Scope, st.Phase)
	if st.Reconciliation != nil {
		fmt.Fprintf(out, "  verdict: %s\n", st.Reconciliation.Verdict)
	}
	if rem := st.Remaining(c); len(rem) > 0 {
		fmt.Fprintln(out, "  remaining:")
		for _, r := range rem {
			fmt.Fprintf(out, "    - %s\n", r)
		}
	} else {
		fmt.Fprintln(out, "  remaining: none")
	}
	return nil
}

func workflowClose(s *workflow.Store, root, id string, out io.Writer) error {
	c, err := s.LoadContract(id)
	if err != nil {
		return err
	}

	// Every source reports its own blindness. A source that cannot see is not a
	// source that saw nothing, so a missing transcript makes agent and skill
	// requirements unverifiable rather than unsatisfied.
	tr, err := workflow.TranscriptObservations(transcriptDir())
	if err != nil {
		return err
	}
	gl, err := workflow.GateLogObservations(
		filepath.Join(root, "local_artifacts", "workflow-runs", id, "gate-status.tsv"))
	if err != nil {
		return err
	}
	outs, err := workflow.OutputObservations(root, c.Outputs)
	if err != nil {
		return err
	}
	ev := workflow.Merge(tr, gl, workflow.Evidence{Observations: outs})

	rec, err := s.Close(id, ev, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", id)
	for _, r := range rec.Results {
		fmt.Fprintf(out, "  %-10s %-34s %s\n", r.Outcome, r.ResourceID, r.Reason)
	}
	for _, n := range rec.Notes {
		fmt.Fprintln(out, "  note:", n)
	}
	fmt.Fprintf(out, "verdict: %s\n", rec.Verdict)

	// A receipt that claims a different verdict fails validation; the computed
	// verdict stands regardless.
	if err := s.ValidateReceiptFile(id, rec.Verdict); err != nil {
		return err
	}
	if rec.Verdict != workflow.Complete {
		return fmt.Errorf("workflow %s is %s, not complete", id, rec.Verdict)
	}
	return nil
}

// transcriptDir is where the harness writes per-session transcripts. It is
// outside the repository, so it is resolved separately and its absence is a
// normal, reportable condition rather than an error: a missing transcript makes
// invocation requirements unverifiable, not unsatisfied.
func transcriptDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "projects")
}
