package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEchoJoinsArguments(t *testing.T) {
	var output bytes.Buffer
	if err := handleEcho([]string{"hello", "world"}, &output); err != nil {
		t.Fatalf("handleEcho: %v", err)
	}
	if got, want := output.String(), "hello world\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestType(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	testCases := []struct {
		name string
		args []string
		want string
	}{
		{name: "IdentifiesEchoAsBuiltin", args: []string{"echo"}, want: "echo is a shell builtin\n"},
		{name: "IdentifiesExitAsBuiltin", args: []string{"exit"}, want: "exit is a shell builtin\n"},
		{name: "IdentifiesTypeAsBuiltin", args: []string{"type"}, want: "type is a shell builtin\n"},
		{name: "IdentifiesPwdAsBuiltin", args: []string{"pwd"}, want: "pwd is a shell builtin\n"},
		{name: "IdentifiesCdAsBuiltin", args: []string{"cd"}, want: "cd is a shell builtin\n"},
		{
			name: "ReportsMultipleBuiltinsInOrder",
			args: []string{"type", "echo"},
			want: "type is a shell builtin\necho is a shell builtin\n",
		},
		{
			name: "ReportsUnknownCommand",
			args: []string{"invalid_command"},
			want: "invalid_command: not found\n",
		},
		{
			name: "ReportsBuiltinThenUnknownCommand",
			args: []string{"echo", "invalid_command"},
			want: "echo is a shell builtin\ninvalid_command: not found\n",
		},
		{
			name: "ReportsUnknownCommandThenBuiltin",
			args: []string{"invalid_command", "echo"},
			want: "invalid_command: not found\necho is a shell builtin\n",
		},
		{
			name: "ReportsMultipleUnknownCommandsInOrder",
			args: []string{"invalid_command", "banana"},
			want: "invalid_command: not found\nbanana: not found\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := handleType(tc.args, &output); err != nil {
				t.Fatalf("handleType(%q): %v", tc.args, err)
			}
			if got := output.String(); got != tc.want {
				t.Errorf("output = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTypeReportsExecutablePath(t *testing.T) {
	// GIVEN a directory on PATH with an executable named mytool
	dir := t.TempDir()
	executable := writeFixture(t, dir, "mytool", "#!/bin/sh\nexit 0\n", 0o755)
	t.Setenv("PATH", dir)

	var output bytes.Buffer
	// WHEN calling handleType with mytool
	if err := handleType([]string{"mytool"}, &output); err != nil {
		t.Fatalf("handleType: %v", err)
	}
	// THEN type prints the path to mytool
	if got, want := output.String(), "mytool is "+executable+"\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestFindExecutable(t *testing.T) {
	const command = "mytool"
	const script = "#!/bin/sh\nexit 0\n"
	separator := string(os.PathListSeparator)

	t.Run("FindsExecutableInFirstDirectory", func(t *testing.T) {
		// GIVEN the first directory on PATH contains mytool
		dir := t.TempDir()
		want := writeFixture(t, dir, command, script, 0o755)
		t.Setenv("PATH", dir)

		assertExecutablePath(t, command, want)
	})

	t.Run("SearchesNextDirectoryWhenCommandIsMissing", func(t *testing.T) {
		// GIVEN mytool is only in the second directory on PATH
		first, second := t.TempDir(), t.TempDir()
		want := writeFixture(t, second, command, script, 0o755)
		t.Setenv("PATH", strings.Join([]string{first, second}, separator))

		assertExecutablePath(t, command, want)
	})

	t.Run("SkipsNonExecutableWithSameName", func(t *testing.T) {
		// GIVEN PATH lists a non-executable mytool before an executable mytool
		first, second := t.TempDir(), t.TempDir()
		writeFixture(t, first, command, script, 0o644)
		want := writeFixture(t, second, command, script, 0o755)
		t.Setenv("PATH", strings.Join([]string{first, second}, separator))

		assertExecutablePath(t, command, want)
	})

	t.Run("ReturnsFirstExecutableOnPath", func(t *testing.T) {
		// GIVEN both directories on PATH contain an executable named mytool
		first, second := t.TempDir(), t.TempDir()
		want := writeFixture(t, first, command, script, 0o755)
		writeFixture(t, second, command, script, 0o755)
		t.Setenv("PATH", strings.Join([]string{first, second}, separator))

		assertExecutablePath(t, command, want)
	})

	t.Run("ReportsCommandMissingWhenFileIsNotExecutable", func(t *testing.T) {
		// GIVEN the only matching file on PATH is non-executable
		dir := t.TempDir()
		writeFixture(t, dir, command, script, 0o644)
		t.Setenv("PATH", dir)

		assertExecutablePath(t, command, "")
	})

	t.Run("SearchesPastLeadingEmptyPathEntry", func(t *testing.T) {
		// GIVEN PATH starts with an empty entry before the directory containing mytool
		t.Chdir(t.TempDir())
		dir := t.TempDir()
		want := writeFixture(t, dir, command, script, 0o755)
		t.Setenv("PATH", separator+dir)

		assertExecutablePath(t, command, want)
	})

	t.Run("FindsExecutableWithTrailingEmptyPathEntry", func(t *testing.T) {
		// GIVEN PATH ends with an empty entry after the directory containing mytool
		t.Chdir(t.TempDir())
		dir := t.TempDir()
		want := writeFixture(t, dir, command, script, 0o755)
		t.Setenv("PATH", dir+separator)

		assertExecutablePath(t, command, want)
	})

	t.Run("ReturnsDotSlashForCurrentDirectory", func(t *testing.T) {
		// GIVEN an executable in the current directory and a relative PATH entry of "."
		dir := t.TempDir()
		t.Chdir(dir)
		writeFixture(t, dir, command, script, 0o755)
		t.Setenv("PATH", ".")

		assertExecutablePath(t, command, "./"+command)
	})
}

func assertExecutablePath(t *testing.T, command, wantPath string) {
	t.Helper()

	// WHEN calling findExecutable with the command name and PATH
	got, found := findExecutable(command, "PATH")

	// THEN it reports whether the command was found and returns the expected path
	wantFound := wantPath != ""
	if got != wantPath || found != wantFound {
		t.Errorf("findExecutable(%q, \"PATH\") = (%q, %t), want (%q, %t)", command, got, found, wantPath, wantFound)
	}
}

func TestPwdPrintsWorkingDirectory(t *testing.T) {
	// GIVEN a known working directory
	dir := resolvedTempDir(t)
	t.Chdir(dir)

	var output bytes.Buffer
	// WHEN calling handlePwd
	if err := handlePwd(&output); err != nil {
		t.Fatalf("handlePwd: %v", err)
	}
	// THEN it writes the current working directory followed by a newline
	if got, want := output.String(), dir+"\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestCd(t *testing.T) {
	// GIVEN a starting directory, a missing path, and a separate directory set as HOME
	startDir := resolvedTempDir(t)
	targetDir := resolvedTempDir(t)
	missingDir := filepath.Join(startDir, "missing")
	t.Setenv("HOME", targetDir)

	testCases := []struct {
		name       string
		arg        string
		wantDir    string
		wantOutput string
	}{
		{name: "ChangesToAbsoluteDirectory", arg: targetDir, wantDir: targetDir},
		{name: "ChangesToParentDirectory", arg: "../", wantDir: filepath.Dir(startDir)},
		{name: "StaysInCurrentDirectoryForDot", arg: "./", wantDir: startDir},
		{name: "ExpandsTildeToHome", arg: "~", wantDir: targetDir},
		{
			name:       "ReportsMissingDirectoryAndStaysInPlace",
			arg:        missingDir,
			wantDir:    startDir,
			wantOutput: "cd: " + missingDir + ": No such file or directory\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(startDir)

			var output bytes.Buffer
			// WHEN changing directory using the case's argument
			if err := handleCd(tc.arg, &output); err != nil {
				t.Fatalf("handleCd(%q): %v", tc.arg, err)
			}
			// THEN cd produces the expected output and leaves the shell in the expected directory
			if got := output.String(); got != tc.wantOutput {
				t.Errorf("output = %q, want %q", got, tc.wantOutput)
			}

			got, err := os.Getwd()
			if err != nil {
				t.Fatalf("get working directory: %v", err)
			}
			if got != tc.wantDir {
				t.Errorf("working directory = %q, want %q", got, tc.wantDir)
			}
		})
	}
}
