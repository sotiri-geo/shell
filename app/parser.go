package main

import (
	"strings"
)

func parseArgs(args string) string {
	output := make([]string, 0, len(args))

	hasSingleQuotes := len(args) > 0 && args[0] == '\'' && args[len(args)-1] == '\''

	if hasSingleQuotes {
		return args[1 : len(args)-1]
	}
	// Remove multiple whitespaces
	for part := range strings.SplitSeq(args, " ") {
		if part != "" {
			output = append(output, strings.ReplaceAll(part, "'", ""))
		}
	}

	return strings.Join(output, " ")
}

