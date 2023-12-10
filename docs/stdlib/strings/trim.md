# trim

`trim` removes all whitepace from the beginning and end of a string. This
includes linebreaks.

## Examples

* `(trim " hello\nworld ")` -> `"hello\nworld"`
* `(trim "\n")` -> `""`

## Forms

### `(trim string)`

* `string` is an arbitrary expression.

`trim` evaluates the first argument and coalesces the result into a string. When
successful, leading and trailing whitespace is removed from the string and then
the resulting string is returned.
