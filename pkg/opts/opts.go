// Package opts parses the options from the configuration.
package opts

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/golang/glog"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the option parsing configuration
type Config struct {
	Flags []Flag `yaml:"flags"`
}

// FType is the type of the flag variable
type FType int

const (
	FTUnknown FType = iota
	FTString
	FTInt
	FTBool
)

// UnmarshalYAML implements yaml.Unmarshaler
func (f *FType) UnmarshalYAML(fn func(interface{}) error) error {
	var s string
	if err := fn(&s); err != nil {
		return fmt.Errorf("not a string: %v", s)
	}
	switch s {
	case "string":
		*f = FTString
	case "bool":
		*f = FTBool
	case "int":
		*f = FTInt
	default:
		*f = FTUnknown
	}
	return nil
}

// Flag represents a definition of a single flag
type Flag struct {
	Name string `yaml:"name"`
	Type FType  `yaml:"type"`
	Help string `yaml:"help"`
	// Default value as read from the configuration.  It needs to be parsed
	// into appropriate type before proceeding.
	RawDefault string `yaml:"default"`
}

// Run parses args based on a configuration supplied in r.  w gets all the
// output.
func Run(r io.Reader, args []string, w io.Writer) error {
	c, err := config(r)
	if err != nil {
		return err
	}

	// Configure the flag set.
	fs, err := flagSet(c)
	if err != nil {
		return err
	}
	if err := fs.Parse(args); err != nil {
		fs.Usage()
		return err
	}

	fmt.Fprintf(w, "# gotopt2:generated:begin\n")
	wrFlags(fs, w)
	wrArgs(args, fs, w)
	fmt.Fprintf(w, "# gotopt2:generated:end\n")
	return nil
}

// parseBool parses a boolean out of a string.  "" and "false" are false,
// "true" is true.
func parseBool(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		fallthrough
	case "":
		return false, nil
	default: // File not found :)
		return false, fmt.Errorf("not a bool value: %q", s)
	}
}

// parseInt parses an int out of a string
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("flag: %q: %v", s, err)
	}
	return v, nil
}

func config(r io.Reader) (Config, error) {
	d := yaml.NewDecoder(r)
	d.SetStrict(true)
	var (
		c   Config
		err error
	)
	if err = d.Decode(&c); err != nil {
		return c, fmt.Errorf("while decoding configuration: %v", err)
	}
	return c, nil
}

func flagSet(c Config) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	for _, f := range c.Flags {
		switch f.Type {
		case FTString:
			fs.String(f.Name, f.RawDefault, f.Help)
		case FTBool:
			def, err := parseBool(f.RawDefault)
			if err != nil {
				return nil, fmt.Errorf("flag: %q: %v", f.Name, err)
			}
			fs.Bool(f.Name, def, f.Help)
		case FTInt:
			def, err := parseInt(f.RawDefault)
			if err != nil {
				return nil, err
			}
			fs.Int(f.Name, def, f.Help)
		default:
			glog.Warningf("unknown value: %+v", f)
			// Skip unknown values?
		}
	}
	return fs, nil
}

func wrFlags(fs *flag.FlagSet, w io.Writer) {
	// Produce the output
	var out []string
	fs.VisitAll(func(f *flag.Flag) {
		var v string
		type isbooler interface {
			IsBoolFlag() bool
		}
		if _, ok := f.Value.(isbooler); !ok || f.Value.String() != "false" {
			v = f.Value.String()
		}
		out = append(out, fmt.Sprintf("readonly gotopt2_%v=%q\n", f.Name, v))
	})
	// Ensure that the output is stable.
	sort.Strings(out)
	for _, s := range out {
		fmt.Fprintf(w, s)
	}
}

// Quote remaining args if nonempty
func wrArgs(args []string, fs *flag.FlagSet, w io.Writer) {
	if fs.NArg() == 0 {
		return
	}
	var a []string
	for _, arg := range fs.Args() {
		a = append(a, fmt.Sprintf("%q", arg))
	}
	// Add leftover args to array.
	fmt.Fprintf(w, "readonly gotopt2_args__=(%s)\n", strings.Join(a, " "))
}
