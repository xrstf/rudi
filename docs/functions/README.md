# Otto Functions

Otto ships with a set of built-in functions. When embedding Otto, this set can
be extended or overwritten as desired to inject custom functions.

## Core

* [`default`](core-default.md) returns a fallback if the given value is
  empty-ish.
* [`do`](core-do.md) evaluates expressions in sequence, sharing a context
  between them.
* [`empty?`](core-empty.md) decides whether a given value is effectively empty
  (for example `0` or `""`).
* [`if`](core-if.md) forms conditions.
* [`try`](core-try.md) returns a fallback value if an expression errors out.
