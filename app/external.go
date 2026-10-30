package main

import (
	"fmt"
	"io"
	"os/exec"
)

func runCommand(name string, out io.Writer) (found bool, err error) {
	if execPath, found := findExecutable(name, "PATH"); found {
		cmd := exec.Command(execPath)
		cmd.Stdout = out

		if err = cmd.Run(); err != nil {
			return false, fmt.Errorf("runnning command %v: %w", name, err)
		}
		return true, nil
	}
	return false, nil
}
