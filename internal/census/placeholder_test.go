package census

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Fields whose values are ordinary short names/labels. Under the revised
// asset policy these are permitted in Git, so a placeholder in one of them
// is a policy misuse — the policy exists to keep substantial assets out of
// the repository, not to hide ordinary reconstruction data.
var nameFields = map[string]bool{
	"name": true, "display_name": true, "label": true, "title": true,
}

// placeholderMarkers are the evasions this test exists to prevent.
var placeholderMarkers = []string{
	"see local extraction",
	"asset policy",
	"redacted",
	"withheld",
	"local only",
	"local-only",
}

// TestNoPlaceholdersInPermittedNameFields walks every tracked data and
// manifest JSON file and fails if a permitted name field carries a
// placeholder instead of its confirmed value.
//
// Genuinely unresolved values are still allowed to say so: the empty
// string, "unknown", and "unnamed" are accepted, because "we have not
// established this" is an honest state and inventing a value would be
// worse. What is rejected is a *known* value replaced by a policy excuse.
func TestNoPlaceholdersInPermittedNameFields(t *testing.T) {
	root := repoRoot(t)
	dirs := []string{"data", "manifests"}

	var offenders []string
	for _, dir := range dirs {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Ext(p) != ".json" {
				return nil //nolint:nilerr // missing optional dirs are not failures
			}
			b, rerr := os.ReadFile(p)
			if rerr != nil {
				return rerr
			}
			var v any
			if jerr := json.Unmarshal(b, &v); jerr != nil {
				return nil // schema validity is a different gate's job
			}
			rel, _ := filepath.Rel(root, p)
			walkJSON(v, "", func(path, key, val string) {
				if !nameFields[key] {
					return
				}
				low := strings.ToLower(val)
				for _, m := range placeholderMarkers {
					if strings.Contains(low, m) {
						offenders = append(offenders, rel+" "+path+": "+val)
						return
					}
				}
			})
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", dir, err)
		}
	}

	if len(offenders) > 0 {
		t.Fatalf("placeholder(s) in policy-permitted name fields (restore the confirmed value; see docs/legal/ASSET_POLICY.md):\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

// TestSpellInventoryIsComplete guards the specific regression that
// motivated this test: all 54 confirmed spell names present and distinct.
func TestSpellInventoryIsComplete(t *testing.T) {
	root := repoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "data", "census", "spells.json"))
	if err != nil {
		t.Fatalf("reading spell inventory: %v", err)
	}
	var doc struct {
		Records []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"records"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatalf("parsing spell inventory: %v", err)
	}
	if len(doc.Records) != 54 {
		t.Fatalf("spell inventory has %d records, want 54 (boundary Confirmed in EXP-0027)", len(doc.Records))
	}
	seen := map[string]int{}
	for _, r := range doc.Records {
		if strings.TrimSpace(r.Name) == "" {
			t.Errorf("spell %d has an empty name", r.ID)
			continue
		}
		if prev, dup := seen[r.Name]; dup {
			t.Errorf("spell %d duplicates the name of spell %d (%q) — a decode error, not a real collision", r.ID, prev, r.Name)
		}
		seen[r.Name] = r.ID
	}
}

// walkJSON visits every string leaf, reporting its JSON path, the object
// key it sits under, and its value.
func walkJSON(v any, path string, fn func(path, key, val string)) {
	switch t := v.(type) {
	case map[string]any:
		for k, child := range t {
			p := path + "." + k
			if s, ok := child.(string); ok {
				fn(p, k, s)
				continue
			}
			walkJSON(child, p, fn)
		}
	case []any:
		for i, child := range t {
			walkJSON(child, path+"["+itoa(i)+"]", fn)
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

// repoRoot walks up from the test's working directory to the module root.
func repoRoot(t *testing.T) string {
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
