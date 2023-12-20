# lte?

`lte?` returns whether the first argument is numerically smaller than or equal to
the second.

## Examples

* `(lte? 3 2)` ➜ `true`
* `(lte? 2.0 3)` ➜ `false`

## Forms

### `(lte? left:any right:any)` ➜ `bool`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

`lte?` evaluates both expressions and coalesces their results to be numbers. If
either evaluation or conversion fail, an error is returned. If the two arguments
are valid numbers, the result of `left <= right` is returned.

## Context

`lte?` executes both expressions in their own contexts, so nothing is shared.
