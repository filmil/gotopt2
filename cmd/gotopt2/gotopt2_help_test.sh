#! /bin/bash
# Args:
#   $1: the name of the gotopt2 binary to execute: this is a bazel quirk.
#   <rest>: the flag arguments to parse, see the BUILD rule for the passed args.

# This test script requires that the name of the "gotopt2" binary be the first
# arg in use.
readonly GOTOPT2="${1}"
shift

GOTOPT2_OUTPUT="$("${GOTOPT2}" "${@}" <<EOF
flags:
- name: "foo"
  type: string
  default: "something"
EOF
)"
EXIT_CODE="$?"
echo "exit code: ${EXIT_CODE}"
if [[ "${EXIT_CODE}" != "11" ]]; then
  echo "Exit status was: ${EXIT_CODE}, expected 11"
  exit 1
fi

