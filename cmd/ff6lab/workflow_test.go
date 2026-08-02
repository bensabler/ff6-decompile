package main

import (
	"bytes"
	"os"
	"os/exec"
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

func workflowFixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runWorkflowGit(t, root, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runWorkflowGit(t, root, "add", "README.md")
	runWorkflowGit(t, root, "-c", "user.name=Workflow Test", "-c", "user.email=workflow@example.invalid",
		"commit", "-q", "-m", "fixture")
	return root
}

func runWorkflowGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmdArgs := append([]string{"-C", root}, args...)
	if out, err := exec.Command("git", cmdArgs...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func TestRunWorkflowLifecycle(t *testing.T) {
	root := workflowFixtureRoot(t)
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

// close must ignore provider-global transcript history and exit non-zero when
// the run ledger does not satisfy the contract.
func TestRunWorkflowCloseDoesNotScanGlobalClaudeTranscripts(t *testing.T) {
	root := workflowFixtureRoot(t)
	path := writeContract(t, root, contractJSON)
	var out bytes.Buffer
	if err := runWorkflow(root, []string{"workflow", "start", path}, &out); err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(root, "home")
	transcriptDir := filepath.Join(home, ".claude", "projects")
	if err := os.MkdirAll(transcriptDir, 0o755); err != nil {
		t.Fatal(err)
	}
	transcript := `{"timestamp":"2026-08-02T00:00:01Z","message":{"content":[` +
		`{"type":"tool_use","name":"Agent","input":{"subagent_type":"graphics-researcher"}}]}}` + "\n"
	if err := os.WriteFile(filepath.Join(transcriptDir, "historical.jsonl"), []byte(transcript), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
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
	root := workflowFixtureRoot(t)
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
