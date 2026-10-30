package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExternalCommand(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		var output bytes.Buffer
		want := "false\n"
		command := "always_false"
		// GIVEN always_false executable in path p1
		dir := t.TempDir()
		p1 := filepath.Join(dir, command)
		if err := os.WriteFile(
			p1,
			[]byte("#!/bin/sh\necho false\n"),
			PermExecutable,
		); err != nil {
			t.Fatalf("create executable fixture: %v", err)
		}

		t.Setenv("PATH", dir)

		found, err := runCommand(command, &output)

		if !found {
			t.Fatal("should exist")
		}

		if err != nil {
			t.Fatalf("should not fail: %v", err)
		}

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("NotExists", func(t *testing.T) {
		var output bytes.Buffer

		found, err := runCommand("invalid_command", &output)
		if err != nil {
			t.Fatalf("should not error: %v", err)
		}

		if found {
			t.Errorf("found unknown command: %q", "invalid_command")
		}
	})

	t.Run("ExistsNotOnPath", func(t *testing.T) {
		var output bytes.Buffer
		command := "invalid_command"
		dir1 := t.TempDir()
		dir2 := t.TempDir()
		p1 := filepath.Join(dir1, command)
		if err := os.WriteFile(
			p1,
			[]byte("#!/bin/sh\necho invalid_command\n"),
			PermExecutable,
		); err != nil {
			t.Fatalf("create executable fixture: %v", err)
		}

		// Path set with no executable
		t.Setenv("PATH", dir2)

		found, err := runCommand(command, &output)
		if err != nil {
			t.Fatalf("should not fail: %v", err)
		}
		if found {
			t.Errorf("should not exist: %v", command)
		}
	})
	// Check when exists but not on path
	// Check runs with multiple commands
}
