package main

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"

	"gophers.dev/cmds/fields"
)

const (
	// delimitersLabel = "cutset"
	newLineLabel = "newline"
)

const (
	exitOK  = 0
	exitErr = 1
)

func createFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("fields-flag-set", flag.ContinueOnError)
	// fs.String(delimitersLabel, "", "set of characters to cut on")
	fs.Bool(newLineLabel, true, "print trailing newline character")
	return fs
}

type inOut struct {
	stdIn  io.Reader
	stdOut io.Writer
	stdErr io.Writer
}

func newCommand(flags *flag.FlagSet, io inOut) *command {
	return &command{
		flags: flags,
		io:    io,
	}
}

type command struct {
	flags *flag.FlagSet
	io    inOut
}

func (c *command) Execute(args []string) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintf(c.io.stdErr, helpText)
		return exitErr
	}

	delimiters, columns, err := setup(c.flags, args)
	if err != nil {
		_, _ = fmt.Fprintf(c.io.stdErr, "fatal: %v\n", err)
		return exitErr
	}

	if err = do(delimiters, columns, c.io.stdIn, c.io.stdOut); err != nil {
		_, _ = fmt.Fprintf(c.io.stdErr, "fatal: %v\n", err)
		return exitErr
	}

	if err := finish(c.flags, c.io.stdOut); err != nil {
		_, _ = fmt.Fprintf(c.io.stdErr, "fatal: %v\n", err)
		return exitErr
	}

	return exitOK
}

func setup(fs *flag.FlagSet, args []string) (string, string, error) {
	if len(args) != 1 {
		return "", "", errors.Errorf("expected 1 argument, got %d", len(args))
	}

	err := fs.Parse(args)
	if err != nil {
		return "", "", errors.Wrapf(err, "parse args: %v", args)
	}

	// todo: support custom cutset
	//separators := fs.Lookup(delimitersLabel).Value.String()
	//if separators != "" {
	//	return "", "", errors.Errorf("custom cutsets not yet supported")
	//}
	separators := ""

	remArgs := fs.Args()
	return separators, remArgs[0], nil
}

// custom delimiters are not supported yet
func do(_, chosen string, input io.Reader, output io.Writer) error {
	cols, err := fields.Components(chosen)

	if err != nil {
		return errors.Wrap(err, "failed to parse column set")
	}

	combo := fields.Combine(cols)

	return fields.Apply(combo, input, output)
}

func finish(fs *flag.FlagSet, output io.Writer) error {
	trailNL, err := strconv.ParseBool(fs.Lookup(newLineLabel).Value.String())
	if err != nil {
		return err
	}

	if !trailNL {
		return errors.Errorf("appending newline not yet supported")
		// _, _ = io.WriteString(output, "\n")
		// we need to control appending a newline to every line, first
	}

	return nil
}
