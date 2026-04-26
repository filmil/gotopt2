## 2025-01-31 - [Command Injection via Double Quotes in eval Payload]
**Vulnerability:** Command injection was possible because `gotopt2` output strings wrapped in double quotes (e.g., `gotopt2_var="$(echo hacked)"`), which caused the host bash script using `eval` to evaluate the inner command substitution.
**Learning:** When generating shell code intended for `eval`, double quotes are unsafe if the content includes user input, because they allow parameter expansion and command substitution.
**Prevention:** Always wrap generated string values in single quotes, and escape internal single quotes as `'"'"'` (or similar), to ensure the string is treated purely as a literal value by the shell parser.
## 2025-02-05 - [Command Injection via Configuration Parameters in eval Payload]
**Vulnerability:** Command injection was possible because `gotopt2` output bash variables without sanitizing the `name`, `prefix`, and `declaration` configuration inputs. An attacker could inject arbitrary commands by providing values like `foo;echo INJECTED`.
**Learning:** Even configurable prefixes and variable names can be vectors for command injection if the output is evaluated (`eval`). Unsanitized configuration strings can break out of variable assignment contexts.
**Prevention:** Strictly sanitize all variables, prefixes, and declarations using an allowlist approach to ensure only valid Bash identifier characters (`[a-zA-Z0-9_]`) or specific declaration characters (`[a-zA-Z0-9_ \-]`) are emitted into shell scripts.
