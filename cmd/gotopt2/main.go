// Package main is the program gotopt2, for parsing command line arguments.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/filmil/gotopt2/pkg/opts"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading from stdin: %v\n", err)
		os.Exit(142)
	}
	// If stdin is a TTY, it means no input is being piped, and we should
	// exit with a helpful error message.
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "gotopt2: configuration must be passed via stdin")
		os.Exit(12)
	}

	if err := opts.Run(os.Stdin, os.Args[1:], os.Stdout); err != nil {
		if err == flag.ErrHelp {
			// flag.ErrHelp means that the flag parser has written out the
			// usage.
			os.Exit(11)
		}
		fmt.Fprintf(os.Stderr, "gotopt2: %v\n", err)
		os.Exit(142)
	}
}
