// Author Seth Hoenig (seth.a.hoenig@gmail.com)

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseSet(t *testing.T) {
	tests := []struct {
		columns string
		exp     []int
		expErr  bool
	}{
		// unspecified (all)
		{columns: "", exp: []int{}, expErr: false},

		// specific column numbers
		{columns: "1", exp: []int{1}, expErr: false},
		{columns: "1,2", exp: []int{1, 2}, expErr: false},
		{columns: "2,5,7", exp: []int{2, 5, 7}, expErr: false},
		{columns: "7,2,1", exp: []int{7, 2, 1}, expErr: false},

		// ranges
		{columns: "1-4", exp: []int{1, 2, 3, 4}, expErr: false},
		{columns: "5-2", exp: []int{5, 4, 3, 2}, expErr: false},

		// mixed
		{columns: "1,3-5,2,9-7,1", exp: []int{1, 3, 4, 5, 2, 9, 8, 7, 1}, expErr: false},
	}

	for _, test := range tests {
		cols, err := parseSet(test.columns)
		require.Equal(t, test.expErr, err != nil)
		require.Equal(t, test.exp, cols)
	}
}
