package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

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

		// We have text to write to writer
		if err == io.EOF {
			// Some text was recieved before reaching delim so we can send this to writer
			if len(line) > 0 {
				line = strings.TrimSuffix(line, "\n")

				fmt.Fprintf(out, "%v: command not found\n", line)
			}
			return nil
		}

		// Exit command is valid
		line = strings.TrimSuffix(line, "\n")
		if line == "exit" {
			os.Exit(0)
		}
		// Echo command is valid
		if args, exists := strings.CutPrefix(line, "echo"); exists {
			fmt.Fprintln(out, strings.TrimSpace(args))
			// return nil
		} else {
			fmt.Fprintf(out, "%v: command not found\n", line)
		}
	}
}

func main() {
	// We need to create the loop phase of this execution
	run(os.Stdin, os.Stdout)
}
