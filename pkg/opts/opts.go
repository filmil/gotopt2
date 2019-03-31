// Package opts parses the options from the configuration.
package opts

import (
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config represents the option parsing configuration
type Config struct {
	// Flags describes each flag to be parsed and its configuration.
	Flags []Flag `yaml:"flags"`
	// FalseValue is the value to generate for a false-valued Boolean variable.
	// It is "" by default, but some scripts may find it more convenient for
	// that text to be "false", or "off" or some such.  Note that this does not
	// affect the way the user can specify the flag.
	FalseValue string `yaml:"falseValue"`
	// AllCaps causes the variable name to be rendered in ALL_CAPS.
	AllCaps bool `yaml:"ALL_CAPS"`
	// Prefix is prepended to the full variable name.
	Prefix string `yaml:"prefix"`
	// Declaration is the default declaration word to use.  For example
	// "readonly" or "local".
	Declaration string `yaml:"declaration"`
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
		return err
	}

	fmt.Fprintf(w, "# gotopt2:generated:begin\n")
	wrFlags(fs, c.FalseValue, c.AllCaps, c.Prefix, c.Declaration, w)
	wrArgs(args, fs, c.Prefix, c.Declaration, c.AllCaps, w)
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
			log.Printf("unknown value: %+v", f)
			// Skip unknown values?
		}
	}
	return fs, nil
}

func declLine(name, value, falseVal, prefix, decl string, toUpper, quote bool) string {
	r := strings.NewReplacer("-", "_")
	name = r.Replace(name)
	fullVarName := fmt.Sprintf("%vgotopt2_%v", prefix, name)
	if toUpper {
		fullVarName = strings.ToUpper(fullVarName)
	}
	assignment := fmt.Sprintf("%v=%q", fullVarName, value)
	if !quote {
		assignment = fmt.Sprintf("%v=%v", fullVarName, value)
	}
	if decl == "" {
		return assignment
	}
	return strings.Join([]string{decl, assignment}, " ")
}

func wrFlags(fs *flag.FlagSet, falseVal string, toUpper bool,
	prefix string, decl string, w io.Writer) {
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
		if v == "" {
			v = falseVal
		}
		dl := declLine(f.Name, v, falseVal, prefix, decl, toUpper, true)
		out = append(out, fmt.Sprintf("%s\n", dl))
	})
	// Ensure that the output is stable.
	sort.Strings(out)
	for _, s := range out {
		fmt.Fprintf(w, s)
	}
}

// Quote remaining args if nonempty
func wrArgs(args []string, fs *flag.FlagSet, prefix, decl string, toUpper bool, w io.Writer) {
	if fs.NArg() == 0 {
		return
	}
	var a []string
	for _, arg := range fs.Args() {
		a = append(a, fmt.Sprintf("%q", arg))
	}
	allArgs := strings.Join(a, " ")
	dl := declLine("args__", fmt.Sprintf("(%s)", allArgs), "()", prefix, decl, toUpper, false)
	fmt.Fprintf(w, "%s\n", dl)
}
