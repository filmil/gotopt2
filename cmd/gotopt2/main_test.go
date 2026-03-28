package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestHelpExitCode(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Mock the stdin so it reads a config
		r, w, _ := os.Pipe()
		w.WriteString("flags:\n- name: foo\n  type: string\n")
		w.Close()
		os.Stdin = r

		os.Args = []string{"gotopt2", "--help"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHelpExitCode")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if e.ExitCode() != 11 {
			t.Fatalf("expected exit code 11, got %d", e.ExitCode())
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 11", err)
}
