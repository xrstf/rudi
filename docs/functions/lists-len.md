# len

`len` return the length of a string/vector or size of an object. The size of an
object is the number of key-value pairs in it.

## Examples

* `(len "")` -> `0`
* `(len " hello ")` -> `7`
* `(len [1 2 3])` -> `3`
* `(len {foo "bar" hello "world"})` -> `2`

## Forms

### `(len value)`

* `value` is an arbitrary expression.

`len` evaluates the value and depending on what it coalesces to, returns either
the length of the string, length of a vector or size of an object. If the value
cannot be coalesced into a suitable type, an error is returned.
