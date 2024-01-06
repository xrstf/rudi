# Rudi

<p align="center">
  <img src="./docs/rudi-portrait.png" alt="">
</p>

<p align="center">
  <img src="https://img.shields.io/github/v/release/xrstf/rudi" alt="last stable release">

  <a href="https://goreportcard.com/report/go.xrstf.de/rudi">
    <img src="https://goreportcard.com/badge/go.xrstf.de/rudi" alt="go report card">
  </a>

  <a href="https://pkg.go.dev/go.xrstf.de/rudi">
    <img src="https://pkg.go.dev/badge/go.xrstf.de/rudi" alt="godoc">
  </a>
</p>

Rudi is a Lisp-like, embeddable programming language that focuses on transforming data structures
like those available in JSON (numbers, bools, objects, vectors etc.). Rudi programs consist of a
series of statements that are evaluated in order:

```lisp
(if (gt? .replicas 5) (error "too many replicas (%d)" .replicas))
(set! .spec.isAdmin (has-suffix? .spec.Email "@initech.com"))
(map! .spec.usages to-lower)
```

Rudi is great for making tasks like manipulating data, implementing policies, setting default values
or normalizing data configurable.

## Contents

* [Features](#features)
* [Installation](#installation)
* [Documentation](#documentation)
  * [Language Description](docs/language.md)
  * [Type Handling](docs/coalescing.md)
  * [Standard Library](docs/stdlib/README.md)
  * [Extended Library](docs/extlib/README.md)
* [Usage](#usage)
  * [Command Line](#command-line)
  * [Embedding](#embedding)
* [Alternatives](#alternatives)
* [Credits](#credits)
* [License](#license)

## Features

* **Safe** evaluation: Rudi is not Turing-complete and so Rudi programs are always guaranteed to
  complete in a reasonable time frame. You can add support for functions defined in Rudi code (i.e.
  at runtime), but this is optional to allow a safe embedded behaviour by default.
* **Lightweight**: Rudi comes with no runtime dependencies besides the Go stdlib.
* **Hackable**: Rudi tries to keep the language itself approachable, so that modifications are
  easier and newcomers have an easier time to get started.
* **Variables** can be pre-defined or set at runtime.
* **JSONPath** expressions are first-class citizens and make referring to the current JSON document
  a breeze.
* **Optional Type Safety**: Choose between pedantic, strict or humane typing for your programs.
  Strict allows nearly no type conversions, humane allows for things like `1` (int) turning into
  `"1"` (string) when needed.
* **Flexible**: The Rudi CLI interpreter (`rudi`) supports reading/writing JSON,
  [JSON5](https://json5.org/), [YAML](https://yaml.org/) and [TOML](https://toml.io/en/).

## Installation

Rudi is primarily meant to be embedded into other Go programs, but a standalone CLI application,
`rudi`, is also available to test your scripts with. `rudi` can be installed using Git & Go. Rudi
requires **Go 1.18** or newer.

```bash
git clone https://github.com/xrstf/rudi
cd rudi
make build
```

This will result in a `rudi` binary in `_build/`; to install system-wide, use `make install`.

Alternatively, you can download the [latest release](https://github.com/xrstf/rudi/releases/latest)
from GitHub.

## Documentation

Make yourself familiar with Rudi using the documentation:

* The [Language Description](docs/language.md) describes the Rudi syntax and semantics.
* All built-in functions are described in the [standard library](docs/stdlib/README.md).
* Additional functions are available in the [extended library](docs/extlib/README.md).
* [Type Handling](docs/coalescing.md) describes how Rudi handles, converts and compares values.

## Usage

### Command Line

Rudi comes with a standalone CLI tool called `rudi`.

```
Usage of rudi:
  -i, --interactive            Start an interactive REPL to run expressions.
  -s, --script string          Load Rudi script from file instead of first argument (only in non-interactive mode).
  -l, --library stringArray    Load additional Rudi file(s) to be be evaluated before the script (can be given multiple times).
      --var stringArray        Define additional global variables (can be given multiple times).
  -f, --stdin-format string    What data format is used for data provided on stdin, one of [raw json json5 yaml yamldocs toml]. (default "yaml")
  -o, --output-format string   What data format to use for outputting data, one of [raw json yaml yamldocs toml]. (default "json")
      --enable-funcs           Enable the func! function to allow defining new functions in Rudi code.
  -c, --coalesce string        Type conversion handling, one of [strict pedantic humane]. (default "strict")
  -h, --help                   Show help and documentation.
  -V, --version                Show version and exit.
      --debug-ast              Output syntax tree of the parsed script in non-interactive mode.
```

`rudi` can run in one of two modes:

* **Interactive Mode** is enabled by passing `--interactive` (or `-i`). This will start a REPL
  session where Rudi scripts are read from stdin and evaluated against the loaded files.
* **Script Mode** is used the an Rudi script is passed either as the first argument or read from a
  file defined by `--script`. In this mode `rudi` will run all statements from the script and print
  the resulting value, then it exits.

    Examples:

    * `rudi '.foo' myfile.json`
    * `rudi '(set! .foo "bar") (set! .users 42) .' myfile.json`
    * `rudi --script convert.rudi myfile.json`

`rudi` has extensive help built right into it, try running `rudi help` to get started.

#### File Handling

Rudi can load JSON, JSON5, YAML and TOML files and will determine the file format based on the
file extension (`.json` for JSON, `.json5` for JSON5, `.yml` and `.yaml` for YAML and `.tml` /
`.toml` for TOML). For data provided via stdin, `rudi` by default assumes YAML (or JSON) encoding.
If you want to use TOML/JSON5 instead, you must use the `--stdin-format` flag.

The first loaded file is known as the "document". Its content is available via path expressions like
`.foo[0]`. All loaded files are also available via the `$files` variable (i.e. `.` is the same as
`$files[0]` for reading, but when writing data, there is a difference between both notations; refer
to the docs for `set` for more information). Additionally the filenames are available in the
`$filenames` variable.

Additional raw files can be loaded using the `--var` flag: To load files, the format for this flag
is `ENCODING:file:FILENAME`, for example `--var "myvar=yaml:file:config.kubeconfig"`. This allows
you to load files regardless of their extension and also allows to load raw files (that will be
kept as strings) using `"myvar=raw:file:logo.png"`. Raw file encoding is not supported for files
given as arguments, those files must have a recognized file extension.

### Embedding

Rudi is well suited to be embedded into Go applications. A clean and simple API makes it a breeze:

```go
package main

import (
   "fmt"
   "log"

   "go.xrstf.de/rudi"
   "go.xrstf.de/rudi/pkg/coalescing"
)

const script = `(set! .foo 42) (+ $myvar 42 .foo)`

func main() {
   // Rudi programs are meant to manipulate a document (path expressions like
   // ".foo" resolve within that document). The document can be anything,
   // but is most often a JSON object.
   documentData := map[string]any{"foo": 9000}

   // parse the script (the name is used when generating error strings)
   program, err := rudi.Parse("myscript", script)
   if err != nil {
      log.Fatalf("The script is invalid: %v", err)
   }

   // evaluate the program;
   // this returns an evaluated value, which is the result of the last expression
   // that was evaluated, plus the final document state (the updatedData) after
   // the script has finished.
   updatedData, result, err := program.Run(
      context.Background(),
      documentData,
      // setup the set of variables available by default in the script
      rudi.NewVariables().Set("myvar", 42),
      // Likewise, setup the functions available (note that this includes
      // functions like "if" and "and", so running with an empty function set
      // is generally not advisable).
      rudi.NewSafeBuiltInFunctions(),
      // Decide what kind of type strictness you would like; pedantic, strict
      // or humane; choose your own adventure (strict is default if you use nil
      // here; humane allows conversions like 1 == "1").
      coalescing.NewStrict(),
   )
   if err != nil {
      log.Fatalf("Script failed: %v", err)
   }

   fmt.Println(result)       // => 126
   fmt.Println(updatedData)  // => {"foo": 42}
}
```

## Alternatives

Rudi doesn't exist in a vacuum; there are many other great embeddable programming/scripting languages
out there, allbeit with slightly different ideas and goals than Rudi:

* [Anko](https://github.com/mattn/anko) – Go-like syntax and allows recursion, making it more
  dangerous and hard to learn for non-developers than I'd like.
* [ECAL](https://github.com/krotik/ecal) – Is an event-based system using rules which are triggered by
  events; comes with recursion as well and is therefore out.
* [Expr](https://github.com/antonmedv/expr), [GVal](https://github.com/PaesslerAG/gval),
  [CEL](https://github.com/google/cel-go) – Great languages for writing a single expression, but not
  suitable for transforming/mutating data structures.
* [Gentee](https://github.com/gentee/gentee) – Is similar to C/Python and allows recursion, so both
  to powerful/dangerous and not my preference in terms of syntax.
* [Jsonnet](https://github.com/google/go-jsonnet) – Probably one of the most obvious alternatives
  among this list. Jsonnet shines when constructing new elements and complexer configurations
  out of smaller pieces of information, less so when manipulating objects. Also I personally really
  am no fan of Jsonnet's syntax, plus: NIH.
* [Starlark](https://github.com/google/starlark-go) – Is the language behind Bazel and actually has
  an optional nun-Turing-complete mode. However I am really no fan of its syntax and have not
  investigated it further.
* [Go Templates](https://pkg.go.dev/text/template) – I really don't like Go's template syntax for
  more than simple one-liners. I liked and copied its concept of ranging over things, as templates
  do not allow unbounded loops (just like Rudi), but apart from being safe to embed, Go templates do
  not offer enough functionality to modify a data structure. Like Jsonnet, templates shine when
  creating/outputting entire _new_ documents.

  _Bonus mention:_ Mastermind's [sprig](https://github.com/Masterminds/sprig) served as inspiration
  for quite a few of the functions in Rudi.

## Credits

Rudi has been named after my grandfather.

Thanks to [@embik](https://github.com/embik) and [@xmudrii](https://github.com/xmudrii) for enduring
my constant questions for feedback :smile:

Rudi has been made possible by the amazing [Pigeon](https://github.com/mna/pigeon) parser generator.

## License

MIT
