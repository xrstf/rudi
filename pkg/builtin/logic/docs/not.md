# not

`not` negates the given value.

## Examples

* `(not true)` -> `false`
* `(not false)` -> `true`

## Forms

### `(not expr)`

* `expr` is an arbitrary expressions.

`not` will evaluate the expression and error out if it errors out or its result
cannot be coalesced to a boolean. Otherwise, the negated boolish value is
returned.
