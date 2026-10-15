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
		fmt.Print("$ ")

		// Will block until user enters input and presses enter
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			fmt.Fprintf(out, "reading input: %s", err)
			return err
		}
		if err == io.EOF {
			// Reached end of file and need to pipe out rest
			fmt.Fprintf(out, "%v: command not found\n", line)
		}

		// We need to strip the end of line so output command is written to stdout on a single line
		line = strings.TrimSuffix(line, "\n")

		fmt.Fprintf(out, "%v: command not found\n", line)
	}
}

func main() {
	// We need to create the loop phase of this execution
	run(os.Stdin, os.Stdout)
}
