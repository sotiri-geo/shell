package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// execute accepts the unparsed user input text and routes to the requested command.
// It returns to the caller a shouldExit boolean which is true if
// the exit command was requested and any errors when writing output to writer.
func execute(line string, out io.Writer) (shouldExit bool, err error) {
	command, args, _ := strings.Cut(line, " ")

	switch command {
	case "exit":
		return true, nil
	case "echo":
		err = handleEcho(args, out)
	case "type":
		err = handleType(args, out)
	// Boundary to introduce external commands
	default:
		if command == "mytool" {
			fmt.Fprint(out, "true")
			return false, err
		}
		// Check its an command found on PATH
		if execPath, found := findExecutable(command, "PATH"); found {
			cmd := exec.Command(execPath)
			// Must write to writer
			cmd.Stdout = out

			if err = cmd.Run(); err != nil {
				return false, fmt.Errorf("runnning command %v: %w", command, err)
			}
			return false, err
		}

		_, err = fmt.Fprintf(out, "%v: command not found\n", strings.TrimSpace(command))
	}
	return false, err
}

// run will read from in and write to out. We accept both the io.Reader and io.Writer interfaces
// for this as it means the dependecy injection allows swapping values making testing easier.
func run(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	for {
		// Routes the prompt through the out (writer)
		if _, err := fmt.Fprint(out, "$ "); err != nil {
			fmt.Fprintf(out, "prompt error: %v", err)
		}

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

		// This is reached when there is text to be printed regardless
		// of if we recieved an io.EOF
		line = strings.TrimSpace(line)
		// When nothing is parsed it should continue for another prompt
		if len(line) == 0 {
			continue
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
