# The Otto interpreter :)

Otti is a command line interpreter for the Otto programming language. Otti can
read multiple JSON/YAML files and then apply JSON paths or scripts to them. For
quicker development, an interactive REPL is also available.

## Modes

Otti can run in one of two modes:

* **Interactive Mode** is enabled by passing `--interactive` (or `-i`). This will
  start a REPL session where Otto scripts are read from stdin and evaluated
  against the loaded files.
* **Script Mode** is used the an Otto script is passed either as the first
  argument or read from a file defined by `--script`. In this mode Otti will
  run all statements from the script and print the resulting value, then it exits.

    Examples:

    * `otti '.foo' myfile.json`
    * `otti '(set .foo "bar") (set .users 42) .' myfile.json`
    * `otti --script convert.otto myfile.json`

## File Handling

The first loaded file is known as the "document". Its content is available via
path expressions like `.foo[0]`. All loaded files are also available via the
`$files` variable (i.e. `.` is the same as `$files[0]` for reading, but when
writing data, there is a difference between both notations; refer to the docs
for `set` for more information).

## Help

Help is available by using `help` as the first argument to Otto. This can be
followed by a topic, like `help if`.
