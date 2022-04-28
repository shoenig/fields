package fields

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func try(t *testing.T, c Columns, in, expOut string, expErr bool) {
	inR := strings.NewReader(in)
	out := bytes.NewBuffer([]byte{})
	err := Apply(c, inR, out)

	if expErr {
		must.Error(t, err)
	} else {
		must.NoError(t, err)
		output := strings.TrimSpace(out.String())
		must.Eq(t, expOut, output)
	}
}

func Test_Apply_individual(t *testing.T) {
	ic2 := &individualColumn{column: 2}

	try(t, ic2, "a b c d e", "b", false)     // base one (col2 is 1th index)
	try(t, ic2, "a", "", true)               // out of bounds
	try(t, ic2, "a\tb\t\t c\td", "b", false) // tabs
	try(t, ic2, "a b", "b", false)           // last element

	ic0 := &individualColumn{column: 1}
	try(t, ic0, "a b c", "a", false) // first element
}

func Test_Apply_range(t *testing.T) {
	roc := &rangeOfColumns{
		columns: []int{2, 3, 4},
	}

	try(t, roc, "a b c d e f g", "c d e", false) // range
	try(t, roc, "a b c", "", true)               // out of bounds
}

func Test_Apply_leftExp(t *testing.T) {
	lec := &leftExpColumns{
		leftIndex: 3,
	}

	try(t, lec, "a b c d e f g", "a b c", false)
	try(t, lec, "a b", "", true) // out of bounds
}

func Test_Apply_rightExp(t *testing.T) {
	rec := &rightExpColumns{
		rightIndex: 3,
	}

	try(t, rec, "a b c d e f g", "c d e f g", false)
	try(t, rec, "a b", "", true) // out of bounds
}

const textBlock = `
a b c d
e f g

h i j
`

func Test_Apply_emptyLines(t *testing.T) {
	ic2 := &individualColumn{column: 2}
	try(t, ic2, textBlock, "b\nf\ni", false)
}
