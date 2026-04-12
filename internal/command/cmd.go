package command

import (
	"fmt"
	"io"
	"os"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/fields/internal/fields"
	"github.com/shoenig/fields/version"
)

type inOut struct {
	stdIn  io.Reader
	stdOut io.Writer
	stdErr io.Writer
}

func defaultIO() inOut {
	return inOut{
		stdIn:  os.Stdin,
		stdOut: os.Stdout,
		stdErr: os.Stderr,
	}
}

func Invoke(args []string) babycli.Code {
	return runWithIO(args, defaultIO())
}

func runWithIO(args []string, io inOut) babycli.Code {
	return babycli.New(&babycli.Configuration{
		Arguments: args,
		Version:   version.Version,
		Top: &babycli.Component{
			Name: "fields",
			Description: `
fields is a tool for selecting columns of text. One or more
individual columns and/or ranges of columns can be selected
by index. Indexes can be positive (counting from the left)
or negative (counting from the right). Ranges can be bounded
or unbounded (e.g. select all input to the left or right of N).
Input is expected on STDIN. Typically input is provided to
fields through pipes or redirection.

Syntax:
- Column indexes are counted starting from 1
- A positive index means start counting from the left
- A negative index means start counting from the right
- Closed ranges are denoted by two indexes separated by ':'
- Right open ranges are denoted by an index and a ':'
- Left open ranges are denoted by a ':' and an index
- High:Low ranges can reverse the order of columns
- Combine selections using comma separated list

Examples:
  fields 3 <<< 'a b c d e f g'              => c
  fields -3 <<< 'a b c d e f g'             => e
  fields 4: <<< 'a b c d e f g'             => d e f g
  fields -2: <<< 'a b c d e f g'            => f g
  fields 2:5 <<< 'a b c d e f g'            => b c d e
  fields 5:2 <<< 'a b c d e f g'            => e d c b
  fields -2,3:5,6: <<< 'a b c d e f g'      => f c d e f g`,
			Flags: babycli.Flags{
				{
					Long:    "newline",
					Type:    babycli.BooleanFlag,
					Default: &babycli.Default{Value: true, Show: true},
					Help:    "print trailing newline character",
				}, //nolint: exhaustruct
			},
			Function: func(c *babycli.Component) babycli.Code {
				return run(c, io)
			},
		},
	}).Run()
}

func run(c *babycli.Component, io inOut) babycli.Code {
	cols := c.Arguments()

	if len(cols) == 0 {
		if isEmpty(io.stdIn) {
			_, _ = fmt.Fprintln(io.stdErr, "no input detected (-h for help)")
			return babycli.Failure
		}
		return babycli.Success
	}

	if len(cols) != 1 {
		return babycli.Failure
	}

	columns, err := fields.Components(cols[0])
	if err != nil {
		_, _ = fmt.Fprintf(io.stdErr, "fatal: %v\n", err)
		return babycli.Failure
	}

	combo := fields.Combine(columns)

	if err := fields.Apply(combo, io.stdIn, io.stdOut); err != nil {
		_, _ = fmt.Fprintf(io.stdErr, "fatal: %v\n", err)
		return babycli.Failure
	}

	return babycli.Success
}

func isEmpty(r io.Reader) bool {
	stat, ok := r.(*os.File)
	if !ok {
		return false
	}
	// stdin is empty if it's not a terminal and has no data
	if stat == os.Stdin {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}
