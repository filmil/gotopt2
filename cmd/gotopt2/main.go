// Package main is the program gotopt2, for parsing command line arguments.
package main

import (
	"os"

	"github.com/filmil/gotopt2/pkg/opts"
	"github.com/golang/glog"
)

func main() {
	if err := opts.Run(os.Stdin, os.Args, os.Stdout); err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}
}
