package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeGo creates a Go file inside a fixture tree.
func writeGo(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckImportBoundariesFlagsROMInTheDemo(t *testing.T) {
	root := t.TempDir()
	writeGo(t, root, "internal/content/font.go", `package content

import (
	"fmt"

	"`+modulePath+`/internal/rom"
)

var _ = fmt.Sprint
var _ = rom.Size
`)
	findings, err := CheckImportBoundaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %v, want exactly 1", findings)
	}
	if !strings.Contains(findings[0].Message, "internal/rom") {
		t.Errorf("finding should name the forbidden import: %s", findings[0].Message)
	}
	if !strings.Contains(findings[0].Message, "->") {
		t.Errorf("finding should show the import path: %s", findings[0].Message)
	}
	// The message must explain the boundary, not merely assert it.
	if !strings.Contains(findings[0].Message, "ASSET_POLICY") {
		t.Errorf("finding should explain why: %s", findings[0].Message)
	}
}

func TestCheckImportBoundariesAllowsTestFiles(t *testing.T) {
	root := t.TempDir()
	// The archive-vs-ROM differential is exactly this shape and must be
	// allowed: the boundary is about what ships in a binary.
	writeGo(t, root, "internal/content/rom_differential_test.go", `package content_test

import _ "`+modulePath+`/internal/rom"
`)
	findings, err := CheckImportBoundaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("test files must be exempt; got %v", findings)
	}
}

func TestCheckImportBoundariesAllowsOtherPackages(t *testing.T) {
	root := t.TempDir()
	// internal/extract is *supposed* to read the ROM.
	writeGo(t, root, "internal/extract/extractors.go", `package extract

import _ "`+modulePath+`/internal/rom"
`)
	writeGo(t, root, "cmd/ff6lab/main.go", `package main

import _ "`+modulePath+`/internal/rom"
`)
	findings, err := CheckImportBoundaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("packages outside the rule's scope must be unaffected; got %v", findings)
	}
}

func TestCheckImportBoundariesFlagsGameInPlatform(t *testing.T) {
	root := t.TempDir()
	writeGo(t, root, "internal/platform/snesaddr/snesaddr.go", `package snesaddr

import _ "`+modulePath+`/internal/game/battle"
`)
	findings, err := CheckImportBoundaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %v, want exactly 1", findings)
	}
	if !strings.Contains(findings[0].Message, "ARCHITECTURE.md") {
		t.Errorf("finding should cite the layering document: %s", findings[0].Message)
	}
}

// TestCheckImportBoundariesIsTransitive is the case a direct-import check
// misses, and the one that actually occurred: internal/content reached
// internal/rom through internal/extract, and nothing noticed.
func TestCheckImportBoundariesIsTransitive(t *testing.T) {
	root := t.TempDir()
	writeGo(t, root, "internal/content/archive.go", `package content

import _ "`+modulePath+`/internal/extract"
`)
	writeGo(t, root, "internal/extract/extract.go", `package extract

import _ "`+modulePath+`/internal/rom"
`)
	writeGo(t, root, "internal/rom/rom.go", "package rom\n")

	findings, err := CheckImportBoundaries(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %v, want exactly 1 transitive violation", findings)
	}
	want := "internal/content -> internal/extract -> internal/rom"
	if !strings.Contains(findings[0].Message, want) {
		t.Errorf("finding = %q, want it to contain %q", findings[0].Message, want)
	}
}

// TestPrefixMatchingIsPathAware guards against a substring match flagging a
// package that merely starts with the same letters.
func TestPrefixMatchingIsPathAware(t *testing.T) {
	if hasPathPrefix(modulePath+"/internal/romhacking", modulePath+"/internal/rom") {
		t.Error("internal/romhacking must not match the internal/rom rule")
	}
	if !hasPathPrefix(modulePath+"/internal/rom", modulePath+"/internal/rom") {
		t.Error("the package itself must match")
	}
	if !hasPathPrefix(modulePath+"/internal/rom/sub", modulePath+"/internal/rom") {
		t.Error("a subpackage must match")
	}
}

// TestRealRepositoryHonorsItsBoundaries runs the check against the actual
// tree, so the rules are load-bearing rather than decorative.
func TestRealRepositoryHonorsItsBoundaries(t *testing.T) {
	findings, err := CheckImportBoundaries(repoRootForTest(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		t.Errorf("%s", f)
	}
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not locate module root")
	return ""
}
