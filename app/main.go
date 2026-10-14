package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// We need to create the loop phase of this execution
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		// Will block until user enters input and presses enter
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			fmt.Printf("reading input: %s", err)
			return
		}

		// We need to strip the end of line so output command is written to stdout on a single line
		line = strings.TrimSuffix(line, "\n")

		fmt.Printf("%v: command not found\n", line)
	}
}
