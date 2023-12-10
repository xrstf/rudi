# gt?

`gt?` returns whether the first argument is numerically larger than the second.

## Examples

* `(gt? 3 2)` -> `true`
* `(gt? 2.0 3)` -> `false`

## Forms

### `(gt? left right)`

* `left` is an arbitrary expression, except for identifiers.
* `right` is likewise an arbitrary expression, except for identifiers.

`gt?` evaluates both expressions and coalesces their results to be numbers. If
either evaluation or conversion fail, an error is returned. If the two arguments
are valid numbers, the result of `left > right` is returned.

## Context

`gt?` executes both expressions in their own contexts, so nothing is shared.
