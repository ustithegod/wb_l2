package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func unpackString(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	var (
		isEscaped bool = false
		buffer    string
		sb        strings.Builder
	)

	for i := 0; i < len(s); i++ {
		if string(s[i]) != `\` && !unicode.IsDigit(rune(s[i])) || isEscaped {
			buffer += string(s[i])
			isEscaped = false
			continue
		} else if string(s[i]) == `\` {
			isEscaped = true
			continue
		}

		if buffer == "" {
			return "", fmt.Errorf("string '%s' contains pair of digits or does not contain non-digit letters. to avoid that use escape-sequence", s)
		}

		sb.WriteString(buffer)
		letter := string(buffer[len(buffer)-1])
		numberOfRepetitions, _ := strconv.Atoi(string(s[i]))
		for range numberOfRepetitions - 1 {
			sb.WriteString(letter)
		}
		buffer = ""
	}

	sb.WriteString(buffer)

	result := sb.String()
	return result, nil
}
