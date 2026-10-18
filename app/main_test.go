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
}

func TestUnknownCommand(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"WithSingleUnknownCommand": {
			input: "banana\n",
			want:  "$ banana: command not found\n$ ",
		},
		"WithMultipleUnknownCommands": {
			input: "banana\napple\n",
			want:  "$ banana: command not found\n$ apple: command not found\n$ ",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer

			// WHEN executing run
			err := run(strings.NewReader(tc.input), &output)
			// THEN no errors AND expected output is written
			if err != nil {
				t.Fatalf("should not fail: %v", err)
			}
			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEchoCommand(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"WithArgs": {
			input: "echo Hello World\n",
			want:  "$ Hello World\n$ ",
		},
		"WithoutArgs": {
			input: "echo\n",
			want:  "$ \n$ ",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer

			// WHEN executing run
			err := run(strings.NewReader(tc.input), &output)
			// THEN no errors AND expected output is written
			if err != nil {
				t.Fatalf("should not fail: %v", err)
			}
			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
