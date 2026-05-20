# gotopt2: From Zero to Hero

Welcome to the `gotopt2` tutorial! This guide will take you from the basics of command-line flag parsing in shell scripts to generating standalone parsers.

## What is gotopt2?

If you write shell scripts, you eventually need to accept command-line arguments like `--file=foo.txt` or `--verbose`. Bash doesn't have a great built-in way to parse these safely and cleanly. `gotopt2` solves this by taking a YAML description of your flags and either:
1. Parsing the arguments directly and giving you safe bash code to evaluate.
2. Generating a standalone bash function that you can drop into your script (no `gotopt2` binary needed at runtime).

## Step 1: Your First Script

Let's create a simple greeting script that accepts a `--name` and a `--loud` flag.

Create a file named `greet.sh`:

```bash
#!/bin/bash

# 1. Define your flags in YAML
YAML_CONFIG="
flags:
- name: name
  type: string
  default: World
  help: Who to greet
- name: loud
  type: bool
  help: Yell the greeting
"

# 2. Run gotopt2 and evaluate the output
# We pass the original script arguments ($@) to gotopt2
GOTOPT2_OUTPUT=$(gotopt2 "$@" <<< "$YAML_CONFIG")

# If gotopt2 exited with 11, it means the user asked for --help.
if [[ "$?" == "11" ]]; then
  exit 1
fi

# 3. Evaluate the generated assignments
eval "${GOTOPT2_OUTPUT}"

# 4. Use the variables! (gotopt2 prepends 'gotopt2_' by default)
GREETING="Hello, ${gotopt2_name}!"

if [[ "${gotopt2_loud}" == "true" ]]; then
  echo "${GREETING^^}" # uppercase
else
  echo "${GREETING}"
fi
```

Make it executable and run it:
```console
$ chmod +x greet.sh
$ ./greet.sh
Hello, World!
$ ./greet.sh --name=Alice
Hello, Alice!
$ ./greet.sh --name=Bob --loud
HELLO, BOB!
$ ./greet.sh --help
Usage:
  --name
         Who to greet (default: "World")
  --loud
         Yell the greeting (default: "false")
```

## Step 2: Types and Lists

`gotopt2` supports several types: `string`, `int`, `bool`, and `stringlist`. Let's see `stringlist` in action.

```bash
#!/bin/bash

YAML_CONFIG="
flags:
- name: items
  type: stringlist
  help: A comma-separated list of items
"

eval "$(gotopt2 "$@" <<< "$YAML_CONFIG")"

echo "You provided ${#gotopt2_items__list[@]} items:"
for item in "${gotopt2_items__list[@]}"; do
  echo " - $item"
done
```

Run it:
```console
$ ./list.sh --items=apples,bananas,pears
You provided 3 items:
 - apples
 - bananas
 - pears
```

Notice that `stringlist` variables are named `gotopt2_<name>__list` and are Bash arrays.

## Step 3: Positional Arguments

What about arguments that aren't flags? (e.g., `script.sh --verbose file1.txt file2.txt`)

`gotopt2` collects these into a special array called `gotopt2_args__`.

```bash
#!/bin/bash
YAML_CONFIG="
flags:
- name: verbose
  type: bool
"

eval "$(gotopt2 "$@" <<< "$YAML_CONFIG")"

if [[ "${gotopt2_verbose}" == "true" ]]; then
  echo "Processing files..."
fi

for file in "${gotopt2_args__[@]}"; do
  echo "Working on: $file"
done
```

Run it:
```console
$ ./process.sh --verbose data.csv config.json
Processing files...
Working on: data.csv
Working on: config.json
```

## Step 4: Standalone Parsers (gotopt2-generator)

Sometimes you don't want to force your users to install `gotopt2` just to run your script. You can use `gotopt2-generator` to create a standalone Bash parser!

Create your YAML config: `config.yaml`
```yaml
flags:
- name: port
  type: int
  default: 8080
- name: debug
  type: bool
```

Run the generator:
```console
$ gotopt2-generator < config.yaml > parser.sh
```

Now, source `parser.sh` in your main script:
```bash
#!/bin/bash
source parser.sh

# Call the generated function
parse_args "$@"
if [[ "$?" == "11" ]]; then exit 1; fi

echo "Starting server on port ${gotopt2_port}..."
[[ "${gotopt2_debug}" == "true" ]] && echo "Debug mode enabled!"
```

The resulting script is 100% pure Bash, requires no external dependencies, and handles `--help`, defaults, and positional arguments exactly like the `gotopt2` binary!

## Step 5: Advanced Configuration

You can tweak how `gotopt2` outputs variables at the top level of your YAML.

```yaml
# Make all variables uppercase (e.g., GOTOPT2_PORT)
ALL_CAPS: true

# Change the prefix (e.g., myapp_port instead of gotopt2_port)
prefix: "myapp_"

# Add variable declarations (useful inside bash functions)
declaration: "local"

flags:
- name: port
  type: int
```

With `declaration: "local"`, the output will look like:
```bash
local myapp_port=8080
```
This is perfect for keeping variables scoped safely inside your script's functions.

## Conclusion

You are now a `gotopt2` hero! You can:
- Parse strings, booleans, integers, and lists safely.
- Auto-generate help text.
- Capture positional arguments.
- Compile standalone, dependency-free bash parsers.

Happy scripting!
