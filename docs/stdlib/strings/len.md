# len

`len` return the length of a string/vector or size of an object. The size of an
object is the number of key-value pairs in it.

## Examples

* `(len "")` ➜ `0`
* `(len " hello ")` ➜ `7`
* `(len [1 2 3])` ➜ `3`
* `(len {foo "bar" hello "world"})` ➜ `2`

## Forms

### `(len value:string)` ➜ `int`

* `value` is an arbitrary expression.

This form returns the length of the string (number of bytes).

### `(len value:vector)` ➜ `int`

* `value` is an arbitrary expression.

This form returns the number of elements in the given vector.

### `(len value:object)` ➜ `int`

* `value` is an arbitrary expression.

This form returns the number of key-value pairs in the given object.
