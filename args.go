// Author Seth Hoenig (seth.a.hoenig@gmail.com)

package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type args struct {
	columns         []int
	noNewline       bool
	outputDelimiter string
}

func parseArgs() (args, error) {
	noNewline := flag.Bool("no-newline", false, "do not print newline after output")

	flag.Parse()

	if len(flag.Args()) == 0 {
		return args{columns: []int{}}, nil
	}

	columns, err := parseSet(flag.Args()[0])
	return args{
		columns:   columns,
		noNewline: *noNewline,
	}, err
}

// eg "1,2,5-7,9"
func parseSet(columns string) ([]int, error) {
	columns = strings.TrimSpace(columns)
	if columns == "" {
		return []int{}, nil
	}

	splitComma := strings.Split(columns, ",")
	nums := []int{}
	for _, split := range splitComma {
		cols, err := parseColumns(split)
		if err != nil {
			return nil, err
		}
		nums = append(nums, cols...)
	}
	return nums, nil
}

// eg "2" or "5-7"
func parseColumns(column string) ([]int, error) {
	dashIdx := strings.Index(column, "-")
	if dashIdx == -1 {
		i, err := parseNumber(column)
		return []int{i}, err
	}

	splitDash := strings.SplitN(column, "-", 2)
	if len(splitDash) != 2 {
		return nil, errors.Errorf("not a valid split: %s", column)
	}
	left := splitDash[0]
	right := splitDash[1]

	leftN, err := parseNumber(left)
	if err != nil {
		return nil, err
	}

	rightN, err := parseNumber(right)
	if err != nil {
		return nil, err
	}

	return fill(leftN, rightN), nil
}

func parseNumber(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Errorf("not a number: %q", value)
	}
	return i, nil
}

func fill(left, right int) []int {
	if left < right {
		length := (right - left) + 1
		nums := make([]int, length, length)
		for i := 0; i < length; i++ {
			nums[i] = left + i
		}
		return nums
	}

	length := (left - right) + 1
	nums := make([]int, length, length)
	for i := length - 1; i >= 0; i-- {
		nums[i] = left - i
	}
	return nums
}
