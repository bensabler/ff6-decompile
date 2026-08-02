package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// contractJSON is an approved, valid contract with one required agent.
const contractJSON = `{"schema_version":"1.0","workflow_id":"WF-0042","workflow":"extract",
 "scope":"Narshe field tiles","state":"approved",
 "stopping_conditions":["scope boundary reached"],
 "requirements":[{"resource_id":"graphics-researcher","resource_type":"agent",
   "requirement":"required","execution_mode":"explicit_agent_call",
   "evidence_rule":"exact subagent_type","failure_policy":"block_completion"}]}`

func writeContract(t *testing.T, dir, body string) string {
	t.Helper()
	p := filepath.Join(dir, "contract.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunWorkflowLifecycle(t *testing.T) {
	root := t.TempDir()
	path := writeContract(t, root, contractJSON)

	var out bytes.Buffer
	if err := runWorkflow(root, []string{"workflow", "validate", path}, &out); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if !strings.Contains(out.String(), "valid") {
		t.Errorf("validate output = %q", out.String())
	}

	out.Reset()
	if err := runWorkflow(root, []string{"workflow", "start", path}, &out); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !strings.Contains(out.String(), "WF-0042 started") {
		t.Errorf("start output = %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "workflows", "runs", "WF-0042", "contract.json")); err != nil {
		t.Errorf("the approved contract must be tracked: %v", err)
	}

	out.Reset()
	if err := runWorkflow(root, []string{"workflow", "status", "WF-0042"}, &out); err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(out.String(), "graphics-researcher") {
		t.Errorf("status must report what remains, got %q", out.String())
	}

	// Starting an existing run must be refused rather than silently overwriting.
	out.Reset()
	if err := runWorkflow(root, []string{"workflow", "start", path}, &out); err == nil {
		t.Error("starting an existing run must fail")
	}
}

// close must exit non-zero when the verdict is anything but complete, so a
// caller cannot ignore a partial or unverifiable run.
func TestRunWorkflowCloseFailsWhenNotComplete(t *testing.T) {
	root := t.TempDir()
	path := writeContract(t, root, contractJSON)
	var out bytes.Buffer
	if err := runWorkflow(root, []string{"workflow", "start", path}, &out); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", filepath.Join(root, "no-home")) // no transcripts reachable
	out.Reset()
	err := runWorkflow(root, []string{"workflow", "close", "WF-0042"}, &out)
	if err == nil {
		t.Fatal("close must fail when the run is not complete")
	}
	if !strings.Contains(out.String(), "verdict:") {
		t.Errorf("close must print the computed verdict, got %q", out.String())
	}
	if strings.Contains(out.String(), "verdict: complete") {
		t.Error("a run with unverifiable invocations must never be complete")
	}
}

func TestRunWorkflowRefusesInvalidContract(t *testing.T) {
	root := t.TempDir()
	path := writeContract(t, root, strings.Replace(contractJSON,
		`"stopping_conditions":["scope boundary reached"],`, `"stopping_conditions":[],`, 1))
	var out bytes.Buffer
	if err := runWorkflow(root, []string{"workflow", "start", path}, &out); err == nil {
		t.Error("an invalid contract must not start a run")
	}
}

func TestRunWorkflowUsage(t *testing.T) {
	var out bytes.Buffer
	if err := runWorkflow(t.TempDir(), []string{"workflow"}, &out); err == nil {
		t.Error("a bare workflow command must print usage as an error")
	}
}
