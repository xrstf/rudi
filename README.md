# Rudi

<p align="center">
  <img src="https://img.shields.io/github/v/release/xrstf/rudi" alt="last stable release">

  <a href="https://goreportcard.com/report/go.xrstf.de/rudi">
    <img src="https://goreportcard.com/badge/go.xrstf.de/rudi" alt="go report card">
  </a>

  <a href="https://pkg.go.dev/go.xrstf.de/rudi">
    <img src="https://pkg.go.dev/badge/go.xrstf.de/rudi" alt="godoc">
  </a>
</p>

Rudi is a Lisp-based, embeddable programming language that focuses on transforming data structures
like those available in JSON (numbers, bools, objects, vectors etc.). A statement in Rudi looks like

```lisp
(set .foo[0] (+ (len .users) 42))
```

Rudi has been named after my grandfather.

## Installation

Rudi is primarily meant to be embedded into other Go programs, but a standalone CLI application,
_Rudi_, is also available to test your scripts with. Rudi can be installed using Git & Go:

```bash
git clone https://github.com/xrstf/rudi
cd rudi
make build
```

Alternatively, you can download the [latest release](https://github.com/xrstf/rudi/releases/latest)
from GitHub.

## Usage

Rudi has extensive help built right into it, try running `rudi help` to get started.

## Embedding

Rudi is well suited to be embedded into Go applications. A clean and simple API makes it a breeze:

```go
package main

import (
   "fmt"
   "log"

   "go.xrstf.de/rudi"
)

const script = `(+ $myvar 42 .foo)`

func main() {
   // setup the set of variables available by default in the script
   vars := rudi.NewVariables().
      Set("myvar", 42)

   // Likewise, setup the functions available (note that this includes functions like "if" and "and",
   // so running with an empty function set is generally not advisable).
   funcs := rudi.NewBuiltInFunctions()

   // Rudi programs are meant to manipulate a document (path expressions like ".foo" resolve within
   // that document). The document can be anything, but is most often a JSON object.
   documentData := map[string]any{"foo": 9000}
   document, err := rudi.NewDocument(documentData)
   if err != nil {
      log.Fatalf("Cannot use %v as the document: %v", documentData, err)
   }

   // combine document, variables and functions into an execution context
   ctx := rudi.NewContext(document, funcs, vars)

   // parse the script (the name is used when generating error strings)
   program, err := rudi.ParseScript("myscript", script)
   if err != nil {
      log.Fatalf("The script is invalid: %v", err)
   }

   // evaluate the program;
   // this returns an evaluated value, which is the result of the last expression that was evaluated,
   // plus a new context, which contains for example newly set runtime variables; in many cases the
   // new context is not that important and you'd focus on the evaluated value.
   newCtx, evaluated, err := rudi.RunProgram(ctx, program)
   if err != nil {
      log.Fatalf("Failed to evaluate script: %v", err)
   }

   fmt.Println(evaluated)
   fmt.Println(newCtx)
}
```

## Credits

Thanks to [@embik](https://github.com/embik) and [@xmudrii](https://github.com/xmudrii) for enduring
my constant questions for feedback :smile:

## License

MIT
