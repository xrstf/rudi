# trim

`trim` removes all whitepace from the beginning and end of a string. This
includes linebreaks.

## Examples

* `(trim " hello\nworld ")` ➜ `"hello\nworld"`
* `(trim "\n")` ➜ `""`

## Forms

### `(trim value:string)` ➜ `string`

* `value` is an arbitrary expression.

Returns a copy of the string with leading and trailing whitespace removed.
