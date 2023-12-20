# has-suffix?

`has-suffix?` checks whether a given string ends with another string.

## Examples

* `(has-suffix? "foobar" "bar")` ➜ `true`
* `(has-suffix? "foobar" "f")` ➜ `false`

## Forms

### `(has-suffix? base:string suffix:string)` ➜ `bool`

* `base` is an arbitrary expression.
* `suffix` is an arbitrary expression.

`has-suffix?` evaluates the first argument and coalesces it into a string. If
successful, it evalates the suffix the same way. If both values are strings,
the function returns true if `base` ends with `suffix`.

## Context

`has-suffix?` executes all expressions in their own contexts, so nothing is
shared.
