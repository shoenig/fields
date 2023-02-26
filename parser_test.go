package fields

import (
	"regexp"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/shoenig/regexplus"
)

func try1(t *testing.T, s string, match bool, exp string, re *regexp.Regexp) {
	matches := re.MatchString(s)

	if match {
		// check the regex matches
		must.True(t, matches)

		// check the regex gets the right value
		result := regexplus.FindNamedSubmatches(
			re, s,
		)["n"]
		must.Eq(t, exp, result)

	} else {
		// check the regex does not match
		must.False(t, matches)
	}
}

func Test_RE_individualRe(t *testing.T) {
	try1(t, "1", true, "1", individualRe)
	try1(t, "321", true, "321", individualRe)
	try1(t, "-5", true, "-5", individualRe)
	try1(t, "-987", true, "-987", individualRe)

	try1(t, ":1", false, "", individualRe)
	try1(t, "1:", false, "", individualRe)
	try1(t, "+1", false, "", individualRe)
	try1(t, "1+", false, "", individualRe)
}

func Test_RE_rangeRe(t *testing.T) {
	tryLR := func(s string, match bool, expL, expR string) {
		matches := rangeRe.MatchString(s)
		if match {
			// check the regex matches
			must.True(t, matches)

			subs := regexplus.FindNamedSubmatches(rangeRe, s)
			l := subs["l"]
			r := subs["r"]
			must.Eq(t, expL, l)
			must.Eq(t, expR, r)
		} else {
			// check the regex does not match
			must.False(t, matches)
		}
	}

	tryLR("1:5", true, "1", "5")
	tryLR("-5:-10", true, "-5", "-10")
	tryLR("-4:2", true, "-4", "2")

	tryLR("1", false, "", "")
	tryLR(":1", false, "", "")
	tryLR("1:", false, "", "")
}

func Test_leftExpRe(t *testing.T) {
	try1(t, ":1", true, "1", leftExpRe) // i.e. 0..1
	try1(t, ":123", true, "123", leftExpRe)
	try1(t, ":-3", true, "-3", leftExpRe) // i.e. 0..[-3]
	try1(t, ":-123", true, "-123", leftExpRe)

	try1(t, "1", false, "", leftExpRe)
	try1(t, "1:", false, "", leftExpRe)
	try1(t, "+3", false, "", leftExpRe)
}

func Test_rightExpRe(t *testing.T) {
	try1(t, "1:", true, "1", rightExpRe) // i.e. 1..[end]
	try1(t, "123:", true, "123", rightExpRe)
	try1(t, "-3:", true, "-3", rightExpRe) // i.e. [-3]..[end]
	try1(t, "-123:", true, "-123", rightExpRe)

	try1(t, "1", false, "", rightExpRe)
	try1(t, ":1", false, "", rightExpRe)
	try1(t, "+3", false, "", rightExpRe)
}

func Test_Interpret_individual(t *testing.T) {
	spanner, err := Single("1")
	must.NoError(t, err)

	_, ok := spanner.(*individualParser)
	must.True(t, ok)
}

func Test_Interpret_range(t *testing.T) {
	spanner, err := Single("2:5")
	must.NoError(t, err)

	_, ok := spanner.(*rangeParser)
	must.True(t, ok)
}

func Test_Interpret_leftExp(t *testing.T) {
	spanner, err := Single(":2")
	must.NoError(t, err)

	_, ok := spanner.(*leftExpParser)
	must.True(t, ok)
}

func Test_Interpret_rightExp(t *testing.T) {
	spanner, err := Single("2:")
	must.NoError(t, err)

	_, ok := spanner.(*rightExpParser)
	must.True(t, ok)
}

func Test_Interpret_malformed(t *testing.T) {
	try := func(input string) {
		_, err := Single(input)
		must.Error(t, err)
	}

	try("")
	try(" ")
	try("-")
	try("_")
	try("abc")
	try("2--3")
	try("-2-")
	try(":")
	try("2::3")
	try("2:-")
	try(":2-") // e.g. -2: would be valid
}

func Test_individualColumn_Columns(t *testing.T) {
	try := func(c, n int, exp []int) {
		ic := &individualColumn{column: c}
		iList := ic.Columns(n)
		must.Eq(t, exp, iList)
	}

	try(1, 3, []int{0})
	try(3, 3, []int{2})
	try(-1, 3, []int{2})
	try(-3, 3, []int{0})
	try(1, 1, []int{0})
}

func Test_rangeOfColumns_Columns(t *testing.T) {
	try := func(cols []int, n int, exp []int) {
		ic := &rangeOfColumns{columns: cols}
		iList := ic.Columns(n)
		must.Eq(t, exp, iList)
	}

	// all it does is copy the list of columns
	try([]int{2, 3, 4}, 10, []int{2, 3, 4})
}

