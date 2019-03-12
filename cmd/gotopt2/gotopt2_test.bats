#!/usr/bin/env bats

# This should not be hardcoded but at the moment there is no good way to pass
# arguments into bats test scripts.
readonly GOTOPT2="./cmd/gotopt2/linux_amd64/stripped/gotopt2"

@test "trying gotopt" {
  result=$(./cmd/gotopt2/linux_amd64_stripped/gotopt2 --foo=bar <<EOF
flags:
- name: "foo"
  type: string
  default: "nothing"
  help: "This is some flag value."
EOF
)
  expected=$'# gotopt2:generated:begin
readonly gotopt2_foo=\"bar\"
# gotopt2:generated:end'
  [ "${result}" == "${expected}" ]
}
