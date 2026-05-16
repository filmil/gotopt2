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

func run(r io.Reader, shell string, w io.Writer) error {
	d := yaml.NewDecoder(r)
	var c opts.Config
	if err := d.Decode(&c); err != nil {
		return fmt.Errorf("decoding configuration: %v", err)
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
		if shell == "fish" {
			if f.Type == opts.FTStringList {
				name := varNameStr + "__list"
				outputs = append(outputs, fmt.Sprintf("  set -l %s_out\n  if test (count $%s) -eq 0\n    set %s_out \"\"\n  else\n    set -l %s_vals\n    for v in $%s\n      set -l esc (string replace -a \"'\" \"\\\\'\" \"$v\")\n      set -a %s_vals \"'$esc'\"\n    end\n    set %s_out (string join \" \" $%s_vals)\n  end\n  echo \"%s\"",
					actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", declLineFish(name, fmt.Sprintf("$%s_out", actualVarName+"__list"), c.Prefix, c.AllCaps, false)))
			} else {
				outputs = append(outputs, fmt.Sprintf("  set -l %s_esc (string replace -a \"'\" \"\\\\'\" \"$%s\")\n  echo \"%s\"", actualVarName, actualVarName, declLineFish(varNameStr, fmt.Sprintf("$%s_esc", actualVarName), c.Prefix, c.AllCaps, true)))
			}
		} else {
			if f.Type == opts.FTStringList {
				name := varNameStr + "__list"
				outputs = append(outputs, fmt.Sprintf("  local %s_out\n  if [ ${#%s[@]} -eq 0 ]; then\n    %s_out=\"()\"\n  else\n    local %s_vals=()\n    for v in \"${%s[@]}\"; do\n      %s_vals+=(\"\\\"$v\\\"\")\n    done\n    %s_out=\"(${%s_vals[*]:-})\"\n  fi\n  echo \"%s\"",
					actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", actualVarName+"__list", declLineBash(name, fmt.Sprintf("${%s_out}", actualVarName+"__list"), c.Prefix, c.Declaration, c.AllCaps, false)))
			} else {
				outputs = append(outputs, fmt.Sprintf("  echo \"%s\"", declLineBash(varNameStr, fmt.Sprintf("${%s}", actualVarName), c.Prefix, c.Declaration, c.AllCaps, true)))
			}
		}
	}

	sort.Strings(outputs)
	data.SortedOutputs = outputs

	if shell == "fish" {
		data.ArgsOutputDecl = declLineFish("args__", "$args_out", c.Prefix, c.AllCaps, false)
	} else {
		data.ArgsOutputDecl = declLineBash("args__", "${args_out}", c.Prefix, c.Declaration, c.AllCaps, false)
	}

	return tmpl.Execute(w, data)
}

func declLineFish(name, value, prefix string, toUpper, quote bool) string {
	r := strings.NewReplacer("-", "_")
	name = r.Replace(name)
	fullVarName := fmt.Sprintf("%sgotopt2_%s", prefix, name)
	if toUpper {
		fullVarName = strings.ToUpper(fullVarName)
	}
	if quote {
		assignment := fmt.Sprintf("set -g %s '%s'", fullVarName, value)
		return assignment
	}
	assignment := fmt.Sprintf("set -g %s %s", fullVarName, value)
	return assignment
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
