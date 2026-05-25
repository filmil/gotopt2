#!/bin/bash -x
/bin/bash
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

  "$GENERATOR" < "$config_file" > "$parser_sh" 
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
  else
      local env_gotopt2="$TEMP_DIR/env_gotopt2_$name"
      local env_generator="$TEMP_DIR/env_generator_$name"

      # We ignore standard output from the generator, but we extract env vars
      bash -c "source \"$out_gotopt2\"; declare -p | grep -E '(gotopt2_|my_prefix)' | grep -v BASH_EXECUTION_STRING | grep -v '_=' | sort" > "$env_gotopt2"
      bash -c "source \"$parser_sh\"; parse_args \"\$@\"; declare -p | grep -E '(gotopt2_|my_prefix)' | grep -v BASH_EXECUTION_STRING | grep -v '_=' | sort" _ "${args[@]}" > "$env_generator"

      if ! diff -u "$env_gotopt2" "$env_generator"; then
        echo "Environment variables differ for test: $name"
        exit 1
      fi
  fi

  if command -v fish >/dev/null 2>&1; then
      local err_generator_fish="$TEMP_DIR/err_generator_fish_$name"
      local parser_fish="$TEMP_DIR/parser_$name.fish"

      "$GENERATOR" --shell=fish < "$config_file" > "$parser_fish" 

      set +e
      # Run it once just to get exit code and stderr
      fish -c "source $parser_fish; parse_args \$argv" -- "${args[@]}" > /dev/null 2> "$err_generator_fish"
      local generator_fish_exit=$?
      set -e

      # Help returns 11 in both generators. If gotopt2_exit=0, we expect 0.
      if [[ $gotopt2_exit -ne $generator_fish_exit ]]; then
          echo "Exit codes differ for fish test $name: gotopt2=$gotopt2_exit, generator=$generator_fish_exit"
          cat "$err_generator_fish"
          exit 1
      fi

      if [[ "$name" != "Help" ]]; then
          local env_gotopt2_fish="$TEMP_DIR/env_gotopt2_fish_$name"
          local env_generator_fish="$TEMP_DIR/env_generator_fish_$name"

          # Convert bash variables from gotopt2 into fish vars
          # (For testing we just check that the variable keys are populated similarly)
          fish -c "cat $out_gotopt2 | source; set | grep -E '^(gotopt2_|my_prefix)' | sort" > "$env_gotopt2_fish"
          fish -c "source $parser_fish; parse_args \$argv; set | grep -E '^(gotopt2_|my_prefix)' | sort" -- "${args[@]}" > "$env_generator_fish"

          if ! diff -u "$env_gotopt2_fish" "$env_generator_fish"; then
              echo "Fish environment variables differ for test: $name"
              exit 1
          fi
      fi
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
trueValue: "on"
flags:
- name: "mybool"
  type: bool
'
run_test "TrueValue" "$CONFIG_4" "--mybool"


CONFIG_5='
flags:
- name: "foo"
  help: "This is foo"
  type: string
- name: "bool-flag"
  type: bool
- name: "list-flag"
  type: stringlist
'
run_test "SpaceSyntax" "$CONFIG_5" "--foo" "bar" "--bool-flag" "pos1" "pos2"

CONFIG_6='
flags:
- name: "list-flag"
  type: stringlist
'
run_test "ListSpaceSyntax" "$CONFIG_6" "--list-flag" "a,b,c"

echo "All tests passed."
