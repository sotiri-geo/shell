package main

import (
	"errors"
	"strings"
)

var (
	ErrUnclosedSingleQuote = errors.New("closing single quote not found")
	ErrUnclosedDoubleQuote = errors.New("closing double quote not found")
	ErrEndingInBackslash   = errors.New("character must follow an escaped backslash")
)

const (
	SingleQuote = '\''
	DoubleQuote = '"'
	WhiteSpace  = ' '
	Backslash   = '\\'
)

// parseArgs preserves text inside single quotes and separates arguments at unquoted spaces.
func parseArgs(args string) ([]string, error) {
	output := []string{}
	var arg strings.Builder
	end := 0
	for end < len(args) {
		switch args[end] {
		case Backslash:
			// We escape the next character and preseve its literal
			// whilst igonring `\`
			if end == len(args)-1 {
				return nil, ErrEndingInBackslash
			}
			end++
			arg.WriteByte(args[end])
			end++
		case SingleQuote:
			closing, found := getClosingIndex(end, SingleQuote, args)
			if !found {
				return nil, ErrUnclosedSingleQuote
			}

			arg.WriteString(args[end+1 : closing])
			end = closing + 1
		case DoubleQuote:
			end++
			isClosed := false
			i := end

			for i < len(args) {
				if args[i] == DoubleQuote {
					isClosed = true
					// write remaining to buffer
					arg.WriteString(args[end:i])
					i++
					end = i // moves to next - "<next>
					break
				}
				if args[i] == Backslash {
					// Write up to backslash
					arg.WriteString(args[end:i])

					if i == len(args)-1 {
						return nil, ErrUnclosedDoubleQuote
					}
					// Decide if we preserve backslash or escape
					if args[i+1] == Backslash || args[i+1] == DoubleQuote {
						i++
					}
					// write the escaped char to buffer
					arg.WriteByte(args[i])
					// move to the next character after escaped character
					i++
					end = i
					continue
				}
				i++
			}

			if !isClosed {
				return nil, ErrUnclosedDoubleQuote
			}
		case WhiteSpace:
			// if there is existing arg - flush to output
			if arg.Len() > 0 {
				output = append(output, arg.String())
				arg.Reset()
				continue
			}
			// Traverse until exhausted non-whitespace characters
			for end < len(args) {
				if args[end] != WhiteSpace {
					break
				}
				end++
			}
		default:
			// Continue to grow buffer
			arg.WriteByte(args[end])
			end++

		}
	}
	// Flush remaining arg
	if arg.Len() > 0 {
		output = append(output, arg.String())
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
