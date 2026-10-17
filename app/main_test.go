package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRunUnknownCommand(t *testing.T) {
	input := strings.NewReader("banana\n")
	var output bytes.Buffer

	err := run(input, &output)
	if err != nil && err != io.EOF {
		t.Fatalf("should not error: %v", err)
	}

	want := "$ banana: command not found\n$ "

	// Note: This is idiomatic in go to assign var inline and
	// use the variable assigned within an expression
	if got := output.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
