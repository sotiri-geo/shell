package main

import (
	"bytes"
	"testing"
)

func TestRunCommand(t *testing.T) {
	t.Run("WritesCommandStdout", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, dir, "print_message", "#!/bin/sh\necho external\n", 0o755)
		t.Setenv("PATH", dir)

		assertCommandOutput(t, "print_message", "external\n")
	})

	t.Run("ReportsMissingCommand", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())

		var output bytes.Buffer
		found, err := runCommand("invalid_command", &output)
		if err != nil {
			t.Fatalf("runCommand: %v", err)
		}
		if found {
			t.Error("found missing command")
		}
	})

	t.Run("IgnoresExecutableOutsidePath", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, dir, "mytool", "#!/bin/sh\nexit 0\n", 0o755)
		t.Setenv("PATH", t.TempDir())

		var output bytes.Buffer
		found, err := runCommand("mytool", &output)
		if err != nil {
			t.Fatalf("runCommand: %v", err)
		}
		if found {
			t.Error("found executable outside PATH")
		}
	})

	t.Run("PassesArgumentsInOrder", func(t *testing.T) {
		dir := t.TempDir()
		writeFixture(t, dir, "print_args", "#!/bin/sh\nprintf '%s\\n' \"$@\"\n", 0o755)
		t.Setenv("PATH", dir)

		assertCommandOutput(t, "print_args", "hello\nworld\n", "hello", "world")
	})

	t.Run("CatConcatenatesFilesInArgumentOrder", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		writeFixture(t, dir, "file1", "mango chutney.", 0o644)
		writeFixture(t, dir, "file2", "blueberry pineapple.", 0o644)

		assertCommandOutput(t, "cat", "mango chutney.blueberry pineapple.", "file1", "file2")
	})
}

func assertCommandOutput(t *testing.T, name, want string, args ...string) {
	t.Helper()

	var output bytes.Buffer
	found, err := runCommand(name, &output, args...)
	if err != nil {
		t.Fatalf("runCommand(%q, %q): %v", name, args, err)
	}
	if !found {
		t.Fatalf("command %q not found", name)
	}
	if got := output.String(); got != want {
		t.Errorf("command %q output = %q, want %q", name, got, want)
	}
}
