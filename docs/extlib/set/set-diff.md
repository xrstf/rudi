# set-diff

This function returns the difference between two sets.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-diff $set (new-set))` ➜ `set{"a", "b", "c"}`
* `(set-diff $set (new-set "b"))` ➜ `set{"a", "c"}`
* `(set-diff $set (new-set "d"))` ➜ `set{"a", "b", "c"}`

## Forms

### `(set-diff base:set other:set)` ➜ `set`

This form returns `base - other`, i.a. a new set that contains all values that
are not part of the `other` set.
