package opts

import (
	"flag"
	"io"
	"testing"
)

func BenchmarkWrArgs(b *testing.B) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	args := []string{"arg1", "arg2", "arg3", "arg4", "arg5", "arg6", "arg7", "arg8", "arg9", "arg10"}
	fs.Parse(args)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrArgs(args, fs, "prefix", "local", true, io.Discard)
	}
}
