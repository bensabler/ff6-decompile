package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// This check exists because the readiness matrix disagreed with itself for
// ten units and nothing noticed.
//
// Unit 0 reported "53 requirements, 33 Unknown, 7 Evidence Ready". Counted
// from the tables it wrote in the same commit, the file carried 55 rows, 36
// Unknown and 6 Evidence Ready. The wrong figure then propagated into the
// milestone record, MILESTONES.md, CURRENT_FOCUS.md and the checkpoint chain,
// because every one of them quoted the summary rather than the tables.
//
// That is the failure mode of deviation D1 exactly: two records of one fact,
// disagreeing, with nothing comparing them. The lesson D1 produced was to
// write the comparison down as a test, so this is that comparison.

// readinessPath is the matrix this check reads.
const readinessPath = "docs/demo/DEMO-0001-READINESS.md"

// readinessRow matches a requirement row: a table row whose first cell is a
// section letter and a number, like "| B17 |" or "| X1 |".
var readinessRow = regexp.MustCompile(`^\|\s*([ETBFAX]\d+)\s*\|`)

// readinessStatus matches the first backticked status token in a row. Status
// is not at a fixed column index: the Engine/Text/Battle/Field/Audio sections
// use eleven columns and the cross-cutting section uses six.
var readinessStatus = regexp.MustCompile("`(Unknown|Researching|Evidence Ready|Extractor Ready|Implemented|Integrated|Validated|Blocked|Deferred)`")

// summaryRow matches a summary line, "| `Integrated` | 0 | **14** |".
var summaryRow = regexp.MustCompile("^\\|\\s*`([A-Za-z ]+)`\\s*\\|[^|]*\\|\\s*\\**(\\d+)\\**\\s*\\|")

// totalRow matches "| **Total rows** | **55** | **57** |".
var totalRow = regexp.MustCompile(`^\|\s*\*\*Total rows\*\*\s*\|[^|]*\|\s*\*\*(\d+)\*\*\s*\|`)

// CheckReadinessSummary asserts the readiness matrix's summary table against
// the rows it summarises.
//
// It compares the "Now" column only. The Unit 0 column is a dated historical
// claim about a past commit and is not recomputable from the current file.
func CheckReadinessSummary(root string) ([]Finding, error) {
	const check = "readiness-summary"

	raw, err := os.ReadFile(filepath.Join(root, readinessPath))
	if os.IsNotExist(err) {
		return nil, nil // the matrix is a DEMO-0001 artifact, not a repo invariant
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", readinessPath, err)
	}

	counted := map[string]int{}
	total := 0
	claimed := map[string]int{}
	claimedTotal := -1
	seenIDs := map[string]bool{}
	var findings []Finding

	for _, line := range strings.Split(string(raw), "\n") {
		if m := readinessRow.FindStringSubmatch(line); m != nil {
			id := m[1]
			if seenIDs[id] {
				findings = append(findings, Finding{check,
					fmt.Sprintf("%s: requirement %s appears more than once; every requirement must appear exactly once", readinessPath, id)})
			}
			seenIDs[id] = true
			total++
			if s := readinessStatus.FindStringSubmatch(line); s != nil {
				counted[s[1]]++
			} else {
				findings = append(findings, Finding{check,
					fmt.Sprintf("%s: requirement %s carries no status token", readinessPath, id)})
			}
			continue
		}
		if m := totalRow.FindStringSubmatch(line); m != nil {
			claimedTotal, _ = strconv.Atoi(m[1])
			continue
		}
		if m := summaryRow.FindStringSubmatch(line); m != nil {
			n, err := strconv.Atoi(m[2])
			if err != nil {
				continue
			}
			claimed[strings.TrimSpace(m[1])] = n
		}
	}

	if total == 0 {
		return append(findings, Finding{check,
			fmt.Sprintf("%s: no requirement rows found; the check or the file's shape is wrong", readinessPath)}), nil
	}
	if len(claimed) == 0 {
		return append(findings, Finding{check,
			fmt.Sprintf("%s: no summary table found to compare the %d rows against", readinessPath, total)}), nil
	}

	for _, status := range []string{
		"Validated", "Integrated", "Implemented", "Evidence Ready",
		"Extractor Ready", "Researching", "Blocked", "Deferred", "Unknown",
	} {
		want, ok := claimed[status]
		if !ok {
			continue
		}
		if got := counted[status]; got != want {
			findings = append(findings, Finding{check,
				fmt.Sprintf("%s: summary claims %d %s, the tables carry %d", readinessPath, want, status, got)})
		}
	}
	if claimedTotal >= 0 && claimedTotal != total {
		findings = append(findings, Finding{check,
			fmt.Sprintf("%s: summary claims %d total rows, the tables carry %d", readinessPath, claimedTotal, total)})
	}
	return findings, nil
}
