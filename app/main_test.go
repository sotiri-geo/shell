package main

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("WithSingleUnknownCommand", func(t *testing.T) {
		input := strings.NewReader("banana\n")
		var output bytes.Buffer

		err := run(input, &output)
		if err != nil && err != io.EOF {
			t.Fatalf("should not error: %v", err)
		}

		want := "$ banana: command not found\n$ "

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("WithMultipleUnknownCommands", func(t *testing.T) {
		input := strings.NewReader("banana\napple\n")
		var output bytes.Buffer

		err := run(input, &output)
		if err != nil && err != io.EOF {
			t.Fatalf("should not error: %v", err)
		}

		want := "$ banana: command not found\n$ apple: command not found\n$ "

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("WithExitCommand", func(t *testing.T) {
		// Build the shell into a temporary directory.
		binary := filepath.Join(t.TempDir(), "shell")
		build := exec.Command("go", "build", "-o", binary, ".")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build failed: %v\n%s", err, output)
		}

		// Fail if the shell takes too long to exit.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, binary)
		cmd.Stdin = strings.NewReader("exit\nshould-not-run\n")

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("expected exit status 0: %v\n%s", err, output)
		}

		// The shell should stop after exit, displaying only its first prompt.
		if got := string(output); got != "$ " {
			t.Errorf("got output %q, want %q", got, "$ ")
		}
	})

	t.Run("WithEchoCommandAndArgs", func(t *testing.T) {
		// Must print to stdout
		// but for testing just write to writer
		input := strings.NewReader("echo Hello World\n")
		var output bytes.Buffer

		err := run(input, &output)
		if err != nil {
			t.Fatalf("should not fail: %v", err)
		}
		// In tests we need to consider that we write the prompt to the writer
		// so well either need to just include the prompt "$ " or transform before
		// we assert. Well just go with the former and include the prompt string.

		want := "$ Hello World\n"

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("WithEchoCommandAndEmptyArg", func(t *testing.T) {
		// Tests that running echo with an empty arg, writes just a new line to
		// the writer.

		input := strings.NewReader("echo\n")
		var output bytes.Buffer

		err := run(input, &output)
		if err != nil {
			t.Fatalf("should not fail: %v", err)
		}

		want := "$ \n$ "

		if got := output.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
