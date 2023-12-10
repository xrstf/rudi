# or

`or` returns the logical disjunction of all given arguments (i.e.
`a || b || c || …`). Each of the arguments must be convertible to booleans.

## Examples

* `(or false false)` -> `false`
* `(or false false true)` -> `true`

## Forms

### `(or expr+)`

* `expr` is one or more arbitrary expressions.

`or` will evaluate each expression in order, stopping at the first that either
errors out or is coalesced to `true`. If coalescing is not possible (e.g. when
using strict mode, `(or 1)` is invalid), an error is returned.

At least one of the expressions must coalesce to `true` for the function to
return `true`.

## Context

`or` executes all expressions in their own contexts, so nothing is shared.