func Test_leftExpColumns_Columns(t *testing.T) {
	try := func(leftIdx, n int, exp []int) {
		lec := &leftExpColumns{leftIndex: leftIdx}
		iList := lec.Columns(n)
		must.Eq(t, exp, iList)
	}

	try(2, 4, []int{0, 1})
	try(4, 8, []int{0, 1, 2, 3})

	try(-2, 5, []int{0, 1, 2, 3})
	try(-3, 3, []int{0})
	try(-2, 3, []int{0, 1})
	try(-1, 3, []int{0, 1, 2})
}

func Test_rightExpColumns_Columns(t *testing.T) {
	try := func(rightIdx, n int, exp []int) {
		rec := &rightExpColumns{rightIndex: rightIdx}
		iList := rec.Columns(n)
		must.Eq(t, exp, iList)
	}

	try(1, 4, []int{0, 1, 2, 3})
	try(2, 4, []int{1, 2, 3})
	try(4, 4, []int{3})
	try(-1, 5, []int{4})
	try(-2, 4, []int{2, 3})
	try(-4, 4, []int{0, 1, 2, 3})
}

func Test_individualParser_Spans(t *testing.T) {
	try := func(input string, n int, exp []int) {
		ip := &individualParser{}
		cols, err := ip.Spans(input)
		must.NoError(t, err)
		iList := cols.Columns(n)
		must.Eq(t, exp, iList)
	}

	try("1", 3, []int{0})
	try("3", 3, []int{2})
	try("-1", 3, []int{2})
	try("-3", 3, []int{0})
}

func Test_fill(t *testing.T) {
	try := func(left, right int, exp []int) {
		nums := fill(left, right)
		must.Eq(t, exp, nums)
	}

	// normal
	try(1, 3, []int{0, 1, 2})
	try(3, 4, []int{2, 3})
	try(3, 3, []int{2})

	// reverse
	try(3, 1, []int{2, 1, 0})
	try(4, 3, []int{3, 2})

	// left to (end-2)
	// try(3, -2, []int{})
	// negative numbers are not possible yet
	// the fill function does not know the number
	// of columns, we need to push that computation into
	// the Column function
}

func Test_rangeParser_Spans(t *testing.T) {
	try := func(input string, n int, exp []int) {
		rp := &rangeParser{}
		cols, err := rp.Spans(input)
		must.NoError(t, err)
		iList := cols.Columns(n)
		must.Eq(t, exp, iList)
	}

	try("1:3", 3, []int{0, 1, 2})
	try("2:4", 6, []int{1, 2, 3})

	try("3:1", 3, []int{2, 1, 0})
	try("4:2", 6, []int{3, 2, 1})
	// negative numbers not yet supported
}

func Test_leftExpParser_Spans(t *testing.T) {
	try := func(input string, n int, exp []int) {
		lep := &leftExpParser{}
		cols, err := lep.Spans(input)
		must.NoError(t, err)
		iList := cols.Columns(n)
		must.Eq(t, exp, iList)
	}

	// normal
	try(":1", 5, []int{0})
	try(":2", 5, []int{0, 1})
	try(":3", 5, []int{0, 1, 2})
	try(":4", 5, []int{0, 1, 2, 3})
	try(":5", 5, []int{0, 1, 2, 3, 4})

	// negatives
	try(":-1", 5, []int{0, 1, 2, 3, 4})
	try(":-2", 5, []int{0, 1, 2, 3})
	try(":-3", 5, []int{0, 1, 2})
	try(":-4", 5, []int{0, 1})
	try(":-5", 5, []int{0})
}

func Test_rightExpParser_Spans(t *testing.T) {
	try := func(input string, n int, exp []int) {
		lep := &rightExpParser{}
		cols, err := lep.Spans(input)
		must.NoError(t, err)
		iList := cols.Columns(n)
		must.Eq(t, exp, iList)
	}

	// normal
	try("1:", 5, []int{0, 1, 2, 3, 4})
	try("2:", 5, []int{1, 2, 3, 4})
	try("3:", 5, []int{2, 3, 4})
	try("4:", 5, []int{3, 4})
	try("5:", 5, []int{4})

	// negatives
	try("-1:", 5, []int{4})
	try("-2:", 5, []int{3, 4})
	try("-3:", 5, []int{2, 3, 4})
	try("-4:", 5, []int{1, 2, 3, 4})
	try("-5:", 5, []int{0, 1, 2, 3, 4})
}

func Test_Components(t *testing.T) {
	s := "1,-3,4:7"
	cols, err := Components(s)
	must.NoError(t, err)
	must.LenSlice(t, 3, cols)
}

func Test_Combine(t *testing.T) {
	s := "1,-3,4:7"
	cols, err := Components(s)
	must.NoError(t, err)
	columns := Combine(cols)
	iList := columns.Columns(10)
	must.Eq(t, []int{0, 7, 3, 4, 5, 6}, iList)
}
