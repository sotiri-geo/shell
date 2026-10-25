package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	var output bytes.Buffer
	want := "hello world\n"

	err := handleEcho("hello world", &output)
	if err != nil {
		t.Fatalf("should not error: %v", err)
	}

	if got := output.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestType(t *testing.T) {
	testCase := map[string]struct {
		input string
		want  string
	}{
		"WithBuitinCommandEcho": {
			input: "echo",
			want:  "echo is a shell builtin\n",
		},
		"WithBuitinCommandExit": {
			input: "exit",
			want:  "exit is a shell builtin\n",
		},
		"WithBuitinCommandType": {
			input: "type",
			want:  "type is a shell builtin\n",
		},
		"WithMulitpleBuitinCommandType": {
			input: "type echo",
			want:  "type is a shell builtin\necho is a shell builtin\n",
		},
		"WithUnknownCommand": {
			input: "invalid_command",
			want:  "invalid_command: not found\n",
		},
		"WithKnownAndUnknownCommand": {
			input: "echo invalid_command",
			want:  "echo is a shell builtin\ninvalid_command: not found\n",
		},
		"WithUnknownAndKnownCommand": {
			input: "invalid_command echo",
			want:  "invalid_command: not found\necho is a shell builtin\n",
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer

			err := handleType(tc.input, &output)
			if err != nil {
				t.Fatalf("should not error: %v", err)
			}

			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTypeFindsExecutable(t *testing.T) {
	// GIVEN a directory d1 on PATH with an executable mytool
	dir := t.TempDir()
	executable := filepath.Join(dir, "mytool")

	if err := os.WriteFile(
		executable,
		[]byte("#!/bin/sh\nexit 0\n"),
		PermExecutable,
	); err != nil {
		t.Fatalf("create executable fixture: %v", err)
	}

	t.Setenv("PATH", dir)

	var output bytes.Buffer
	// WHEN calling handleType with mytool
	if err := handleType("mytool", &output); err != nil {
		t.Fatalf("handleType: %v", err)
	}

	want := "mytool is " + executable + "\n"
	// THEN mytool
	if got := output.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFindExecutable(t *testing.T) {
	t.Run("InFirstDir", func(t *testing.T) {
		// GIVEN the first directory d1 in path containing executable e1
		dir := t.TempDir()
		executable := filepath.Join(dir, "mytool")
		if err := os.WriteFile(
			executable,
			[]byte("#!/bin/sh\nexit 0\n"),
			PermExecutable,
		); err != nil {
			t.Fatalf("create executable fixture: %v", err)
		}

		t.Setenv("PATH", dir)
		// WHEN calling findExecutable with e1 and PATH
		got, found := findExecutable("mytool", "PATH")
		// THEN we return the full path and found as true
		if !found {
			t.Fatal("executable mytool should exist")
		}
		if got != executable {
			t.Errorf("got %q, want %q", got, executable)
		}
	})

	t.Run("InSecondDir", func(t *testing.T) {
		// GIVEN a directory d1 located second in path with executable e1
		// e.g. "file/to/path1:d1"
		dirDummy := t.TempDir()
		dir := t.TempDir()
		path := strings.Join([]string{dirDummy, dir}, ":")
		executable := filepath.Join(dir, "mytool")
		if err := os.WriteFile(
			executable,
			[]byte("#!/bin/sh\nexit 0\n"),
			PermExecutable,
		); err != nil {
			t.Fatalf("create executable fixture: %v", err)
		}

		t.Setenv("PATH", path)
		// WHEN calling findExecutable with e1 and PATH
		got, found := findExecutable("mytool", "PATH")
		// THEN we return the full path and found as true
		if !found {
			t.Fatal("executable mytool should exist")
		}
		if got != executable {
			t.Errorf("got %q, want %q", got, executable)
		}
	})
}
