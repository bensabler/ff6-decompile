package audit

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// modulePath is this repository's Go module path, the prefix every internal
// import carries.
const modulePath = "github.com/bensabler/ff6-decompile"

// importRule forbids a package prefix from being imported by anything under
// the listed package prefixes, except by the named exemptions.
type importRule struct {
	// Forbidden is the import path prefix that must not appear.
	Forbidden string
	// Within lists the package-path prefixes the rule applies to.
	Within []string
	// Except lists package-path prefixes allowed to import it anyway.
	Except []string
	// Why is included in the finding, so the message explains the boundary
	// rather than just asserting it.
	Why string
}

// importRules are the architectural boundaries that are cheap to state and
// expensive to rediscover. They are checked mechanically because a comment
// asking people not to import something does not survive contact with a
// deadline.
//
// Test files are exempt from every rule: the boundary is about what ships in
// a binary. The archive-vs-ROM differential in internal/content is exactly
// the case — it must read both sides, and lives in an external test package
// so the production package stays clean.
var importRules = []importRule{
	{
		Forbidden: modulePath + "/internal/rom",
		Within: []string{
			modulePath + "/cmd/ff6demo",
			modulePath + "/internal/content",
			modulePath + "/internal/engine",
		},
		Why: "the demo loads game data only from the generated archive, so the binary structurally cannot read a ROM (docs/legal/ASSET_POLICY.md)",
	},
	{
		Forbidden: modulePath + "/internal/game",
		Within:    []string{modulePath + "/internal/platform"},
		Why:       "internal/platform models SNES hardware and must not depend on FF6-specific packages (ARCHITECTURE.md)",
	},
}

// CheckImportBoundaries reports imports that cross an architectural boundary,
// following the module's own import graph transitively.
//
// Transitivity is the whole point. "The demo binary cannot read a ROM" is a
// claim about what gets linked, and a direct-import check does not establish
// it: internal/content reached internal/rom through internal/extract until
// the manifest model was split into internal/assetmanifest, and a
// direct-import check called that clean.
//
// It uses go/parser rather than `go list`, so it needs no build cache, works
// against a fixture tree, and cannot be defeated by a build tag.
func CheckImportBoundaries(root string) ([]Finding, error) {
	graph, err := buildImportGraph(root)
	if err != nil {
		return nil, err
	}

	var findings []Finding
	for _, r := range importRules {
		for pkg := range graph {
			if !anyPrefix(pkg, r.Within) || anyPrefix(pkg, r.Except) {
				continue
			}
			if path := findForbidden(graph, pkg, r); path != "" {
				findings = append(findings, Finding{
					Check:   "import-boundaries",
					Message: path + " — " + r.Why,
				})
			}
		}
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].Message < findings[j].Message })
	return findings, nil
}

// buildImportGraph maps each in-module package path to the in-module packages
// it imports. Test files are excluded: the boundary is about what ships in a
// binary, and a differential test legitimately reads both sides.
func buildImportGraph(root string) (map[string][]string, error) {
	graph := map[string][]string{}
	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "local_artifacts", "node_modules", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		pkg := modulePath + "/" + filepath.ToSlash(filepath.Dir(rel))

		f, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			// A file that does not parse is a build problem, which the build
			// gate reports better than this check could.
			return nil
		}
		if _, ok := graph[pkg]; !ok {
			graph[pkg] = nil
		}
		for _, spec := range f.Imports {
			imp, uerr := strconv.Unquote(spec.Path.Value)
			if uerr != nil || !hasPathPrefix(imp, modulePath) {
				continue
			}
			graph[pkg] = append(graph[pkg], imp)
		}
		return nil
	})
	return graph, err
}

// findForbidden walks the import graph from pkg and returns a human-readable
// path to the first forbidden package, or "" if none is reachable.
func findForbidden(graph map[string][]string, pkg string, r importRule) string {
	type node struct {
		pkg   string
		trail []string
	}
	seen := map[string]bool{pkg: true}
	queue := []node{{pkg, []string{short(pkg)}}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		imports := append([]string(nil), graph[cur.pkg]...)
		sort.Strings(imports) // deterministic reporting
		for _, imp := range imports {
			if hasPathPrefix(imp, r.Forbidden) {
				return strings.Join(append(cur.trail, short(imp)), " -> ")
			}
			if seen[imp] {
				continue
			}
			seen[imp] = true
			queue = append(queue, node{imp, append(append([]string(nil), cur.trail...), short(imp))})
		}
	}
	return ""
}

// short trims the module prefix so findings read as repository paths.
func short(pkg string) string { return strings.TrimPrefix(pkg, modulePath+"/") }

// hasPathPrefix reports whether importPath is prefix or a package under it,
// so "…/internal/rom" matches but "…/internal/romhacking" does not.
func hasPathPrefix(importPath, prefix string) bool {
	return importPath == prefix || strings.HasPrefix(importPath, prefix+"/")
}

func anyPrefix(pkg string, prefixes []string) bool {
	for _, p := range prefixes {
		if hasPathPrefix(pkg, p) {
			return true
		}
	}
	return false
}
