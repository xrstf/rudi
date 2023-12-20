# reverse

`reverse` returns a copy of the given string/vector with the characters/items
inreverse order.

## Examples

* `(reverse "abc")` ➜ `"cba"`
* `(reverse [1 2 3])` ➜ `[3 2 1]`

## Forms

### `(reverse source:string)` ➜ `string`

* `source` is an arbitrary expression.

Returns the reverse of the input strings, i.e. a string with the order of all
bytes flipped.

### `(reverse source:vector)` ➜ `vector`

* `source` is an arbitrary expression.

Returns a copy of the source vector with items in the opposite order of the
source.
