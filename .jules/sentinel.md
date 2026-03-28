## 2025-01-31 - [Command Injection via Double Quotes in eval Payload]
**Vulnerability:** Command injection was possible because `gotopt2` output strings wrapped in double quotes (e.g., `gotopt2_var="$(echo hacked)"`), which caused the host bash script using `eval` to evaluate the inner command substitution.
**Learning:** When generating shell code intended for `eval`, double quotes are unsafe if the content includes user input, because they allow parameter expansion and command substitution.
**Prevention:** Always wrap generated string values in single quotes, and escape internal single quotes as `'"'"'` (or similar), to ensure the string is treated purely as a literal value by the shell parser.

## 2025-04-14 - [Command Injection via Unsanitized Variable Names in eval Payload]
**Vulnerability:** Malicious input in configuration fields like `name`, `prefix`, and `declaration` could inject arbitrary shell commands (e.g., `prefix: "p; echo HACKED; "`) because these fields were used unsanitized in the generated bash script that is typically executed via `eval`.
**Learning:** When generating bash code from user-provided input, not only the values but also the variable names and command prefixes must be sanitized to prevent command injection.
**Prevention:** Sanitize all components used in generating shell scripts (variable names, prefixes, declaration keywords) to only include safe characters like alphanumeric and underscores ([a-zA-Z0-9_]).
