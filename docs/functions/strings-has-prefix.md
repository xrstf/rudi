# has-prefix?

`has-prefix?` checks whether a given string begins with another string.

## Examples

* `(has-prefix? "hello" "hell")` -> `true`
* `(has-prefix? "hallo" "halloween")` -> `false`
* `(has-prefix? "foo" "bar")` -> `false`

## Forms

### `(has-prefix? base prefix)`

* `base` is an arbitrary expression.
* `prefix` is an arbitrary expression.

`has-prefix?` evaluates the first argument and coalesces it into a string. If
successful, it evalates the prefix the same way. If both values are strings,
the function returns true if `base` begins with `prefix`.

## Context

`has-prefix?` executes all expressions in their own contexts, so nothing is
shared.
