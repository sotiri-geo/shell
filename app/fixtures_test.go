package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, dir, name, contents string, mode os.FileMode) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write fixture %q: %v", path, err)
	}
	// Chmod sets the requested mode independently of the process umask.
	if err := os.Chmod(path, mode); err != nil {
		t.Fatalf("set fixture mode for %q: %v", path, err)
	}
	return path
}

func resolvedTempDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("resolve temporary directory %q: %v", dir, err)
	}
	return resolved
}
