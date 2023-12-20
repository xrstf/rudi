# trim-suffix?

`trim-suffix?` checks whether a given string ends with another string and if
it does, returns a copy of the string with the suffix removed. Note that this
removal is done once, so removing the suffix `"o"` from `"foo"` will yield
`"fo"`.

## Examples

* `(trim-suffix? "12344" "4")` ➜ `"1234"`
* `(trim-suffix? "12344" "x")` ➜ `"12344"`

## Forms

### `(trim-suffix? base:string suffix:string)` ➜ `string`

* `base` is an arbitrary expression.
* `suffix` is an arbitrary expression.

`trim-suffix?` evaluates the first argument and coalesces it into a string. If
successful, it evalates the suffix the same way. If both values are strings,
the function checks if the base string has the suffix, and if so removes it once
and returns the resulting string.

## Context

`trim-suffix?` executes all expressions in their own contexts, so nothing is
shared.
