# set-symdiff

This function returns a set of values which are in either of the sets, but not
in their intersection.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-symdiff $set (new-set))` ➜ `set{"a", "b", "c"}`
* `(set-symdiff $set (new-set "b"))` ➜ `set{"a", "c"}`
* `(set-symdiff $set (new-set "d"))` ➜ `set{"a", "b", "c", "d"}`

## Forms

### `(set-symdiff base:set other:set)` ➜ `set`

This form returns a set of values which are in either of the sets, but not
in their intersection.
