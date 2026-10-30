package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	PermExecutable = 0o755
	PermReadWrite  = 0o666
)

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

// findExecutable returns the first executable path for `name`.
// For example if the first executable path of `mytool` in PATH is `file/to/bin` then
// execPath would return file/to/bin/mytool. For relative path ".", we join the file separator
// i.e. `./mytool`.
func findExecutable(name, pathEnv string) (execPath string, found bool) {
	path := os.Getenv(pathEnv)

	// Order matters, we return the first found executable path for name
	for _, dir := range filepath.SplitList(path) {
		execPath := filepath.Join(dir, name)
		// Explictly constructing path for relative path as
		// its not handled in filepath.Join
		if dir == "." {
			execPath = strings.Join([]string{dir, name}, "/")
		}
		fi, err := os.Stat(execPath)

		if err == nil && fi.Mode().Perm() == PermExecutable {
			return execPath, true
		}
	}

	return "", false
}
