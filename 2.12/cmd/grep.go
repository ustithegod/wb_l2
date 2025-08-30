package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Утилита, предназначенная для поиска строк, соответствующих заданному шаблону",
	Long: `Утилита, предназначенная для поиска строк, соответствующих заданному шаблону.
	
Флаги:
-A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).
-B N — вывести N строк до каждой найденной строки.
-C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).
-c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).
-i — игнорировать регистр.
-v — инвертировать фильтр: выводить строки, не содержащие шаблон.
-F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).
-n — выводить номер строки перед каждой найденной строкой.`,

	Run: runGrep,
}

func init() {
	rootCmd.AddCommand(grepCmd)

	grepCmd.Flags().IntP("after-context", "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё (контекст)")
	grepCmd.Flags().IntP("before-context", "B", 0, "вывести N строк до каждой найденной строки")
	grepCmd.Flags().IntP("context", "C", 0, "вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N)")
	grepCmd.Flags().BoolP("count", "c", false, "выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число)")
	grepCmd.Flags().BoolP("ignore-case", "i", false, "игнорировать регистр")
	grepCmd.Flags().BoolP("invert-match", "v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон")
	grepCmd.Flags().BoolP("fixed-string", "F", false, "воспринимать шаблон как фиксированную строку, а не регулярное выражение")
	grepCmd.Flags().BoolP("line-number", "n", false, "выводить номер строки перед каждой найденной строкой")
}

func runGrep(cmd *cobra.Command, args []string) {
	after, _ := cmd.Flags().GetInt("after-context")
	before, _ := cmd.Flags().GetInt("before-context")
	context, _ := cmd.Flags().GetInt("context")
	isCount, _ := cmd.Flags().GetBool("count")
	isIgnore, _ := cmd.Flags().GetBool("ignore-case")
	isInvert, _ := cmd.Flags().GetBool("invert-match")
	isFixed, _ := cmd.Flags().GetBool("fixed-string")
	isLineNumber, _ := cmd.Flags().GetBool("line-number")

	substr := args[0]

	var data []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if context < 0 || after < 0 || before < 0 {
		fmt.Println("context must not be less than zero")
		return
	}

	if context > 0 {
		if context > after {
			after = context
		}
		if context > before {
			before = context
		}
	}

	fmt.Printf("before: %d, after: %d\n", before, after)
	fmt.Println(data)

	var matchArr []string
	for k, str := range data {
		findStr := substr
		if isIgnore {
			str = strings.ToLower(str)
			findStr = strings.ToLower(findStr)
		}

		if matches(findStr, str, isFixed, isInvert) {
			matchArr = append(matchArr, appendContext(data, findStr, isFixed, isInvert, isLineNumber, before, k, after)...)
		}
	}

	if isCount {
		fmt.Println(len(matchArr))
	} else {
		for _, v := range matchArr {
			fmt.Println(v)
		}
	}
}

func appendContext(data []string, substr string, isFixed, isInvert, isLineNumber bool, before, idx, after int) []string {
	startIndex := idx - before
	endIndex := idx + after

	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex >= len(data) {
		endIndex = len(data) - 1
	}

	if idx > 0 && idx < len(data)-1 {
		for i := idx - 1; i >= startIndex; i-- {
			if matches(substr, data[i], isFixed, isInvert) {
				startIndex = i + 1
				break
			}
		}

		for i := idx + 1; i <= endIndex; i++ {
			if matches(substr, data[i], isFixed, isInvert) {
				endIndex = i - 1
				break
			}
		}
	}

	result := make([]string, 0, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		if isLineNumber {
			result = append(result, fmt.Sprintf("%d:%s", i+1, data[i])) // Номер строки начинается с 1
		} else {
			result = append(result, data[i])
		}
	}

	return result
}

func matches(substr, str string, isFixed, isInvert bool) bool {
	if isFixed {
		words := strings.Fields(str)
		if isInvert {
			if !slices.Contains(words, substr) {
				return true
			}
		} else {
			if slices.Contains(words, substr) {
				return true
			}
		}
	} else {
		if isInvert {
			if matched, _ := regexp.MatchString(substr, str); !matched {
				return true
			}
		} else {
			if matched, _ := regexp.MatchString(substr, str); matched {
				return true
			}
		}
	}
	return false
}
