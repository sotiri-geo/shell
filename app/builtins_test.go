package main

import (
	"bytes"
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
