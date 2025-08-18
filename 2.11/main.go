package main

import (
	"fmt"
	"sort"
	"strings"
)

func sortString(s string) string {
	temp := strings.Split(s, "")
	sort.Strings(temp)
	return strings.Join(temp, "")
}

func compare(s1, s2 string) bool {
	return sortString(s1) == sortString(s2)
}

func findAnagrams(data []string) map[string][]string {
	anagrams := make(map[string][]string)

	for _, str := range data {
		isNew := true
		for key := range anagrams {
			if compare(str, key) {
				anagrams[key] = append(anagrams[key], str)
				isNew = false
			}
		}
		if isNew {
			anagrams[str] = []string{str}
		}
	}

	for key, val := range anagrams {
		if len(val) == 1 {
			delete(anagrams, key)
		}
	}

	return anagrams
}

func main() {
	testData := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(findAnagrams(testData))
}
