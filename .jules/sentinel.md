## 2025-01-31 - [Command Injection via Double Quotes in eval Payload]
**Vulnerability:** Command injection was possible because `gotopt2` output strings wrapped in double quotes (e.g., `gotopt2_var="$(echo hacked)"`), which caused the host bash script using `eval` to evaluate the inner command substitution.
**Learning:** When generating shell code intended for `eval`, double quotes are unsafe if the content includes user input, because they allow parameter expansion and command substitution.
**Prevention:** Always wrap generated string values in single quotes, and escape internal single quotes as `'"'"'` (or similar), to ensure the string is treated purely as a literal value by the shell parser.
