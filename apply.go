package fields

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Apply will actually select the desired columns from the input and write them
// back out into the output in the same order in which they were requested.
func Apply(columns Columns, input io.Reader, output io.Writer) error {
	// Need to know how many words there are so columns can know how many /
	// which words to select on to forward into the output.
	// There's a more efficient way to do this that involves reading words one
	// at a time and keeping track until we resolve the columns selection,
	// but meh.

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		if line := next(scanner); line != "" {
			if err := process(
				line,
				output,
				columns,
			); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func next(scanner *bufio.Scanner) string {
	return strings.TrimSpace(scanner.Text())
}

func process(line string, output io.Writer, columns Columns) error {
	// get the words of the input
	words, err := read(line)
	if err != nil {
		return err
	}

	// compute which columns are going to be selected
	cols := columns.Columns(len(words))

	// check the selected columns are possible
	if err := validate(cols, len(words)); err != nil {
		return err
	}

	// select the words of the wanted columns
	chosen := zip(cols, words)

	// write the selected columns to the output
	if err := write(chosen, output); err != nil {
		return err
	}

	// append a newline character after each line
	if _, err := io.WriteString(output, "\n"); err != nil {
		return err
	}

	return nil
}

func write(words []string, output io.Writer) error {
	o := strings.Join(words, " ")
	_, err := io.WriteString(output, o)
	return err
}

func validate(cols []int, n int) error {
	if len(cols) == 0 {
		return fmt.Errorf("no columns match input length")
	}

	max := n - 1
	for _, col := range cols {
		if col > max {
			return fmt.Errorf("column %d is out of range [0,%d]", col, max)
		}
	}
	return nil
}

func zip(cols []int, words []string) []string {
	chosen := make([]string, 0, len(cols))
	for _, col := range cols {
		chosen = append(chosen, words[col])
	}
	return chosen
}

func read(line string) ([]string, error) {
	words := make([]string, 0, 10)

	input := strings.NewReader(line)
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	return words, scanner.Err()
}
