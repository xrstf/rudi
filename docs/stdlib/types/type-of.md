# type-of

`type-of` returns the type name of a given value.

## Examples

* `(type-of null)` -> `"null"`
* `(type-of true)` -> `"bool"`
* `(type-of 1)` -> `"number"`
* `(type-of 1.2)` -> `"number"`
* `(type-of "")` -> `"string"`
* `(type-of [])` -> `"vector"`
* `(type-of {})` -> `"object"`

## Forms

### `(type-of value)`

* `value` is an arbitrary expression.

`type-of` evaluates the given expression and then returns the datatype name of
the resulting value.
