# empty?

`empty?` returns `true` if the given value is "empty-ish". The following values
are considered "empty-ish":

* `null`
* `bool`: `false`
* `number`: `0` (integer) and `0.0` (float)
* `string`: `""` (string of zero length)
* `vector`: `[]` (vector with 0 elements)
* `object`: `{}` (object with 0 elements)

`empty?` is primarily intended to deal with malformed/maltyped JSON data.

## Examples

* `(empty? "")` -> `true`
* `(empty? 42)` -> `false`

## Forms

### `(empty? expr)`

* `expr` is exactly 1 expression.

The expression is evaluated and if it's empty-ish according to the list above,
`true` is returned, else `false`. The expression must be evaluatable, so
identifiers cannot be used (`(empty? to-upper)` is invalid).

If the expression errors out, `empty?` also returns an error.

## Context

`empty?` evaluates the expression in its own scope, so variables defined in it
do not leak.
