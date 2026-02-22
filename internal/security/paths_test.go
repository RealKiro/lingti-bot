package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathChecker_NoRestrictions(t *testing.T) {
	pc := NewPathChecker(nil)
	if pc.HasRestrictions() {
		t.Fatal("expected no restrictions")
	}
	if !pc.IsAllowed("/any/path") {
		t.Fatal("expected all paths allowed when no restrictions")
	}
}

func TestPathChecker_AllowedPaths(t *testing.T) {
	dir := t.TempDir()
	pc := NewPathChecker([]string{dir})

	if !pc.IsAllowed(filepath.Join(dir, "subdir", "file.txt")) {
		t.Error("path inside allowed dir should be allowed")
	}
	if pc.IsAllowed("/some/other/path") {
		t.Error("path outside allowed dir should be blocked")
	}
}

func TestPathChecker_TildeExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	pc := NewPathChecker([]string{"~/testdir"})
	expected := filepath.Join(home, "testdir")
	if len(pc.AllowedPaths()) != 1 || pc.AllowedPaths()[0] != expected {
		t.Errorf("expected %s, got %v", expected, pc.AllowedPaths())
	}
}

func TestPathChecker_TraversalAttempt(t *testing.T) {
	dir := t.TempDir()
	pc := NewPathChecker([]string{dir})

	traversal := filepath.Join(dir, "..", "..", "etc", "passwd")
	if pc.IsAllowed(traversal) {
		t.Error("traversal path should be blocked")
	}
}

func TestPathChecker_ExactMatch(t *testing.T) {
	dir := t.TempDir()
	pc := NewPathChecker([]string{dir})

	if !pc.IsAllowed(dir) {
		t.Error("exact allowed dir should be allowed")
	}
}

func TestPathChecker_CheckPath(t *testing.T) {
	dir := t.TempDir()
	pc := NewPathChecker([]string{dir})

	if err := pc.CheckPath(filepath.Join(dir, "ok.txt")); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if err := pc.CheckPath("/blocked/path"); err == nil {
		t.Error("expected error for blocked path")
	}
}
