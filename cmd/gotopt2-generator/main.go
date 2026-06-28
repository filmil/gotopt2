package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/filmil/gotopt2/pkg/opts"
)

//go:embed parser.sh.tmpl parser.fish.tmpl
var tmplFS embed.FS

func main() {
	shell := flag.String("shell", "bash", "The shell to generate code for (bash, fish)")
	flag.Parse()

	if err := run(os.Stdin, *shell, os.Stdout); err != nil {
		if err == flag.ErrHelp {
			os.Exit(11)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type TemplateData struct {
	Usage       string
	UsageQuoted string
	Flags       []TemplateFlag
	ArgsVarName string
}

type TemplateFlag struct {
	Name               string
	Type               opts.FType
	ActualVarName      string
	DefaultValue       string
	DefaultValueQuoted string
	TrueValue          string
	TrueValueQuoted    string
	Help               string
	HelpLineQuoted     string
}

// fishQuote wraps a string in single quotes and safely escapes existing backslashes and single quotes for fish shell.
func fishQuote(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	return "'" + s + "'"
}

func run(r io.Reader, shell string, w io.Writer) error {
	c, err := opts.ParseConfig(r)
	if err != nil {
		return err
	}

	return generateShell(c, w, shell)
}

func generateShell(c opts.Config, w io.Writer, shell string) error {
	tmplName := "parser.sh.tmpl"
	if shell == "fish" {
		tmplName = "parser.fish.tmpl"
	}

	tmpl, err := template.ParseFS(tmplFS, tmplName)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	usageQuoted := ""
	if shell == "fish" {
		usageQuoted = fishQuote(c.Usage)
	} else {
		usageQuoted = opts.ShellQuote(c.Usage)
	}

	data := TemplateData{
		Usage:       c.Usage,
		UsageQuoted: usageQuoted,
		ArgsVarName: varName("args__", c.Prefix, c.AllCaps),
	}

	trueVal := c.TrueValue
	if trueVal == "" {
		trueVal = "true"
	}

	for _, f := range c.Flags {
		actualVarName := varName(f.Name, c.Prefix, c.AllCaps)
		def := f.RawDefault
		if f.Type == opts.FTBool {
			if def == "" || def == "false" {
				def = c.FalseValue
			} else if def == "true" {
				def = trueVal
			}
		}

		defQuoted := ""
		trueValQuoted := ""
		helpLineQuoted := ""

		helpLine := fmt.Sprintf("        %s (default: %q)", f.Help, def)

		if shell == "fish" {
			defQuoted = fishQuote(def)
			trueValQuoted = fishQuote(trueVal)
			helpLineQuoted = fishQuote(helpLine)
		} else {
			defQuoted = opts.ShellQuote(def)
			trueValQuoted = opts.ShellQuote(trueVal)
			helpLineQuoted = opts.ShellQuote(helpLine)
		}

		data.Flags = append(data.Flags, TemplateFlag{
			Name:               f.Name,
			Type:               f.Type,
			ActualVarName:      actualVarName,
			DefaultValue:       def,
			DefaultValueQuoted: defQuoted,
			TrueValue:          trueVal,
			TrueValueQuoted:    trueValQuoted,
			Help:               f.Help,
			HelpLineQuoted:     helpLineQuoted,
		})
	}

	return tmpl.Execute(w, data)
}

func varName(name, prefix string, allCaps bool) string {
	r := strings.NewReplacer("-", "_")
	name = r.Replace(name)
	name = opts.SanitizeBashIdentifier(name)
	prefix = opts.SanitizeBashIdentifier(prefix)
	fullVarName := fmt.Sprintf("%sgotopt2_%s", prefix, name)
	if allCaps {
		fullVarName = strings.ToUpper(fullVarName)
	}
	return fullVarName
}
