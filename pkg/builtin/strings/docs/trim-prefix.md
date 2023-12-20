# trim-prefix?

`trim-prefix?` checks whether a given string begins with another string and if
it does, returns a copy of the string with the prefix removed. Note that this
removal is done once, so removing the prefix `"o"` from `"oof"` will yield
`"of"`.

## Examples

* `(trim-prefix? "11234" "1")` ➜ `"1234"`
* `(trim-prefix? "11234" "x")` ➜ `"11234"`

## Forms

### `(trim-prefix? base:string prefix:string)` ➜ `string`

* `base` is an arbitrary expression.
* `prefix` is an arbitrary expression.

`trim-prefix?` evaluates the first argument and coalesces it into a string. If
successful, it evalates the prefix the same way. If both values are strings,
the function checks if the base string has the prefix, and if so removes it once
and returns the resulting string.

## Context

`trim-prefix?` executes all expressions in their own contexts, so nothing is
shared.
