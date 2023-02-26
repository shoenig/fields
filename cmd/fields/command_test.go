package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_createFlagSet(t *testing.T) {
	fs := createFlagSet()
	must.NotNil(t, fs)
}

func Test_setup(t *testing.T) {
	fs := createFlagSet()
	args := []string{"1,3,5"}

	seps, rem, err := setup(fs, args)
	must.NoError(t, err)
	must.Eq(t, "1,3,5", rem)
	must.Eq(t, "", seps)
}

func Test_setup_undefined_arg(t *testing.T) {
	fs := createFlagSet()
	args := []string{"--foo"}
	_, _, err := setup(fs, args)
	must.EqError(t, err, "parse args: [--foo]: flag provided but not defined: -foo")
}

func Test_setup_cutset_nys(t *testing.T) {
	t.Skip("custom cutset not yet supported")

	fs := createFlagSet()
	args := []string{"--cutset", "a"}
	_, _, err := setup(fs, args)
	must.EqError(t, err, "custom cutsets not yet supported")
}

func Test_setup_wrong_nargs(t *testing.T) {
	fs := createFlagSet()
	args := []string{"1", "2", "3"}
	_, _, err := setup(fs, args)
	must.EqError(t, err, "expected 1 argument, got 3")
}

func Test_do(t *testing.T) {
	s := "2,4"
	input := strings.NewReader("a b c d e f")
	var output bytes.Buffer
	err := do("", s, input, &output)
	must.NoError(t, err)

	result := output.String()
	must.Eq(t, "b d\n", result)
}

func Test_do_bad_fields(t *testing.T) {
	s := "huh"
	input := strings.NewReader("a b c d e f")
	var output bytes.Buffer
	err := do("", s, input, &output)
	must.EqError(t, err, `failed to parse column set: not valid syntax "huh"`)
}

func Test_finish(t *testing.T) {
	fs := createFlagSet()
	var output bytes.Buffer
	err := finish(fs, &output)
	must.NoError(t, err)
}

func Test_Execute(t *testing.T) {
	fs := createFlagSet()

	stdIn := strings.NewReader("a b c d e f")
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd := newCommand(fs, inOut{
		stdIn:  stdIn,
		stdOut: &stdOut,
		stdErr: &stdErr,
	})

	exitCode := cmd.Execute([]string{"2:4,1"})
	must.Eq(t, exitOK, exitCode)

	outStr := stdOut.String()
	must.Eq(t, "b c d a\n", outStr)

	errStr := stdErr.String()
	must.Eq(t, "", errStr)
}
