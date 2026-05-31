#!/bin/bash
# A test script that confirms that the usage string is printed when --help is used.

# Source the generated parser
source "${1}"

# Call the generated parse_args function with --help.
# It should exit with code 11 and print the usage string.
output=$(parse_args --help 2>&1)
exit_code=$?

echo "Captured output: $output"
echo "Exit code: $exit_code"

if [[ $exit_code -ne 11 ]]; then
  echo "Expected exit code 11, got $exit_code"
  exit 1
fi

if [[ "$output" == *"Usage: my_script.sh [options]"* ]]; then
  echo "Success: Usage string found"
  exit 0
else
  echo "Failure: Usage string not found"
  exit 1
fi
