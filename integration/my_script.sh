#!/bin/bash
# A script that uses the generated parser

# Source the generated parser
source "${1}"
shift

# Call the generated parse_args function
parse_args "$@"

# Check if the parsing succeeded and do something with the output
if [ -n "${gotopt2_my_flag:-}" ]; then
  echo "Success: ${gotopt2_my_flag}"
else
  echo "Failure: my_flag was not set"
  exit 1
fi
