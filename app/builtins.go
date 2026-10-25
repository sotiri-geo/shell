package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
			// Fetch path
			dir := os.Getenv("PATH")

			executable := filepath.Join(dir, arg)
			// check existance and then call stats
			fi, err := os.Stat(executable)
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(out, "%s: not found\n", arg)
				continue
			}
			if err != nil {
				return fmt.Errorf("reading file stats: %w", err)
			}
			// Executable file
			if fi.Mode().Perm() == 0o755 {
				_, err = fmt.Fprintf(out, "%s is %s\n", arg, executable)
				return err

			}
		}
	}
	return err
}
