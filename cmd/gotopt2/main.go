// Package main is the program gotopt2, for parsing command line arguments.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/google/glog"
)

func Run(r io.Reader, args []string, w io.Writer) error {
	return fmt.Errorf("tbd")
}

func main() {
	if err := Run(os.Stdin, os.Args, os.Stdout); err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}
}
