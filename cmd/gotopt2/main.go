// Package main is the program gotopt2, for parsing command line arguments.
package main

import (
	"os"

	"github.com/filmil/gotopt2/pkg/opts"
)

func main() {
	if err := opts.Run(os.Stdin, os.Args[1:], os.Stdout); err != nil {
		os.Exit(142)
	}
}
