# gte?

`gte?` returns whether the first argument is numerically larger than or equal to
the second.

## Examples

* `(gte? 3 2)` ➜ `true`
* `(gte? 2.0 3)` ➜ `false`

## Forms

### `(gte? left:any right:any)` ➜ `bool`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

`gte?` evaluates both expressions and coalesces their results to be numbers. If
either evaluation or conversion fail, an error is returned. If the two arguments
are valid numbers, the result of `left >= right` is returned.

## Context

`gte?` executes both expressions in their own contexts, so nothing is shared.
