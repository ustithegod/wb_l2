package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var cutCmd = &cobra.Command{
	Use:   "cut",
	Short: "Утилита, которая разбивает каждую строку по заданному разделителюи и выводит определённые поля.",
	Long: `Утилита, которая считывает входные данные (STDIN) и разбивает каждую строку по заданному разделителю, 
после чего выводит определённые поля (колонки).
	 
Флаги:
-f "fields" — указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.
Например: «-f 1,3-5» — вывести 1-й и с 3-го по 5-й столбцы.
-d "delimiter" — использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').
-s – (separated) только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).`,
	Run: runCut,
}

func init() {
	rootCmd.AddCommand(cutCmd)

	cutCmd.Flags().StringP("fields", "f", "", "указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.")
	cutCmd.Flags().StringP("delimiter", "d", "\t", "использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').")
	cutCmd.Flags().BoolP("separated", "s", false, "только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются.")
}

func runCut(cmd *cobra.Command, args []string) {
	fieldsString, _ := cmd.Flags().GetString("fields")
	delimiter, _ := cmd.Flags().GetString("delimiter")
	isSeparated, _ := cmd.Flags().GetBool("separated")

	fields, err := parseFields(fieldsString)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, v := range lines {
		if isSeparated {
			if !strings.Contains(v, delimiter) {
				continue
			}
		}
		var result []string
		splittedStr := strings.Split(v, delimiter)
		if len(fields) == 0 {
			fmt.Println(v)
		}
		for _, fieldNumber := range fields {
			if fieldNumber <= len(splittedStr) && fieldNumber > 0 {
				result = append(result, splittedStr[fieldNumber-1])
			}
		}
		fmt.Println(strings.Join(result, delimiter))
	}

}

func parseFields(fString string) ([]int, error) {
	var fields []int

	if fString == "" {
		return []int{}, nil
	}

	fieldsSlice := strings.Split(fString, ",")
	if matched, _ := regexp.MatchString(`^[0-9,-]*$`, fString); !matched {
		return nil, fmt.Errorf("invalid field value for: '%s'", fString)
	}
	for _, v := range fieldsSlice {
		if strings.Contains(v, "-") {
			diapason := strings.Split(v, "-")
			first, err := strconv.Atoi(diapason[0])
			if err != nil {
				return nil, err
			}
			second, err := strconv.Atoi(diapason[1])
			if err != nil {
				return nil, err
			}

			for i := first; i <= second; i++ {
				if !slices.Contains(fields, i) {
					fields = append(fields, i)
				}
			}
		} else {
			field, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		}
	}

	slices.Sort(fields)
	return fields, nil
}
