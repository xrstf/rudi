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

	fmt.Println(result)      // => 126
	fmt.Println(updatedData) // => {"foo": 42}
}
