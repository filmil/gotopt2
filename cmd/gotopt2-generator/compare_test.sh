#!/bin/bash
set -euo pipefail

GOTOPT2="$1"
GENERATOR="$2"

TEMP_DIR="test_tmp_$$"
mkdir -p "$TEMP_DIR"
trap "rm -rf '$TEMP_DIR'" EXIT

run_test() {
  local name="$1"
  local config_content="$2"
  shift 2
  local args=("$@")

  echo "Testing: $name"

  local config_file="$TEMP_DIR/config_$name.yaml"
  echo "$config_content" > "$config_file"

  local out_gotopt2="$TEMP_DIR/out_gotopt2_$name"
  local err_gotopt2="$TEMP_DIR/err_gotopt2_$name"
  local out_generator="$TEMP_DIR/out_generator_$name"
  local err_generator="$TEMP_DIR/err_generator_$name"
  local parser_sh="$TEMP_DIR/parser_$name.sh"

  set +e
  "$GOTOPT2" "${args[@]}" < "$config_file" > "$out_gotopt2" 2> "$err_gotopt2"
  local gotopt2_exit=$?

  "$GENERATOR" < "$config_file" > "$parser_sh" 2>/dev/null
  bash -c "source $parser_sh; parse_args \"\$@\"" _ "${args[@]}" > "$out_generator" 2> "$err_generator"
  local generator_exit=$?
  set -e

  if [[ $gotopt2_exit -ne $generator_exit ]]; then
      echo "Exit codes differ for test $name: gotopt2=$gotopt2_exit, generator=$generator_exit"
      exit 1
  fi

  if [[ "$name" == "Help" ]]; then
      # Help output to stderr might be slightly different stylistically, 
      # but we can check if both exited 11 and printed something.
      if [[ $gotopt2_exit -ne 11 ]]; then
         echo "gotopt2 didn't exit with 11 on help!"
         exit 1
      fi
      if ! grep -q "Usage" "$err_generator"; then
         echo "generator didn't print usage!"
         exit 1
      fi
      return 0
  fi

  if ! diff -u "$out_gotopt2" "$out_generator"; then
    echo "Outputs differ for test: $name"
    exit 1
  fi
}

CONFIG_1='
flags:
- name: "foo"
  help: "This is foo"
  type: string
- name: "bool-flag"
  type: bool
- name: "list-flag"
  type: stringlist
'
run_test "Basic" "$CONFIG_1" "--foo=bar" "--bool-flag" "pos1" "pos2"

CONFIG_2='
ALL_CAPS: true
prefix: "my_prefix_"
flags:
- name: "foo-bar"
  help: "This is foo"
  type: string
'
run_test "PrefixAndCaps" "$CONFIG_2" "--foo-bar=baz" "arg1"

CONFIG_3='
flags:
- name: "foo"
  type: string
  help: "Foo configuration"
'
run_test "Help" "$CONFIG_3" "--help"

CONFIG_4='
programName: "my_cool_prog"
flags:
- name: "foo"
  type: string
  help: "Foo configuration"
'
run_test "HelpProgramName" "$CONFIG_4" "--help"

echo "All tests passed."
