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

func BenchmarkWrFlags(b *testing.B) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.String("str1", "val1", "help")
	fs.Bool("bool1", true, "help")
	fs.Int("int1", 42, "help")
	v := StringListFlag{}
	fs.Var(&v, "list1", "help")

	fs.Parse([]string{"-str1=hello", "-bool1=false", "-int1=100", "-list1=a,b,c"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrFlags(fs, "false", "true", true, "prefix", "local", io.Discard)
	}
}
