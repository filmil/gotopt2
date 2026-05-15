#!/bin/bash
# A test script to verify my_script.sh

SCRIPT_PATH="${1}"
PARSER_PATH="${2}"

# Run the script with the parser path and a test flag
OUTPUT=$("${SCRIPT_PATH}" "${PARSER_PATH}" --my-flag=test_value)

if [[ "${OUTPUT}" == *"Success: test_value"* ]]; then
  exit 0
else
  echo "Expected to contain 'Success: test_value', got '${OUTPUT}'"
  exit 1
fi
