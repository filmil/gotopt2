// Package main is the program gotopt2, for parsing command line arguments.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/filmil/gotopt2/pkg/opts"
)

func main() {
	go func() {
		// This goroutine will not run if the program exits before the
		// pause expires.
		const pause = 10
		<-time.After(pause * time.Second)
		fmt.Fprintf(os.Stderr, "gotopt2: %v seconds passed, did you forget to pass config as stdin?\n", pause)
		os.Exit(145)
	}()
	if err := opts.Run(os.Stdin, os.Args[1:], os.Stdout); err != nil {
		os.Exit(142)
	}
}
