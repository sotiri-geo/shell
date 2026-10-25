package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const PermExecutable = 0o755

// handleEcho will write args directly to out
func handleEcho(args string, out io.Writer) error {
	// Using Fprintf for consistency so its obvious when printing a new line
	_, err := fmt.Fprintf(out, "%s\n", strings.TrimSpace(args))
	return err
}

// handleType will separate args by whitespace
// and loop over each string, writing to out whether
// its an identifiable builtin command or not recognised
func handleType(args string, out io.Writer) error {
	builtinTypes := map[string]bool{
		"echo": true,
		"exit": true,
		"type": true,
	}
	var err error
	for arg := range strings.SplitSeq(strings.TrimSpace(args), " ") {
		if builtinTypes[arg] {
			_, err = fmt.Fprintf(out, "%s is a shell builtin\n", arg)
			if err != nil {
				return err
			}
		} else {
			// This needs expanding to search across dir's listed in PATH
			path, found := findExecutable(arg, "PATH")
			if found {
				_, err = fmt.Fprintf(out, "%s is %s\n", arg, path)
				return err
			}
			fmt.Fprintf(out, "%s: not found\n", arg)
		}
	}
	return err
}

// findExecutable returns the first dir in PATH of the requested executable
func findExecutable(name, pathEnv string) (path string, found bool) {
	dir := os.Getenv(pathEnv)

	execPath := filepath.Join(dir, name)
	// check existance and then call stats
	fi, err := os.Stat(execPath)

	// We only care about file existing being executable
	if err == nil && fi.Mode().Perm() == PermExecutable {
		return execPath, true
	}

	return "", false
}
