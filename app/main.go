package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// execute accepts the unparsed user input text and routes to the requested command.
// It returns to the caller a shouldExit boolean which is true if
// the exit command was requested and any errors when writing output to writer.
func execute(line string, out io.Writer) (shouldExit bool, err error) {
	// Separate command to args to validate downstream
	command, args, _ := strings.Cut(line, " ")
	command = strings.TrimSpace(command)

	switch command {
	case "exit":
		return true, nil
	case "echo":
		_, err = fmt.Fprintln(out, strings.TrimSpace(args))
	default:
		_, err = fmt.Fprintf(out, "%v: command not found\n", strings.TrimSpace(line))
	}
	return false, err
}

// run will read from in and write to out. We accept both the io.Reader and io.Writer interfaces
// for this as it means the dependecy injection allows swapping values making testing easier.
func run(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	for {
		// Routes the prompt through the out (writer)
		fmt.Fprint(out, "$ ")

		// Will block until user enters input and presses enter
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			fmt.Fprintf(out, "reading input: %s", err)
			return err
		}
		// Early return if we reach EOF and no more bytes to read
		if err == io.EOF && len(line) == 0 {
			return nil
		}

		shouldExit, execErr := execute(line, out)
		if execErr != nil {
			return execErr
		}

		if shouldExit || err == io.EOF {
			return nil
		}
	}
}

func main() {
	// We need to create the loop phase of this execution
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "running shell program: %v", err)
		// Exit with non-zero status code to indicate failure
		os.Exit(1)
	}
}
