package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func unpackString(s string) (result string) {
	for i := 0; i < len(s)-1; i++ {
		current := s[i]
		if !unicode.IsDigit(rune(current)) {
			next := s[i+1]
			if unicode.IsDigit(rune(next)) {
				number, _ := strconv.Atoi(string(next))
				for i := 0; i < number; i++ {
					result += string(current)
				}
			} else {
				result += string(current)
			}
		}
	}
	return
}

func main() {
	str := "a4bc2d5e"
	fmt.Println(unpackString(str))
}

// a4bc2d5e
