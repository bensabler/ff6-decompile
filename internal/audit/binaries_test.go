package audit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gitFixture builds a tiny repository so CheckTrackedBinaries has an index to
// read. The check consults `git ls-files`, so a bare directory is not enough.
func gitFixture(t *testing.T, files map[string][]byte) string {
	t.Helper()
	root := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "test@example.invalid")
	run("config", "user.name", "test")
	for name, body := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("add", "-A")
	return root
}

func TestCheckTrackedBinariesFlagsExecutables(t *testing.T) {
	tests := []struct {
		name  string
		head  []byte
		which string
	}{
		// The exact case that occurred: a Go binary with no extension.
		{"mach-o 64-bit", []byte{0xCF, 0xFA, 0xED, 0xFE, 0x07, 0x00, 0x00, 0x01}, "Mach-O"},
		{"elf", []byte{0x7F, 'E', 'L', 'F', 0x02, 0x01, 0x01, 0x00}, "ELF"},
		{"pe", []byte{'M', 'Z', 0x90, 0x00}, "PE"},
		{"universal", []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x02}, "universal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := gitFixture(t, map[string][]byte{"ff6demo": tt.head})
			findings, err := CheckTrackedBinaries(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(findings) != 1 {
				t.Fatalf("findings = %v, want exactly 1", findings)
			}
			if !strings.Contains(findings[0].Message, "ff6demo") {
				t.Errorf("finding should name the file: %s", findings[0].Message)
			}
			if !strings.Contains(findings[0].Message, ".gitignore") {
				t.Errorf("finding should say what to do about it: %s", findings[0].Message)
			}
		})
	}
}

// TestCheckTrackedBinariesAllowsRealSources is the false-positive guard. Shell
// scripts are executable and legitimately tracked; so is every text file in
// the repository.
func TestCheckTrackedBinariesAllowsRealSources(t *testing.T) {
	root := gitFixture(t, map[string][]byte{
		"scripts/check.sh": []byte("#!/usr/bin/env bash\nset -eu\n"),
		"main.go":          []byte("package main\n\nfunc main() {}\n"),
		"README.md":        []byte("# hello\n"),
		"data/x.json":      []byte("{\"a\":1}\n"),
		"empty":            {},
		"one-byte":         []byte("x"),
		// A PNG is caught by CheckRestrictedTracked, not this one; it must
		// not be double-reported here.
		"img.png": []byte{0x89, 'P', 'N', 'G'},
	})
	findings, err := CheckTrackedBinaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("no source file should be flagged; got %v", findings)
	}
}

// TestRealRepositoryTracksNoBinaries runs the check against the actual tree.
// This is the assertion that would have failed on 2026-08-02.
func TestRealRepositoryTracksNoBinaries(t *testing.T) {
	findings, err := CheckTrackedBinaries(repoRootForTest(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		t.Errorf("%s", f)
	}
}
