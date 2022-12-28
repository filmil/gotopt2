#!/usr/bin/env bats

# The GOTOPT2 variable should be set in the BUILD.bazel script.
GOTOPT2="${GOTOPT2:-./cmd/gotopt2/gotopt2_/gotopt2}"

@test "Is program available" {
  [ "${GOTOPT2}" != "" ]
}

@test "Basic string flag parsing" {
  echo PWD:     ${PWD}
  echo GOTOPT2: ${GOTOPT2}
  result=$("${GOTOPT2}" --foo=bar <<EOF
flags:
- name: "foo"
  type: string
  default: "nothing"
  help: "This is some flag value."
EOF
)
  expected=$'# gotopt2:generated:begin
gotopt2_foo=\"bar\"
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

@test "Stringlist printing" {
  run "${GOTOPT2}" --list=eenie,meenie,minie,moe <<EOF
flags:
- name: "list"
  type: stringlist
  default: ""
  help: ""
EOF
  echo "${status}"
  echo "${lines[0]}"
  echo "${lines[1]}"
  echo "${lines[2]}"
  [ "${status}" -eq 0 ]
  [ "${lines[0]}" == "# gotopt2:generated:begin" ]
  [ "${lines[1]}" == "gotopt2_list__list=(\"eenie\" \"meenie\" \"minie\" \"moe\")" ]
  [ "${lines[2]}" == "# gotopt2:generated:end" ]
  [ "${#lines[@]}" -eq 3 ]
}

@test "Help" {
  run "${GOTOPT2}" --help <<EOF
flags:
- name: "list"
  type: stringlist
  default: ""
  help: ""
EOF
  echo "${status}"
  echo "${lines[0]}"
  echo "${lines[1]}"
  [ "${status}" -eq 11 ] # The exit code when --help is specified is 11
  [ "${lines[0]}" == "Usage:" ]
  [ "${lines[1]}" == "  -list value" ]
}

