package cmd

import "github.com/spf13/cobra"

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

	fields := parseFields(fieldsString)

}

func parseFields(fString string) []int {
	for _, ch := range fString {

	}
}
