package main

import (
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
			args: []string{"--foo=bar"},
			input: `
flags:
- name: "foo"
  shortname: "f"
  help: "This is foo"
  type: string
`,
			expected: `# getopt:generated:begin
readonly gotopt2_foo="bar"
# getopt:generated:end
`,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			if err := Run(strings.NewReader(test.input), test.args, &b); err != nil {
				if test.wantError == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), test.wantError.Error()) {
					t.Errorf("Run(...)=%v, want %v", err, test.wantError)
				}
				actuals := strings.Split(b.String(), "\n")
				expects := strings.Split(test.expected, "\n")
				if !cmp.Equal(expects, actuals) {
					t.Errorf("diff:\n%v\nactual:\n%+v\nwant:\n%+v",
						cmp.Diff(expects, actuals), actuals, expects)
				}
			}
		})
	}
}
