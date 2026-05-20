# gotopt2 TODO List

This document tracks potential features, improvements, and maintenance tasks for the `gotopt2` project.

## 1. Generator Improvements (High Impact)
The `gotopt2-generator` allows for zero-dependency shell scripts, but its current implementation is more restrictive than the Go binary.
- [ ] **Support `--flag value` syntax:** Currently, the generated bash function only supports `--flag=value`.
- [ ] **Short flag aliases:** Add support for one-letter aliases (e.g., `-f` for `--foo`).
- [ ] **Positional argument handling:** Improve how remaining arguments are collected and quoted to prevent injection or word-splitting issues in complex cases.
- [ ] **Zsh/Fish Support:** Generate shell-specific completion or parsing logic for other popular shells.

## 2. Core Feature Additions
Enhance the YAML configuration schema to support common CLI patterns.
- [ ] **Required Flags:** Add a `required: true` field to the YAML config to fail if a mandatory flag is missing.
- [ ] **Flag Aliases:** Explicitly support multiple names for the same flag.
- [ ] **Validation Rules:**
    - [ ] `choices`: Restrict a string flag to a specific list of values.
    - [ ] `min`/`max`: Boundary checks for integer flags.
    - [ ] `regex`: Pattern matching for string flags.
- [ ] **Environment Variable Mapping:** Allow flags to automatically pull values from environment variables (e.g., `env: DB_PASSWORD`).

## 3. Tooling & DX (Developer Experience)
- [ ] **`validate` command:** Add a way to validate the YAML configuration file without attempting to parse arguments.
- [ ] **Customizable Help:** Allow users to provide a `header` and `footer` for the `--help` output in the YAML config.
- [ ] **Version Information:** Add a standard `--version` flag to the `gotopt2` and `gotopt2-generator` binaries.
- [ ] **Interactive Mode:** A simple helper to generate the YAML boilerplate based on user prompts.

## 4. Internal Refactoring & Maintenance
- [ ] **Template Cleanup:** The `parser.sh.tmpl` uses hardcoded integer types (e.g., `eq .Type 4`). These should be replaced with string names or symbolic constants for better readability.
- [ ] **Refactor `pkg/opts`:** Separate the configuration parsing from the flag set generation to make it more testable and reusable.
- [ ] **Standardize `stringlist` parsing:** Ensure both the Go binary and the generated bash script handle commas and whitespace in list items identically.

## 5. Documentation & Examples
- [ ] **Advanced Examples:** Add examples for using `gotopt2` with `local` variables in bash functions.
- [x] **Tutorial:** A "From Zero to Hero" guide in `docs/` or an expanded `README.md`.
- [ ] **Comparison Matrix:** Explicitly show the differences between `getopt`, `argbash`, and `gotopt2`.
