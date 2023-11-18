# Otto

<p align="center">
  <img src="https://img.shields.io/github/v/release/xrstf/otto" alt="last stable release">

  <a href="https://goreportcard.com/report/go.xrstf.de/otto">
    <img src="https://goreportcard.com/badge/go.xrstf.de/otto" alt="go report card">
  </a>

  <a href="https://pkg.go.dev/go.xrstf.de/otto">
    <img src="https://pkg.go.dev/badge/go.xrstf.de/otto" alt="godoc">
  </a>
</p>

Otto is a Lisp-based, embeddable programming language that focuses on transforming data structures
like those available in JSON (numbers, bools, objects, vectors etc.). A statement in Otto looks like

```lisp
(set .foo[0] (+ (len .users) 42))
```

## Installation

Otto is primarily meant to be embedded into other Go programs, but a standalone CLI application,
_Otti_, is also available to test your scripts with. Otti can be installed using Go:

```bash
go install go.xrstf.de/otto/cmd/otti
```

Alternatively, you can download the [latest release](https://github.com/xrstf/otto/releases/latest)
from GitHub.

## Embedding

Otto is well suited to be embedded into Go applications. A clean and simple API makes it a breeze:

```go
package main

import (
   "fmt"
   "log"

   "go.xrstf.de/otto"
)

const script = `(+ $myvar 42 .foo)`

func main() {
   // setup the set of variables available by default in the script
   vars := otto.NewVariables().
      Set("myvar", 42)

   // Likewise, setup the functions available (note that this includes functions like "if" and "and",
   // so running with an empty function set is generally not advisable).
   funcs := otto.NewBuiltInFunctions()

   // Otto programs are meant to manipulate a document (path expressions like ".foo" resolve within
   // that document). The document can be anything, but is most often a JSON object.
   documentData := map[string]any{"foo": 9000}
   document, err := otto.NewDocument(documentData)
   if err != nil {
      log.Fatalf("Cannot use %v as the document: %v", documentData, err)
   }

   // parse the script (the name is used when generating error strings)
   program, err := otto.ParseScript("myscript", script)
   if err != nil {
      log.Fatalf("The script is invalid: %v", err)
   }

   // evaluate the program;
   // this returns an evaluated value, which is the result of the last expression that was evaluated,
   // plus a new context, which contains for example newly set runtime variables; in many cases the
   // new context is not that important and you'd focus on the evaluated value.
   newCtx, evaluated, err := otto.RunProgram(ctx, program)
   if err != nil {
      log.Fatalf("Failed to evaluate script: %v", err)
   }

   fmt.Println(evaluated)
   fmt.Println(newCtx)
}
```

## License

MIT
