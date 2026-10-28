package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoutingCommand(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"WithExit": {
			input: "exit\n",
			want:  "$ ",
		},
		"WithEcho": {
			input: "echo Hello World\n",
			want:  "$ Hello World\n$ ",
		},
		"WithType": {
			input: "type echo\n",
			want:  "$ echo is a shell builtin\n$ ",
		},
		"WithUnknownCommand": {
			input: "invalid_command\n",
			want:  "$ invalid_command: command not found\n$ ",
		},
		"WithEmptyCommand": {
			input: "",
			want:  "$ ",
		},
		"WithNewLineCommand": {
			input: "\n",
			want:  "$ $ ",
		},
		"WithoutNewLine": {
			input: "echo Hello World",
			want:  "$ Hello World\n",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			err := run(strings.NewReader(tc.input), &output)
			if err != nil {
				t.Fatalf("should not fail: %v", err)
			}

			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Run("WithMultipleCommands", func(t *testing.T) {
		var output bytes.Buffer
		want := "$ Hello World\n$ echo is a shell builtin\n$ "

		err := run(strings.NewReader("echo Hello World\n type echo\n"), &output)
		if err != nil {
			t.Fatalf("should not fail %v", err)
		}

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("WithMultipleArgs", func(t *testing.T) {
		var output bytes.Buffer
		want := "$ echo is a shell builtin\ntype is a shell builtin\n"

		err := run(strings.NewReader("type echo type"), &output)
		if err != nil {
			t.Fatalf("should not fail %v", err)
		}

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("ForExternalCommand", func(t *testing.T) {
		// GIVEN a writer and an external executable command always_false that always returns false
		dir := t.TempDir()
		executable := filepath.Join(dir, "always_false")
		if err := os.WriteFile(
			executable,
			[]byte("#!/bin/sh\necho false\n"),
			PermExecutable,
		); err != nil {
			t.Fatalf("create executable fixture: %v", err)
		}

		t.Setenv("PATH", dir)

		var output bytes.Buffer
		want := "$ false\n"

		// WHEN we call run with
		err := run(strings.NewReader("always_false hello"), &output)
		if err != nil {
			t.Fatalf("should not fail: %v", err)
		}

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
