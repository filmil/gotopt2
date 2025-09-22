#! /bin/bash
# Copyright (C) 2020 Google Inc.
#
# This file has been licensed under Apache 2.0 license.  Please see the LICENSE
# file at the root of the repository.

# This magic was copied from runfiles by consulting:
#   https://stackoverflow.com/questions/53472993/how-do-i-make-a-bazel-sh-binary-target-depend-on-other-binary-targets

# --- begin runfiles.bash initialization ---
# Copy-pasted from Bazel's Bash runfiles library (tools/bash/runfiles/runfiles.bash).
if [[ ! -d "${RUNFILES_DIR:-/dev/null}" && ! -f "${RUNFILES_MANIFEST_FILE:-/dev/null}" ]]; then
  if [[ -f "$0.runfiles_manifest" ]]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles_manifest"
  elif [[ -f "$0.runfiles/MANIFEST" ]]; then
    export RUNFILES_MANIFEST_FILE="$0.runfiles/MANIFEST"
  elif [[ -f "$0.runfiles/bazel_tools/tools/bash/runfiles/runfiles.bash" ]]; then
    export RUNFILES_DIR="$0.runfiles"
  fi
fi
if [[ -f "${RUNFILES_DIR:-/dev/null}/bazel_tools/tools/bash/runfiles/runfiles.bash" ]]; then
  source "${RUNFILES_DIR}/bazel_tools/tools/bash/runfiles/runfiles.bash"
elif [[ -f "${RUNFILES_MANIFEST_FILE:-/dev/null}" ]]; then
  source "$(grep -m1 "^bazel_tools/tools/bash/runfiles/runfiles.bash " \
            "$RUNFILES_MANIFEST_FILE" | cut -d ' ' -f 2-)"
else
  echo >&2 "ERROR: cannot find @bazel_tools//tools/bash/runfiles:runfiles.bash"
  exit 1
fi
# --- end runfiles.bash initialization ---

# This is somewhat of a hack: we're trying to find the runfiles path of a binary
# whose location is with respect to its position in the dependency repo. So,
# strip 'external/' where it applies.
if [[ "${_binary_path}" == "external/"* ]]; then
    _binary_path="${_binary_path#external/}"
fi
readonly _gotopt2_binary="$(rlocation ${_binary_path})"


# Exit quickly if the binary isn't found. This may happen if the binary location
# moves internally in bazel.
if [[ ! -x "${_gotopt2_binary}" ]]; then
  echo "gotopt2 binary not found"
  exit 240
fi

GOTOPT2_OUTPUT=$($_gotopt2_binary "${@}" <<EOF
flags:
- name: "container"
  type: string
  help: "The name of the container to run"
EOF
)
if [[ "$?" == "11" ]]; then
  # When --help option is used, gotopt2 exits with code 11.
  exit 0
fi

# Evaluate the output of the call to gotopt2, shell vars assignment is here.
eval "${GOTOPT2_OUTPUT}"

