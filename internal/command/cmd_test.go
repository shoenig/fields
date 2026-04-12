package command

import (
	"bytes"
	"strings"
	"testing"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/test/must"
)

func Test_Invoke_basic(t *testing.T) {
	io := inOut{
		stdIn:  strings.NewReader("a b c d e f"),
		stdOut: &bytes.Buffer{},
		stdErr: &bytes.Buffer{},
	}

	code := runWithIO([]string{"2:4,1"}, io)
	must.Eq(t, babycli.Success, code)
	must.Eq(t, "b c d a\n", io.stdOut.(*bytes.Buffer).String())
}

func Test_Invoke_invalid_columns(t *testing.T) {
	io := inOut{
		stdIn:  strings.NewReader("a b c d e f"),
		stdOut: &bytes.Buffer{},
		stdErr: &bytes.Buffer{},
	}

	code := runWithIO([]string{"huh"}, io)
	must.Eq(t, babycli.Failure, code)
}

func Test_Invoke_no_args(t *testing.T) {
	io := inOut{
		stdIn:  strings.NewReader(""),
		stdOut: &bytes.Buffer{},
		stdErr: &bytes.Buffer{},
	}

	code := runWithIO([]string{}, io)
	must.Eq(t, babycli.Success, code)
}
