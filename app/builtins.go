package main

import (
	"fmt"
	"io"
	"strings"
)

// handleEcho will write args directly to out
func handleEcho(args string, out io.Writer) error {
	_, err := fmt.Fprintln(out, strings.TrimSpace(args))
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
			_, err = fmt.Fprintf(out, "%s: not found\n", arg)
			if err != nil {
				return err
			}
		}
	}
	return err
}
