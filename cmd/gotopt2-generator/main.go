package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/filmil/gotopt2/pkg/opts"
	yaml "gopkg.in/yaml.v3"
)

//go:embed parser.sh.tmpl
var tmplFS embed.FS

func main() {
	if err := run(os.Stdin, os.Args[1:], os.Stdout); err != nil {
		if err == flag.ErrHelp {
			os.Exit(11)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type TemplateData struct {
	Flags          []TemplateFlag
	ArgsVarName    string
	SortedOutputs  []string
	ArgsOutputDecl string
}

type TemplateFlag struct {
	Name          string
	Type          opts.FType
	ActualVarName string
	DefaultValue  string
	TrueValue     string
	Help          string
}

func run(r io.Reader, args []string, w io.Writer) error {
	d := yaml.NewDecoder(r)
	var c opts.Config
	if err := d.Decode(&c); err != nil {
		return fmt.Errorf("decoding configuration: %v", err)
	}

	return generateBash(c, w)
}

func generateBash(c opts.Config, w io.Writer) error {
	tmpl, err := template.ParseFS(tmplFS, "parser.sh.tmpl")
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	data := TemplateData{
		ArgsVarName: varName("args__", c.Prefix, c.AllCaps),
	}

	var outputs []string

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

		data.Flags = append(data.Flags, TemplateFlag{
			Name:          f.Name,
			Type:          f.Type,
			ActualVarName: actualVarName,
			DefaultValue:  def,
			TrueValue:     trueVal,
			Help:          f.Help,
		})
		varNameStr := f.Name
		if f.Type == opts.FTStringList {
			name := varNameStr + "__list"
			outputs = append(outputs, fmt.Sprintf("  local %s_out\n  if [ ${#%s[@]} -eq 0 ]; then\n    %s_out=\"()\"\n  else\n    local %s_vals=()\n    for v in \"${%s[@]}\"; do\n      %s_vals+=(\"\\\"$v\\\"\")\n    done\n    %s_out=\"(${%s_vals[*]:-})\"\n  fi\n  echo \"%s\"",
				actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", declLineBash(name, fmt.Sprintf("${%s_out}", actualVarName+"__list"), c.Prefix, c.Declaration, c.AllCaps, false)))
		} else {
			outputs = append(outputs, fmt.Sprintf("  echo \"%s\"", declLineBash(varNameStr, fmt.Sprintf("${%s}", actualVarName), c.Prefix, c.Declaration, c.AllCaps, true)))
		}
	}

	sort.Strings(outputs)
	data.SortedOutputs = outputs

	data.ArgsOutputDecl = declLineBash("args__", "${args_out}", c.Prefix, c.Declaration, c.AllCaps, false)

	return tmpl.Execute(w, data)
}

func declLineBash(name, value, prefix, decl string, toUpper, quote bool) string {
	r := strings.NewReplacer("-", "_")
	name = r.Replace(name)
	fullVarName := fmt.Sprintf("%sgotopt2_%s", prefix, name)
	if toUpper {
		fullVarName = strings.ToUpper(fullVarName)
	}
	assignment := fmt.Sprintf("%s='%s'", fullVarName, value)
	if !quote {
		assignment = fmt.Sprintf("%s=%s", fullVarName, value)
	}
	if decl == "" {
		return assignment
	}
	return strings.Join([]string{decl, assignment}, " ")
}

func varName(name, prefix string, allCaps bool) string {
	r := strings.NewReplacer("-", "_")
	name = r.Replace(name)
	fullVarName := fmt.Sprintf("%sgotopt2_%s", prefix, name)
	if allCaps {
		fullVarName = strings.ToUpper(fullVarName)
	}
	return fullVarName
}
