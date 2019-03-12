package opts

import (
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
			args: []string{"-foo=bar"},
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
readonly gotopt2_baz="value"
readonly gotopt2_foo="bar"
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
readonly gotopt2_foo="true"
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
readonly gotopt2_foo=""
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
readonly gotopt2_foo="true"
readonly gotopt2_args__=("file1" "file2")
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
readonly gotopt2_boolarg="true"
readonly gotopt2_intarg="10"
readonly gotopt2_strarg2="bar"
readonly gotopt2_strarg="foo"
readonly gotopt2_args__=("param1" "param2")
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
