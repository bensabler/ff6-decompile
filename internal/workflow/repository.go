package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type repositoryContext struct {
	Root     string
	Identity string
	Branch   string
	Head     string
}

// resolveRepository records a canonical Git identity rather than trusting the
// checkout directory name. A normalized remote is preferred. A repository
// without a remote is identified honestly by its canonical common Git dir,
// which is stable across linked worktrees but not falsely portable to clones.
func resolveRepository(root string) (repositoryContext, error) {
	top, err := gitOutput(root, "rev-parse", "--show-toplevel")
	if err != nil {
		return repositoryContext{}, fmt.Errorf("resolve repository root: %w", err)
	}
	top, err = canonicalPath(top)
	if err != nil {
		return repositoryContext{}, fmt.Errorf("canonicalize repository root: %w", err)
	}

	identity, err := resolveRepositoryIdentity(top)
	if err != nil {
		return repositoryContext{}, err
	}
	branch, err := gitOutput(top, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			return repositoryContext{}, fmt.Errorf("resolve starting branch: %w", err)
		}
		branch = "(detached)"
	}
	head, err := gitOutput(top, "rev-parse", "HEAD")
	if err != nil {
		return repositoryContext{}, fmt.Errorf("resolve starting HEAD: %w", err)
	}
	return repositoryContext{Root: top, Identity: identity, Branch: branch, Head: head}, nil
}

func resolveRepositoryIdentity(root string) (string, error) {
	remote, err := gitOutput(root, "config", "--get", "remote.origin.url")
	if err == nil {
		return "remote:" + canonicalRemote(remote, root), nil
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		return "", fmt.Errorf("resolve origin remote: %w", err)
	}

	remotes, err := gitOutput(root, "remote")
	if err != nil {
		return "", fmt.Errorf("list repository remotes: %w", err)
	}
	var identities []string
	for _, name := range strings.Fields(remotes) {
		raw, err := gitOutput(root, "remote", "get-url", name)
		if err != nil {
			return "", fmt.Errorf("resolve remote %s: %w", name, err)
		}
		identities = append(identities, canonicalRemote(raw, root))
	}
	if len(identities) > 0 {
		sort.Strings(identities)
		sum := sha256.Sum256([]byte(strings.Join(identities, "\n")))
		return "remote-set-sha256:" + hex.EncodeToString(sum[:]), nil
	}

	common, err := gitOutput(root, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve local Git identity: %w", err)
	}
	if !filepath.IsAbs(common) {
		common = filepath.Join(root, common)
	}
	common, err = canonicalPath(common)
	if err != nil {
		return "", fmt.Errorf("canonicalize common Git dir: %w", err)
	}
	return "local-git-common-dir:" + common, nil
}

var scpRemote = regexp.MustCompile(`^(?:[^@/]+@)?([^:/]+):(.+)$`)

// canonicalRemote strips URL credentials, queries, fragments, trailing .git,
// and spelling differences that do not identify a different repository.
func canonicalRemote(raw, root string) string {
	raw = strings.TrimSpace(raw)
	if m := scpRemote.FindStringSubmatch(raw); m != nil && !strings.Contains(raw, "://") {
		return "ssh://" + strings.ToLower(m[1]) + "/" + trimRemotePath(m[2])
	}
	if u, err := url.Parse(raw); err == nil && u.Scheme != "" {
		u.User = nil
		u.RawQuery = ""
		u.Fragment = ""
		u.Scheme = strings.ToLower(u.Scheme)
		u.Host = strings.ToLower(u.Host)
		u.Path = "/" + trimRemotePath(u.Path)
		return u.String()
	}
	if !filepath.IsAbs(raw) {
		raw = filepath.Join(root, raw)
	}
	if p, err := canonicalPath(raw); err == nil {
		return "file://" + trimRemotePath(filepath.ToSlash(p))
	}
	return "unresolved-local:" + trimRemotePath(filepath.ToSlash(raw))
}

func trimRemotePath(path string) string {
	path = strings.Trim(strings.TrimSpace(path), "/")
	return strings.TrimSuffix(path, ".git")
}

func canonicalPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func gitOutput(root string, args ...string) (string, error) {
	cmdArgs := append([]string{"-C", root}, args...)
	b, err := exec.Command("git", cmdArgs...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
