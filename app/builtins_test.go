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
		"WithUnKnownAndUnknownCommand": {
			input: "invalid_command banana",
			want:  "invalid_command: not found\nbanana: not found\n",
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

func TestTypeReturnsExecutable(t *testing.T) {
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
	testCases := map[string]struct {
		fileName         string
		executableExists bool
		setup            func() string
	}{
		"InFirstDirectory": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
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
				return executable
			},
		},
		"InSecondDirectory": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
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
				return executable
			},
		},
		// This verifies when we have two files named the same
		// under two different directories in PATH, we return the
		// first one that is executable
		"WithSameFilenameOnlyOneExecutable": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
				// GIVEN a directory d1 with filename f1 with non-exec permissions
				// GIVEN a directory d2 with filename f1 with exec permissions
				// GIVEN an env PATH d1:d2
				d1 := t.TempDir()
				d2 := t.TempDir()
				p1 := filepath.Join(d1, "mytool")
				p2 := filepath.Join(d2, "mytool")
				if err := os.WriteFile(
					p1,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermReadWrite,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}
				if err := os.WriteFile(
					p2,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermExecutable,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}

				t.Setenv("PATH", strings.Join([]string{d1, d2}, ":"))
				return p2
			},
		},
		// This tests we return first executable when there are multiple
		// under different directories in PATH
		"WithMultilpeDirSameExecutable": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
				// GIVEN a directory d1 with filename f1 with non-exec permissions
				// GIVEN a directory d2 with filename f2 with exec permissions
				// GIVEN a directory d3 with filename f2 with exec permissions
				// GIVEN an env PATH d1:d2:d3
				d1 := t.TempDir()
				d2 := t.TempDir()
				d3 := t.TempDir()
				p1 := filepath.Join(d1, "badTool")
				p2 := filepath.Join(d2, "mytool")
				p3 := filepath.Join(d3, "mytool")
				if err := os.WriteFile(
					p1,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermReadWrite,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}
				if err := os.WriteFile(
					p2,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermExecutable,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}
				if err := os.WriteFile(
					p3,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermExecutable,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}

				t.Setenv("PATH", strings.Join([]string{d1, d2, d3}, ":"))
				return p2
			},
		},
		"WithNoExecutableFiles": {
			fileName:         "invalid_command",
			executableExists: false,
			setup: func() string {
				// GIVEN a directory d1 with filename f1 with non-exec permissions
				d1 := t.TempDir()
				p1 := filepath.Join(d1, "invalid_command")
				if err := os.WriteFile(
					p1,
					[]byte("#!/bin/sh\nexit 0\n"),
					PermReadWrite,
				); err != nil {
					t.Fatalf("create executable fixture: %v", err)
				}
				t.Setenv("PATH", d1)
				return p1
			},
		},
		"WithLeadingPathListSeparator": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
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

				// GIVEN leading path list separator
				t.Setenv("PATH", strings.Join([]string{"", dir}, ":"))
				return executable
			},
		},
		"WithTrailingPathListSeparator": {
			fileName:         "mytool",
			executableExists: true,
			setup: func() string {
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

				// GIVEN leading path list separator
				t.Setenv("PATH", strings.Join([]string{dir, ""}, ":"))
				return executable
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			executableDir := tc.setup()
			// WHEN calling findExecutable with executable name and PATH
			got, found := findExecutable(tc.fileName, "PATH")
			// THEN we return the found executable path (if exists)
			if tc.executableExists {
				if !found {
					t.Fatalf("executable %s should exist", tc.fileName)
				}

				if got != executableDir {
					t.Errorf("got %q, want %q", got, executableDir)
				}
			} else {
				if found {
					t.Fatalf("should not find executable: %s", tc.fileName)
				}
				if got != "" {
					t.Errorf("got path: %q", got)
				}
			}
		})
	}
}

func TestFindExecutableReturnsFileSeparatorWithRelativePath(t *testing.T) {
	// GIVEN an executable in the current directory,
	// with PATH containing a relative directory.
	dir := t.TempDir()
	t.Chdir(dir)
	want := "./mytool"

	const command = "mytool"

	if err := os.WriteFile(
		command,
		[]byte("#!/bin/sh\nexit 0\n"),
		PermExecutable,
	); err != nil {
		t.Fatalf("create executable fixture: %v", err)
	}

	t.Setenv("PATH", ".")

	// WHEN finding the executable.
	got, found := findExecutable(command, "PATH")

	// THEN it is found and with "/" file separator attached
	if !found {
		t.Fatalf("should find command: %s", command)
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
