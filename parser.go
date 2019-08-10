package fields

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Components returns a Spanner which is a superset of one or more other Spanner
// which can used to describe the columns to be selected.
//
// Example: "1,-3:-4,5:"
func Components(s string) ([]Columns, error) {
	splitComma := strings.Split(s, ",")

	var cols []Columns

	for _, single := range splitComma {
		spanner, err := Single(single)
		if err != nil {
			return nil, err
		}

		col, err := spanner.Spans(single)
		if err != nil {
			return nil, err
		}

		cols = append(cols, col)
	}

	return cols, nil
}

type comboColumns struct {
	cols []Columns
}

func Combine(cols []Columns) Columns {
	return &comboColumns{cols: cols}
}

func (cc *comboColumns) Columns(i int) []int {
	var indexes []int
	for _, c := range cc.cols {
		selections := c.Columns(i)
		indexes = append(indexes, selections...)
	}
	return indexes
}

// Single returns a Spanner which is capable of interpreting s as one item
// that can be parsed as one or more columns.
//
// Examples: "1", "-3:-4", "5:"
func Single(s string) (Spanner, error) {
	switch {
	case individualRe.MatchString(s):
		return &individualParser{}, nil
	case rangeRe.MatchString(s):
		return &rangeParser{}, nil
	case leftExpRe.MatchString(s):
		return &leftExpParser{}, nil
	case rightExpRe.MatchString(s):
		return &rightExpParser{}, nil
	default:
		return nil, errors.Errorf("not valid syntax %q", s)
	}
}

// Columns is used to actually select on columns of an input.
type Columns interface {
	// Columns returns the set of selected columns given a length of input.
	Columns(int) []int
}

type individualColumn struct {
	column int
}

func (sc *individualColumn) Columns(n int) []int {
	if sc.column >= 0 {
		return []int{sc.column - 1}
	}
	return []int{n + sc.column}
}

type rangeOfColumns struct {
	columns []int
}

func (roc *rangeOfColumns) Columns(int) []int {
	dest := make([]int, len(roc.columns))
	copy(dest, roc.columns)
	return dest
}

type leftExpColumns struct {
	leftIndex int
}

func (lec *leftExpColumns) Columns(n int) []int {
	if lec.leftIndex >= 0 {
		dest := make([]int, lec.leftIndex)
		for i := 0; i < lec.leftIndex; i++ {
			dest[i] = i
		}
		return dest
	}

	size := n - (-lec.leftIndex) + 1
	dest := make([]int, size)
	for i := 0; i < size; i++ {
		dest[i] = i
	}
	return dest
}

type rightExpColumns struct {
	rightIndex int
}

func (rec *rightExpColumns) Columns(n int) []int {
	if rec.rightIndex >= 0 {
		cardinality := n - rec.rightIndex + 1
		dest := make([]int, cardinality)
		for i := 0; i < cardinality; i++ {
			dest[i] = (rec.rightIndex - 1) + i
		}
		return dest
	}

	cardinality := -rec.rightIndex
	dest := make([]int, cardinality)
	for i := 0; i < cardinality; i++ {
		dest[i] = n - (cardinality - i)
	}
	return dest
}

var (
	individualRe = regexp.MustCompile(`^(?P<n>-?[\d]+)$`)
	rangeRe      = regexp.MustCompile(`^(?P<l>-?[\d]+):(?P<r>-?[\d]+)$`)
	leftExpRe    = regexp.MustCompile(`^:(?P<n>-?[\d]+)$`)
	rightExpRe   = regexp.MustCompile(`^(?P<n>-?[\d]+):$`)
)

// Spanner needs to parse things in the formats:
// 1, -1 (single) (- means count back from right)
// 1:5 (range)
// 1: (right-expansion) (i.e. from this column to the right)
// :5 (left-expansion) (i.e. from this column to the left)
// and any combination of thereof.
type Spanner interface {
	Spans(string) (Columns, error)
}

type individualParser struct{}

func (ip *individualParser) Spans(s string) (Columns, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}

	return &individualColumn{column: i}, nil
}

type rangeParser struct{}

func (rp *rangeParser) Spans(s string) (Columns, error) {
	numbers := strings.SplitN(s, ":", 2)
	if len(numbers) != 2 {
		return nil, errors.Errorf("not a valid range %q", s)
	}

	left := numbers[0]
	right := numbers[1]

	leftN, err := parseNumber(left)
	if err != nil {
		return nil, err
	}

	rightN, err := parseNumber(right)
	if err != nil {
		return nil, err
	}

	return &rangeOfColumns{
		columns: fill(leftN, rightN),
	}, nil
}

func parseNumber(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Errorf("not a number: %q", value)
	}
	return i, nil
}

func fill(left, right int) []int {
	if left <= right {
		length := (right - left) + 1
		nums := make([]int, length, length)
		for i := 0; i < length; i++ {
			nums[i] = left + i - 1
		}
		return nums
	}

	length := (left - right) + 1
	nums := make([]int, length, length)
	for i := length - 1; i >= 0; i-- {
		nums[i] = left - i - 1
	}
	return nums
}

type leftExpParser struct{}

func (lep *leftExpParser) Spans(s string) (Columns, error) {
	n, err := parseNumber(s[1:]) // e.g. ":3", ":-3"
	if err != nil {
		return nil, err
	}

	return &leftExpColumns{leftIndex: n}, nil
}

type rightExpParser struct{}

func (rep *rightExpParser) Spans(s string) (Columns, error) {
	num := s[0 : len(s)-1] // e.g. "3:", "-3:"
	n, err := parseNumber(num)
	if err != nil {
		return nil, err
	}

	return &rightExpColumns{rightIndex: n}, nil
}
