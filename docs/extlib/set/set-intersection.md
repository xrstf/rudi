# set-intersection

This function returns the intersection of two sets.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-intersection $set (new-set))` ➜ `set{}`
* `(set-intersection $set (new-set "b"))` ➜ `set{"b"}`
* `(set-intersection $set (new-set "d"))` ➜ `set{}`

## Forms

### `(set-intersection base:set other:set)` ➜ `set`

This form returns a new set that contains all values that exist in both sets.
