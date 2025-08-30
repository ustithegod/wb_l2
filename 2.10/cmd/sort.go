package cmd

import (
	"bufio"
	"fmt"
	"slices"
	"sort"

	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "упрощённый аналог UNIX-утилиты sort.",
	Long: `Программа читает строки из STDIN и выводит их отсортированными.
	
Флаги:
-k [N] — сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию).
Например, «sort -k 2» отсортирует строки по второму столбцу каждой строки.
-n — сортировать по числовому значению (строки интерпретируются как числа).
-r — сортировать в обратном порядке (reverse).
-u — не выводить повторяющиеся строки (только уникальные).
-b — игнорировать хвостовые пробелы (trailing blanks).
-c — проверить, отсортированы ли данные; если нет, вывести сообщение об этом.`,

	Run: runSort,
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().IntP("column", "k", 0, "сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию)")
	sortCmd.Flags().BoolP("numeric", "n", false, "сортировать по числовому значению (строки интерпретируются как числа)")
	sortCmd.Flags().BoolP("reverse", "r", false, "сортировать в обратном порядке (reverse)")
	sortCmd.Flags().BoolP("unique", "u", false, "не выводить повторяющиеся строки (только уникальные)")
	sortCmd.Flags().BoolP("blank", "b", false, "игнорировать хвостовые пробелы (trailing blanks)")
	sortCmd.Flags().BoolP("check", "c", false, "проверить, отсортированы ли данные; если нет, вывести сообщение об этом")

}

func runSort(cmd *cobra.Command, args []string) {
	column, _ := cmd.Flags().GetInt("column")
	isNumeric, _ := cmd.Flags().GetBool("numeric")
	isReversed, _ := cmd.Flags().GetBool("reverse")
	isUnique, _ := cmd.Flags().GetBool("unique")
	ignoreBlanks, _ := cmd.Flags().GetBool("blank")
	check, _ := cmd.Flags().GetBool("check")

	if column < 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "column must be bigger than zero")
		return
	}

	scanner := bufio.NewScanner(cmd.InOrStdin())
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var result []string

	if isUnique {
		uniqueMap := make(map[string]struct{})
		for _, v := range lines {
			uniqueMap[v] = struct{}{}
		}
		for k := range uniqueMap {
			result = append(result, k)
		}
	} else {
		result = lines
	}

	if check {
		checkRes := slices.IsSorted(result)
		if checkRes {
			fmt.Fprintln(cmd.OutOrStdout(), "данные отсортированы!")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "данные не отсортированы!")
		}
		return
	}

	sortSlice := CustomStringSlice{
		lines:        result,
		column:       column - 1,
		isNumeric:    isNumeric,
		ignoreBlanks: ignoreBlanks,
	}

	if isReversed {
		sort.Sort(sort.Reverse(sortSlice))
	} else {
		sort.Sort(sortSlice)
	}

	for _, v := range result {
		fmt.Fprintln(cmd.OutOrStdout(), v)
	}
}
