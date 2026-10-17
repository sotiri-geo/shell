package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
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
}
