#! /bin/bash
# Args:
#   $1: the name of the gotopt2 binary to execute: this is a bazel quirk.
#   <rest>: the flag arguments to parse, see the BUILD rule for the passed args.

# This test script requires that the name of the "gotopt2" binary be the first
# arg in use.
readonly GOTOPT2="${1}"
shift

readonly output=$("${GOTOPT2}" "${@}" <<EOF
flags:
- name: "foo"
  type: string
  default: "something"
- name: "bar"
  type: int
  default: 42
- name: "baz"
  type: bool
  default: true
EOF
)

# Evaluate the output of the call to gotopt2, shell vars assignment is here.
eval "${output}"

# Quick check of the result.
if [[ "${gotopt2_foo}" != "bar" ]]; then
  echo "Want: bar; got: '${gotopt_foo}'"
  exit 1
fi
if [[ "${gotopt2_bar}" != "42" ]]; then
  echo "Want: 42; got: '${gotopt_bar}'"
  exit 1
fi
if [[ "${gotopt2_baz}" != "true" ]]; then
  echo "Want: true; got: '${gotopt_baz}'"
  exit 1
fi
