# set-union

This function returns the union of two or more sets.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-union $set (new-set))` ➜ `set{"a", "b", "c"}`
* `(set-union $set (new-set "b"))` ➜ `set{"a", "b", "c"}`
* `(set-union $set (new-set "d" "e"))` ➜ `set{"a", "b", "c", "d" "e"}`

## Forms

### `(set-union base:set other:set+)` ➜ `set`

This form returns a new set that contains all values that exist in any of the
given sets. More than two sets may be merged into a union at the same time.
