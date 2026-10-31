package main

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
)

func runCommand(name string, out io.Writer, arg ...string) (found bool, err error) {
	if execPath, found := findExecutable(name, "PATH"); found {

		command := filepath.Base(execPath)
		cmd := exec.Command(command, arg...)
		cmd.Stdout = out

		if err = cmd.Run(); err != nil {
			return false, fmt.Errorf("runnning command %v: %w", name, err)
		}
		return true, nil
	}
	return false, nil
}
