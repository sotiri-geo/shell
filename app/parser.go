package main

import (
	"errors"
)

var ErrUnclosedSingleQuote = errors.New("closing single quote not found")

func parseArgs(args string) ([]string, error) {
	output := []string{}
	// Identify each word within a single quote
	// when we identify a single quote, we essentially need to traverse until we find
	// a closing quote, whilst preserving the characters found.

	i := 0
	// Behaves like a while loop
	for i < len(args) {
		// Handle single quotes in a single batch
		if args[i] == '\'' {
			if end, found := getClosingIndex(i, '\'', args); found {
				if i+1 < end {
					// We can also concat if we discover no whitespace between two single quotes
					// e.g. 'hello''world' -> "helloworld"
					if i > 0 && args[i-1] == args[i] {
						output[len(output)-1] = output[len(output)-1] + args[i+1:end]
					} else {
						output = append(output, args[i+1:end])
					}
				} else {
					// zero length quotes triggers a concat with previous word found
					// this can sometimes go all the way till end of string if sep not found
					if len(output) > 0 {
						start := end + 1
						if start >= len(args) {
							break
						}
						if args[start] == '\'' {
							i = start
							continue
						}
						if end, found := getClosingIndex(start, ' ', args); found {
							output[len(output)-1] = output[len(output)-1] + args[start:end]
						} else {
							output[len(output)-1] = output[len(output)-1] + args[start:]
							break
						}
					}
				}
				// move to the char next to the closing quote
				i = end + 1
			} else {
				return nil, ErrUnclosedSingleQuote
			}
			continue
		}
		// We now need to handle non-quoted chacacters that are not whitespace
		if args[i] != ' ' {
			// Walk up to a whitespace or a single quote
			start := i
			for i < len(args) {
				// Break out of current buffer and flush
				if args[i] == ' ' || args[i] == '\'' {
					output = append(output, args[start:i])
					break
				}
				i++
			}

			// If exhausted all characters: flush to buffer
			if i == len(args) {
				output = append(output, args[start:i])
			}

			continue
		}
		i++
	}

	return output, nil
}

func getClosingIndex(start int, sep byte, args string) (end int, found bool) {
	// We first inspect the next element in array
	end = start + 1
	for end < len(args) {
		if args[end] == sep {
			return end, true
		}
		end++
	}
	return end, false
}
