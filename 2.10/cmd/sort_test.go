package cmd

import (
	"bytes"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// setupCommand создаёт и настраивает rootCmd и sortCmd для каждого теста
func setupCommand(t *testing.T, input string, args []string) (*cobra.Command, *bytes.Buffer) {
	rootCmd := &cobra.Command{Use: "mysort"}
	// Создаём новую sortCmd для полной изоляции
	newSortCmd := &cobra.Command{
		Use:   "sort",
		Short: sortCmd.Short,
		Long:  sortCmd.Long,
		Run:   sortCmd.Run,
	}
	// Регистрируем флаги заново
	newSortCmd.Flags().IntP("column", "k", 0, "сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию)")
	newSortCmd.Flags().BoolP("numeric", "n", false, "сортировать по числовому значению (строки интерпретируются как числа)")
	newSortCmd.Flags().BoolP("reverse", "r", false, "сортировать в обратном порядке (reverse)")
	newSortCmd.Flags().BoolP("unique", "u", false, "не выводить повторяющиеся строки (только уникальные)")
	newSortCmd.Flags().BoolP("blank", "b", false, "игнорировать хвостовые пробелы (trailing blanks)")
	newSortCmd.Flags().BoolP("check", "c", false, "проверить, отсортированы ли данные; если нет, вывести сообщение об этом")

	rootCmd.AddCommand(newSortCmd)

	bufOut := bytes.NewBufferString("")
	rootCmd.SetOut(bufOut)
	rootCmd.SetIn(bytes.NewBufferString(input))
	rootCmd.SetArgs(args)

	return rootCmd, bufOut
}

// stripDebug удаляет отладочный вывод, начинающийся с "DEBUG:"
func stripDebug(output string) string {
	if strings.HasPrefix(output, "DEBUG:") {
		lines := strings.SplitN(output, "\n", 2)
		if len(lines) > 1 {
			return lines[1]
		}
		return ""
	}
	return output
}

func TestDefaultSort(t *testing.T) {
	input := "banana\napple\ncherry\n"
	args := []string{"sort"}
	expected := "apple\nbanana\ncherry\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSortBySecondColumn(t *testing.T) {
	input := "apple\t10\tred\nbanana\t5\tyellow\ncherry\t15\tred\n"
	args := []string{"sort", "-k", "2"}
	// Лексикографически: "10" < "15" < "5"
	expected := "apple\t10\tred\ncherry\t15\tred\nbanana\t5\tyellow\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestNumericSortBySecondColumn(t *testing.T) {
	input := "apple\t10\tred\nbanana\t5\tyellow\ncherry\t15\tred\n"
	args := []string{"sort", "-k", "2", "-n"}
	// Числовой порядок: 5 < 10 < 15
	expected := "banana\t5\tyellow\napple\t10\tred\ncherry\t15\tred\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestReverseSort(t *testing.T) {
	input := "banana\napple\ncherry\n"
	args := []string{"sort", "-r"}
	expected := "cherry\nbanana\napple\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestUniqueSort(t *testing.T) {
	input := "apple\nbanana\napple\ncherry\n"
	args := []string{"sort", "-u"}
	expected := "apple\nbanana\ncherry\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestCheckSortedData(t *testing.T) {
	input := "apple\nbanana\ncherry\n"
	args := []string{"sort", "-c"}
	expected := "данные отсортированы!\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestCheckUnsortedData(t *testing.T) {
	input := "banana\napple\ncherry\n"
	args := []string{"sort", "-c"}
	expected := "данные не отсортированы!\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSortByColumnWithUniqueAndReverse(t *testing.T) {
	input := "apple\t10\tred\nbanana\t5\tyellow\napple\t10\tred\n"
	args := []string{"sort", "-k", "2", "-u", "-r"}
	// Лексикографически в обратном порядке: "5" > "10"
	expected := "banana\t5\tyellow\napple\t10\tred\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestEmptyInput(t *testing.T) {
	input := ""
	args := []string{"sort"}
	expected := ""

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestInvalidColumn(t *testing.T) {
	input := "apple\t10\tred\n"
	args := []string{"sort", "-k", "-1"}
	expected := "column must be bigger than zero\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSortByColumnWithInsufficientColumns(t *testing.T) {
	input := "apple\t10\nbanana\ncherry\t15\tred\n"
	args := []string{"sort", "-k", "3"}
	// Строки с < 3 столбцов считаются пустыми
	expected := "apple\t10\nbanana\ncherry\t15\tred\n"

	rootCmd, bufOut := setupCommand(t, input, args)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got := stripDebug(bufOut.String())
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestCustomStringSliceSortBySecondColumn(t *testing.T) {
	lines := []string{
		"apple\t10\tred",
		"banana\t5\tyellow",
		"cherry\t15\tred",
	}
	result := make([]string, len(lines))
	copy(result, lines)
	sortSlice := CustomStringSlice{lines: result, column: 1, isNumeric: false, ignoreBlanks: false}
	sort.Sort(sortSlice)

	expected := []string{
		"apple\t10\tred",
		"cherry\t15\tred",
		"banana\t5\tyellow",
	}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCustomStringSliceNumericSortBySecondColumn(t *testing.T) {
	lines := []string{
		"apple\t10\tred",
		"banana\t5\tyellow",
		"cherry\t15\tred",
	}
	result := make([]string, len(lines))
	copy(result, lines)
	sortSlice := CustomStringSlice{lines: result, column: 1, isNumeric: true, ignoreBlanks: false}
	sort.Sort(sortSlice)

	expected := []string{
		"banana\t5\tyellow",
		"apple\t10\tred",
		"cherry\t15\tred",
	}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCustomStringSliceIgnoreBlanks(t *testing.T) {
	lines := []string{
		"apple\t10\tred  ",
		"banana\t5\tyellow\t",
		"cherry\t15\tred",
	}
	result := make([]string, len(lines))
	copy(result, lines)
	sortSlice := CustomStringSlice{lines: result, column: 1, isNumeric: true, ignoreBlanks: true}
	sort.Sort(sortSlice)

	expected := []string{
		"banana\t5\tyellow\t",
		"apple\t10\tred  ",
		"cherry\t15\tred",
	}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCustomStringSliceSortWithInsufficientColumns(t *testing.T) {
	lines := []string{
		"apple\t10",
		"banana",
		"cherry\t15\tred",
	}
	result := make([]string, len(lines))
	copy(result, lines)
	sortSlice := CustomStringSlice{lines: result, column: 2, isNumeric: false, ignoreBlanks: false}
	sort.Sort(sortSlice)

	expected := []string{
		"apple\t10",
		"banana",
		"cherry\t15\tred",
	}
	if !slices.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
