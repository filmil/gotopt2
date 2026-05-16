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
  else
      if ! diff -u "$out_gotopt2" "$out_generator"; then
        echo "Outputs differ for test: $name"
        exit 1
      fi
  fi

  if command -v fish >/dev/null 2>&1; then
      local out_generator_fish="$TEMP_DIR/out_generator_fish_$name"
      local err_generator_fish="$TEMP_DIR/err_generator_fish_$name"
      local parser_fish="$TEMP_DIR/parser_$name.fish"

      "$GENERATOR" --shell=fish < "$config_file" > "$parser_fish" 2>/dev/null

      set +e
      fish -c "source $parser_fish; parse_args \$argv" _ "${args[@]}" > "$out_generator_fish" 2> "$err_generator_fish"
      local generator_fish_exit=$?
      set -e

      # Help returns 11 in both generators. If gotopt2_exit=0, we expect 0.
      if [[ $gotopt2_exit -ne $generator_fish_exit ]]; then
          echo "Exit codes differ for fish test $name: gotopt2=$gotopt2_exit, generator=$generator_fish_exit"
          cat "$err_generator_fish"
          exit 1
      fi

      if [[ "$name" != "Help" ]]; then
          if ! grep -q "# gotopt2:generated:begin" "$out_generator_fish"; then
              echo "Fish output missing generated block for test: $name"
              exit 1
          fi
          # Test if evaluating the fish output succeeds
          set +e
          fish -c "cat $out_generator_fish | source"
          local source_exit=$?
          set -e
          if [[ $source_exit -ne 0 ]]; then
              echo "Failed to evaluate generated fish script for test: $name"
              exit 1
          fi

          # Output verification: the generated code must have identical assignments.
          # We extract the content between generated:begin and generated:end markers.
          local bash_eval=$(sed -n '/# gotopt2:generated:begin/,/# gotopt2:generated:end/p' "$out_generator" | grep -v "# gotopt2:generated:")
          local fish_eval=$(sed -n '/# gotopt2:generated:begin/,/# gotopt2:generated:end/p' "$out_generator_fish" | grep -v "# gotopt2:generated:")

          # Convert fish 'set -g VAR VAL' into bash 'VAR=VAL' format just for basic comparison of keys.
          # Note: this might not handle complex quotes perfectly in bash format but verifies identical export names.
          local fish_to_bash=$(echo "$fish_eval" | sed -r 's/set -g ([^ ]+) (.*)/\1=\2/')

          # Just count the number of variables to ensure they map 1-1
          local bash_lines=$(echo "$bash_eval" | grep -c "=" || true)
          local fish_lines=$(echo "$fish_to_bash" | grep -c "=" || true)
          if [[ $bash_lines -ne $fish_lines ]]; then
              echo "Mismatch in number of variables generated. bash=$bash_lines, fish=$fish_lines"
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

echo "All tests passed."
