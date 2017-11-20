// Author Seth Hoenig (seth.a.hoenig@gmail.com)

package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_indexes(t *testing.T) {
	tests := []struct {
		columns []int
		exp     map[int][]int
	}{
		{
			columns: []int{1},
			exp: map[int][]int{
				1: {0},
			},
		},
		{
			columns: []int{3, 4, 5},
			exp: map[int][]int{
				3: {0},
				4: {1},
				5: {2},
			},
		},
		{
			columns: []int{9, 3, 3, 5, 0, 1, 3, 0},
			exp: map[int][]int{
				9: {0},
				3: {1, 2, 6},
				5: {3},
				0: {4, 7},
				1: {5},
			},
		},
	}

	for _, test := range tests {
		result := indexes(test.columns)
		require.Equal(t, test.exp, result)
	}
}

func Test_noBlanks(t *testing.T) {
	tests := []struct {
		fields []string
		exp    []string
	}{
		{
			fields: []string{"foo", "bar", "baz"},
			exp:    []string{"foo", "bar", "baz"},
		},
		{
			fields: []string{"foo", "", "baz"},
			exp:    []string{"foo", "baz"},
		},
		{
			fields: []string{"", "foo", "bar", ""},
			exp:    []string{"foo", "bar"},
		},
	}

	for _, test := range tests {
		result := noBlanks(test.fields)
		require.Equal(t, test.exp, result)
	}
}

func Test_apply(t *testing.T) {
	tests := []struct {
		columns []int
		input   io.Reader
		exp     string
	}{
		{
			columns: []int{1},
			input:   strings.NewReader("foo bar baz"),
			exp:     "bar",
		},
		{
			columns: []int{0, 2},
			input:   strings.NewReader("foo bar baz"),
			exp:     "foo baz",
		},
		{
			columns: []int{3, 2, 2, 6, 1, 0, 1},
			input:   strings.NewReader("taco food truck is best thing ever yay woot"),
			exp:     "is truck truck ever food taco food",
		},
		{
			// there is no column 3
			columns: []int{1, 2, 3},
			input:   strings.NewReader("foo bar baz"),
			exp:     "bar baz",
		},
		{
			columns: []int{3, 9, 0, 1},
			input:   strings.NewReader("taco nights good for you"),
			exp:     "for taco nights",
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		apply(test.columns, test.input, writer)
		writer.Flush()
		s := buf.String()
		require.Equal(t, test.exp, s)
	}
}
