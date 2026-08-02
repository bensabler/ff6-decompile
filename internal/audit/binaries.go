package audit

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// executableMagic lists the leading bytes of native executable and object
// formats. A compiled binary is recognisable by its header regardless of what
// it is called, which is the property this check needs: the artifact that
// motivated it, a 10 MB `ff6demo`, has no file extension at all.
//
// Shell scripts are unaffected — they begin `#!` — so scripts/ stays tracked.
var executableMagic = []struct {
	magic []byte
	name  string
}{
	{[]byte{0x7F, 'E', 'L', 'F'}, "ELF"},
	{[]byte{0xFE, 0xED, 0xFA, 0xCE}, "Mach-O 32-bit"},
	{[]byte{0xFE, 0xED, 0xFA, 0xCF}, "Mach-O 64-bit"},
	{[]byte{0xCE, 0xFA, 0xED, 0xFE}, "Mach-O 32-bit (LE)"},
	{[]byte{0xCF, 0xFA, 0xED, 0xFE}, "Mach-O 64-bit (LE)"},
	{[]byte{0xCA, 0xFE, 0xBA, 0xBE}, "Mach-O universal or Java class"},
	{[]byte{'M', 'Z'}, "PE/DOS"},
}

// CheckTrackedBinaries rejects compiled executables in the git index.
//
// It exists because of a real defect: `ff6demo` — a 10 MB Mach-O binary — was
// committed on 2026-08-02 and reached the remote. `.gitignore` covered
// `/ff6lab` but nobody extended it when the second executable appeared, and
// neither existing gate could catch it. CheckRestrictedTracked matches file
// *extensions*, and a build output has none; CheckRequiredTracked only asserts
// that certain files are present, never that others are absent.
//
// Extension lists protect against the files you thought of. This checks what a
// file actually is.
func CheckTrackedBinaries(root string) ([]Finding, error) {
	files, err := gitLsFiles(root)
	if err != nil {
		return []Finding{{"tracked-binaries", "git unavailable: " + err.Error()}}, nil
	}

	var fs []Finding
	buf := make([]byte, 4)
	for _, rel := range files {
		f, err := os.Open(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			// A tracked file missing from the worktree is a different
			// problem, and not this check's to report.
			continue
		}
		n, _ := f.Read(buf)
		f.Close()
		if n < 2 {
			continue
		}
		head := buf[:n]
		for _, m := range executableMagic {
			if len(m.magic) <= n && bytes.HasPrefix(head, m.magic) {
				fs = append(fs, Finding{"tracked-binaries", fmt.Sprintf(
					"tracked %s executable: %s — build outputs belong in .gitignore, not the repository", m.name, rel)})
				break
			}
		}
	}
	return fs, nil
}
