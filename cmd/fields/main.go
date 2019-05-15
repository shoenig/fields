package main

import (
	"os"
)

func osIO() inOut {
	return inOut{
		stdIn:  os.Stdin,
		stdOut: os.Stdout,
		stdErr: os.Stderr,
	}
}

func main() {
	fs := createFlagSet()
	cmd := newCommand(fs, osIO())
	rc := cmd.Execute(os.Args[1:])
	os.Exit(rc)
}
