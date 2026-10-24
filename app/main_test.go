package main

import (
	"bytes"
	"strings"
	"testing"
)

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
		"WithCommandsSeparatedByWhitespace": {
			input: "banana apple\n",
			want:  "$ banana: command not found\n$ ", // only first word identified as command
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

func TestEmptyCommandIsIgnored(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"WithEmptyString": {
			input: "",
			want:  "$ ",
		},
		"WithNewLine": {
			input: "\n",
			want:  "$ $ ",
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

func TestRunExit(t *testing.T) {
	input := strings.NewReader("exit\necho should-not-run\n")
	var output bytes.Buffer

	if err := run(input, &output); err != nil {
		t.Fatalf("expected successful exit, got: %v", err)
	}

	// Only the initial prompt should be printed; input after exit must not run.
	if got, want := output.String(), "$ "; got != want {
		t.Errorf("got %q, want %q", got, want)
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
		"WithWhitespaceBefore": {
			input: " echo Hello World\n",
			want:  "$ Hello World\n$ ",
		},
		"WithoutNewline": {
			input: "echo Hello World",
			want:  "$ Hello World\n",
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

func TestTypeBuiltin(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		"WithEcho": {
			// GIVEN a user inputs "type echo"
			input: "type echo\n",
			want:  "$ echo is a shell builtin\n$ ",
		},
		"WithType": {
			// GIVEN a user inputs "type type"
			input: "type type\n",
			want:  "$ type is a shell builtin\n$ ",
		},
		"WithExit": {
			// GIVEN a user inputs "type exit"
			input: "type exit\n",
			want:  "$ exit is a shell builtin\n$ ",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer

			// WHEN run() is called
			err := run(strings.NewReader(tc.input), &output)
			if err != nil {
				t.Fatalf("should not fail: %v", err)
			}

			// THEN command is identified as a builtin
			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTypeNotFound(t *testing.T) {
	testCases := map[string]struct {
		input string
		want  string
	}{
		// GIVEN a user calls type on an unknown command
		"WithUnknownCommand": {
			input: "type invalid_command\n",
			want:  "$ invalid_command: not found\n$ ",
		},
		"WithMultipleUnknownCommands": {
			input: "type banana apple\n",
			want:  "$ banana: not found\napple: not found\n$ ",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer

			// WHEN run is executed
			err := run(strings.NewReader(tc.input), &output)
			if err != nil {
				t.Fatalf("should not fail: %v", err)
			}

			// THEN command will be returned as not found
			if got := output.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
