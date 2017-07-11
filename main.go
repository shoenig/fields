// Author Seth Hoenig (seth.a.hoenig@gmail.com)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args, err := parseArgs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse columns: %v\n", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	if err := apply(args.columns, os.Stdin, writer); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse fields: %v\n", err)
		os.Exit(1)
	}

	if !args.noNewline {
		if _, err := writer.WriteString("\n"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write newline: %v\n", err)
			os.Exit(1)
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to flush output: %v\n", err)
		os.Exit(1)
	}

	return
}
