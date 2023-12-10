# reverse

`reverse` returns a copy of the given string/vector with the characters/items
inreverse order.

## Examples

* `(reverse "abc")` -> `"cba"`
* `(reverse [1 2 3])` -> `[3 2 1]`

## Forms

### `(reverse source)`

* `source` is an arbitrary expression.

`reverse` evaluates the source expression. If tries to coalesce the value to
string and if successful, returns a reverse of the string. If unsuccessful, it
tries to coalesce to a vector and if successful, returns a copy of the vector
with items in reverse order. Otherwise an error is returned.
