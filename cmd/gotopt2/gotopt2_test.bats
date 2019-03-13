#!/usr/bin/env bats

# This should not be hardcoded but at the moment there is no good way to pass
# arguments into bats test scripts.
readonly GOTOPT2="./cmd/gotopt2/linux_amd64_stripped/gotopt2"

@test "Basic string flag parsing" {
  result=$("${GOTOPT2}" --foo=bar <<EOF
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

@test "Usage printing" {
  run "${GOTOPT2}" --unknown <<EOF
flags:
- name: "foo"
  type: string
  default: "nothing"
  help: "This is some flag value."
EOF
  echo "${status}"
  echo "${lines[0]}"
  echo "${lines[1]}"
  echo "${lines[2]}"
  echo "${lines[3]}"
  [ "${status}" -eq 142 ] # The exit code is randomly set to 142
  [ "${lines[0]}" == "flag provided but not defined: -unknown" ]
  [ "${lines[1]}" == "Usage:" ]
  [ "${lines[2]}" == "  -foo string" ]
  [ "${#lines[@]}" -eq 4 ]
}
