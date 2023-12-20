# and

`and` returns the logical conjunction of all given arguments (i.e.
`a && b && c && …`). Each of the arguments must be convertible to booleans.

## Examples

* `(and true true)` ➜ `true`
* `(and true false (and true true))` ➜ `false`

## Forms

### `(and expr:bool…)` ➜ `bool`

* `expr` is one or more arbitrary expressions.

`and` will evaluate each expression in order, stopping at the first that either
errors out or is coalesced to `false`. If coalescing is not possible (e.g. when
using strict mode, `(and 1)` is invalid), an error is returned.

All of the expressions must coalesce to `true` for this function to return
`true`.

## Context

`and` executes all expressions in their own contexts, so nothing is shared.
