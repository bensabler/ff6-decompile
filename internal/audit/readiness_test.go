package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readinessFixture writes a minimal readiness matrix with the given rows and summary.
func readinessFixture(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "docs", "demo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DEMO-0001-READINESS.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

const agreeing = "" +
	"# matrix\n\n" +
	"| # | Requirement | Status | Next |\n|---|---|---|---|\n" +
	"| E1 | a | `Integrated` | — |\n" +
	"| E2 | b | `Integrated` | — |\n" +
	"| T1 | c | `Validated` | — |\n" +
	"| F1 | d | `Unknown` | — |\n" +
	"| X1 | e | `Unknown` | — |\n\n" +
	"## Summary\n\n" +
	"| Status | Unit 0 | Now |\n|---|---|---|\n" +
	"| `Validated` | 0 | **1** |\n" +
	"| `Integrated` | 0 | **2** |\n" +
	"| `Unknown` | 5 | **2** |\n" +
	"| **Total rows** | **5** | **5** |\n"

func TestReadinessSummaryAcceptsAgreement(t *testing.T) {
	got, err := CheckReadinessSummary(readinessFixture(t, agreeing))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("a self-consistent matrix produced findings: %v", got)
	}
}

// The regression case. These are the exact shapes of the drift that went
// unnoticed for ten units.
func TestReadinessSummaryCatchesDrift(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "a status count is wrong",
			body: strings.Replace(agreeing, "| `Integrated` | 0 | **2** |", "| `Integrated` | 0 | **7** |", 1),
			want: "summary claims 7 Integrated, the tables carry 2",
		},
		{
			name: "the total is wrong",
			body: strings.Replace(agreeing, "| **Total rows** | **5** | **5** |", "| **Total rows** | **5** | **53** |", 1),
			want: "summary claims 53 total rows, the tables carry 5",
		},
		{
			name: "a row was added without updating the summary",
			body: strings.Replace(agreeing, "| X1 | e | `Unknown` | — |", "| X1 | e | `Unknown` | — |\n| X2 | f | `Unknown` | — |", 1),
			want: "summary claims 2 Unknown, the tables carry 3",
		},
		{
			name: "a row carries a status from another vocabulary",
			body: strings.Replace(agreeing, "| F1 | d | `Unknown` | — |", "| F1 | d | `Located` | — |", 1),
			want: "requirement F1 carries no status token",
		},
		{
			name: "a requirement id appears twice",
			body: strings.Replace(agreeing, "| X1 | e | `Unknown` | — |", "| E1 | e | `Unknown` | — |", 1),
			want: "requirement E1 appears more than once",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckReadinessSummary(readinessFixture(t, tt.body))
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range got {
				if strings.Contains(f.Message, tt.want) {
					return
				}
			}
			t.Errorf("want a finding containing %q, got %v", tt.want, got)
		})
	}
}

func TestReadinessSummaryHandlesMissingOrEmptyFiles(t *testing.T) {
	// Absent file: the matrix is a DEMO-0001 artifact, not a repo invariant,
	// so its absence is not a finding.
	got, err := CheckReadinessSummary(t.TempDir())
	if err != nil {
		t.Fatalf("an absent matrix should not error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("an absent matrix should produce no findings, got %v", got)
	}

	// Present but shapeless: that IS a finding, because it means the check
	// has silently stopped checking.
	got, err = CheckReadinessSummary(readinessFixture(t, "# matrix\n\nnothing here\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Error("a matrix with no requirement rows should be reported, not silently passed")
	}
}

// The real file must pass. This is the assertion that would have failed on
// 2026-08-02 before Unit 10 recounted it.
func TestRealReadinessMatrixAgreesWithItself(t *testing.T) {
	got, err := CheckReadinessSummary("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range got {
		t.Errorf("%s", f)
	}
}
