package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPassesParsedArgumentsToEcho(t *testing.T) {
	t.Run("JoinsFragmentsAroundEmptyQuotes", func(t *testing.T) {
		assertRunOutput(t, "echo hello''world\n", "$ helloworld\n$ ")
	})

	t.Run("PreservesSpacesInsideSingleQuotes", func(t *testing.T) {
		assertRunOutput(t, "echo 'hello   world'\n", "$ hello   world\n$ ")
	})
}

func TestExecutePreservesExternalArgumentBoundaries(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "print_args", "#!/bin/sh\nprintf '<%s>\\n' \"$@\"\n", 0o755)
	t.Setenv("PATH", dir)

	var output bytes.Buffer
	shouldExit, err := execute("print_args 'x y' z", &output)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if shouldExit {
		t.Error("external command exited the shell")
	}
	if got, want := output.String(), "<x y>\n<z>\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestRunDispatchesCommands(t *testing.T) {
	t.Run("DispatchesEcho", func(t *testing.T) {
		assertRunOutput(t, "echo Hello World\n", "$ Hello World\n$ ")
	})

	t.Run("DispatchesType", func(t *testing.T) {
		assertRunOutput(t, "type echo\n", "$ echo is a shell builtin\n$ ")
	})

	t.Run("TypeReportsEachArgument", func(t *testing.T) {
		assertRunOutput(t, "type echo type", "$ echo is a shell builtin\ntype is a shell builtin\n")
	})

	t.Run("ReportsUnknownCommand", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		assertRunOutput(t, "invalid_command\n", "$ invalid_command: command not found\n$ ")
	})

	t.Run("DispatchesBuiltinPwd", func(t *testing.T) {
		// GIVEN a known working directory and an empty directory on PATH
		dir := resolvedTempDir(t)
		t.Chdir(dir)
		t.Setenv("PATH", t.TempDir())

		// WHEN we execute pwd
		// THEN the builtin prints the working directory
		assertRunOutput(t, "pwd\n", "$ "+dir+"\n$ ")
	})

	t.Run("PrefersBuiltinPwdOverExternalCommand", func(t *testing.T) {
		// GIVEN a known working directory and a competing pwd executable on PATH
		dir := resolvedTempDir(t)
		t.Chdir(dir)
		binDir := t.TempDir()
		writeFixture(t, binDir, "pwd", "#!/bin/sh\necho EXTERNAL\n", 0o755)
		t.Setenv("PATH", binDir)

		// WHEN we execute pwd
		// THEN the builtin prints the working directory
		assertRunOutput(t, "pwd\n", "$ "+dir+"\n$ ")
	})

	t.Run("CdPrintsNoOutputOnSuccess", func(t *testing.T) {
		// GIVEN an existing absolute path and a different working directory
		targetDir := resolvedTempDir(t)
		t.Chdir(t.TempDir())

		// WHEN we execute cd with the absolute path
		// THEN the shell prints only the prompts
		assertRunOutput(t, "cd '"+targetDir+"'\n", "$ $ ")
	})

	t.Run("PwdPrintsDirectoryAfterCd", func(t *testing.T) {
		// GIVEN an existing absolute path and a different working directory
		targetDir := resolvedTempDir(t)
		t.Chdir(t.TempDir())

		// WHEN we execute cd followed by pwd
		// THEN pwd prints the new working directory
		assertRunOutput(t, "cd '"+targetDir+"'\n pwd\n", "$ $ "+targetDir+"\n$ ")
	})

	t.Run("CdReportsMissingDirectory", func(t *testing.T) {
		// GIVEN a known working directory and a path that does not exist
		startDir := resolvedTempDir(t)
		t.Chdir(startDir)
		missingDir := filepath.Join(startDir, "missing")

		// WHEN we execute cd with the missing path
		// THEN the shell reports that the directory does not exist
		assertRunOutput(t, "cd '"+missingDir+"'\n", "$ cd: "+missingDir+": No such file or directory\n$ ")
	})

	t.Run("RunsExternalCommand", func(t *testing.T) {
		// GIVEN an external executable on PATH that prints a message
		dir := t.TempDir()
		writeFixture(t, dir, "print_message", "#!/bin/sh\necho external\n", 0o755)
		t.Setenv("PATH", dir)

		// WHEN we run the external command
		// THEN its output appears after the prompt
		assertRunOutput(t, "print_message hello", "$ external\n")
	})
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "StopsAfterExit",
			input: "exit\n",
			want:  "$ ",
		},
		{
			name:  "PrintsPromptBeforeEOF",
			input: "",
			want:  "$ ",
		},
		{
			name:  "PromptsAgainAfterBlankLine",
			input: "\n",
			want:  "$ $ ",
		},
		{
			name:  "ExecutesFinalLineWithoutNewline",
			input: "echo Hello World",
			want:  "$ Hello World\n",
		},
		{
			name:  "ExecutesMultipleCommandsInOrder",
			input: "echo Hello World\n type echo\n",
			want:  "$ Hello World\n$ echo is a shell builtin\n$ ",
		},
		{
			name:  "ContinuesAfterUnclosedSingleQuote",
			input: "echo 'unclosed\necho still running\n",
			want:  "$ error: closing single quote not found\n$ still running\n$ ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertRunOutput(t, tc.input, tc.want)
		})
	}
}

func assertRunOutput(t *testing.T, input, want string) {
	t.Helper()

	var output bytes.Buffer
	if err := run(strings.NewReader(input), &output); err != nil {
		t.Fatalf("run(%q): %v", input, err)
	}
	if got := output.String(); got != want {
		t.Errorf("run(%q) output = %q, want %q", input, got, want)
	}
}
