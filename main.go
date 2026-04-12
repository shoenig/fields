package main

import (
	"os"

	"cattlecloud.net/go/babycli"
	"github.com/shoenig/fields/internal/command"
)

func main() {
	args := babycli.Arguments()
	rc := command.Invoke(args)
	os.Exit(rc)
}
