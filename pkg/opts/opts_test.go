package opts

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGotopt2(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		input     string
		expected  string
		wantError error
	}{
		{
			name: "Basic",
			args: []string{"-foo=bar", "arg"},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: string
- name: "baz"
  help: "This is baz"
  type: string
  default: "value"
`,
			expected: `# gotopt2:generated:begin
gotopt2_baz="value"
gotopt2_foo="bar"
gotopt2_args__=("arg")
# gotopt2:generated:end
`,
		},
		{
			name: "Help",
			args: []string{"--help"},
			input: `
flags:
- name: "foo"
  type: string
`,
			wantError: flag.ErrHelp,
		},
		{
			name: "String list",
			args: []string{"-foo=bar,baz,bat", "arg"},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: stringlist
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo__list=("bar" "baz" "bat")
gotopt2_args__=("arg")
# gotopt2:generated:end
`,
		},
		{
			// Note, an empty list is indistinguishable from an undefined variable.
			name: "Empty string list",
			args: []string{"-foo="},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: stringlist
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo__list=()
# gotopt2:generated:end
`,
		},
		{
			name: "Basic with declaration",
			args: []string{"-foo=bar", "arg"},
			input: `
declaration: "some_declaration"
flags:
- name: "foo"
  help: "This is foo"
  type: string
`,
			expected: `# gotopt2:generated:begin
some_declaration gotopt2_foo="bar"
some_declaration gotopt2_args__=("arg")
# gotopt2:generated:end
`,
		},
		{
			name: "Basic with prefix",
			args: []string{"-foo=bar", "arg"},
			input: `
prefix: "some_prefix_"
flags:
- name: "foo"
  help: "This is foo"
  type: string
`,
			expected: `# gotopt2:generated:begin
some_prefix_gotopt2_foo="bar"
some_prefix_gotopt2_args__=("arg")
# gotopt2:generated:end
`,
		},
		{
			name: "Basic with ALL_CAPS",
			args: []string{"-foo=bar", "arg"},
			input: `
ALL_CAPS: true
flags:
- name: "foo"
  help: "This is foo"
  type: string
- name: "baz"
  help: "This is baz"
  type: string
  default: "value"
`,
			expected: `# gotopt2:generated:begin
GOTOPT2_BAZ="value"
GOTOPT2_FOO="bar"
GOTOPT2_ARGS__=("arg")
# gotopt2:generated:end
`,
		},
		{
			name: "Unknown flag value",
			args: []string{"-foo=bar"},
			input: `
flags:
- name: "baz"
  help: "This is baz"
  default: "default_value"
  type: string
`,
			wantError: fmt.Errorf("not defined: -foo"),
		},
		{
			name: "Bool True",
			args: []string{"--foo"},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo="true"
# gotopt2:generated:end
`,
		},
		{
			name: "Hyphen in flag name",
			args: []string{"--foo-bar"},
			input: `
flags:
- name: "foo-bar"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo_bar="true"
# gotopt2:generated:end
`,
		},
		{
			name: "Arg with hyphen",
			args: []string{"--", "--foo"},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo=""
gotopt2_args__=("--foo")
# gotopt2:generated:end
`,
		},
		{
			name: "Bool False",
			args: []string{},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo=""
# gotopt2:generated:end
`,
		},
		{
			name: "Bool False with custom value for false",
			args: []string{},
			input: `
falseValue: "false"
flags:
- name: "foo"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo="false"
# gotopt2:generated:end
`,
		},
		{
			name: "Bool with arg",
			args: []string{"--foo", "file1", "file2"},
			input: `
flags:
- name: "foo"
  help: "This is foo"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_foo="true"
gotopt2_args__=("file1" "file2")
# gotopt2:generated:end
`,
		},
		{
			name: "One of each",
			args: []string{
				"--strarg", "foo",
				"--strarg2=bar",
				"--intarg=10",
				"--boolarg=true",
				"param1",
				"param2",
			},
			input: `
flags:
- name: "strarg"
  type: string
- name: "strarg2"
  type: string
- name: "intarg"
  type: int
- name: "boolarg"
  type: bool
`,
			expected: `# gotopt2:generated:begin
gotopt2_boolarg="true"
gotopt2_intarg="10"
gotopt2_strarg2="bar"
gotopt2_strarg="foo"
gotopt2_args__=("param1" "param2")
# gotopt2:generated:end
`,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			err := Run(strings.NewReader(test.input), test.args, &b)
			if err != nil {
				if test.wantError == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), test.wantError.Error()) {
					t.Fatalf("Run(...)=%v, want %v", err, test.wantError)
				}
			}
			actuals := strings.Split(b.String(), "\n")
			expects := strings.Split(test.expected, "\n")
			if !cmp.Equal(expects, actuals) {
				t.Errorf("diff:\n%v\nactual:\n%+v\nwant:\n%+v",
					cmp.Diff(expects, actuals), actuals, expects)
			}
		})
	}
}
