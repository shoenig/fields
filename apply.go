// Author Seth Hoenig (seth.a.hoenig@gmail.com)

package main

import (
	"bufio"
	"io"
)

func apply(columns []int, input io.Reader, output io.Writer) error {
	fields := make([]string, len(columns), len(columns))
	c2i := indexes(columns)

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)
	col := 0
	for scanner.Scan() {
		if indexes, exists := c2i[col]; exists {
			field := scanner.Text()
			for _, index := range indexes {
				fields[index] = field
			}
		}
		col++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return write(noBlanks(fields), output)
}

func indexes(columns []int) map[int][]int {
	col2indexes := make(map[int][]int)
	for index, column := range columns {
		col2indexes[column] = append(col2indexes[column], index)
	}
	return col2indexes
}

func noBlanks(fields []string) []string {
	clean := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			clean = append(clean, field)
		}
	}
	return clean
}

func write(fields []string, output io.Writer) error {
	for i, field := range fields {
		if _, err := io.WriteString(output, field); err != nil {
			return err
		}

		if i < len(fields)-1 {
			if _, err := io.WriteString(output, " "); err != nil {
				return err
			}
		}
	}
	return nil
}
