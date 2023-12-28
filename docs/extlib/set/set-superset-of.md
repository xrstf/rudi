# set-superset-of?

This function returns true if the first set is a superset of the second set,
meaning the first set contains at least all values of the second set.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-superset-of? (new-set) (new-set))` ➜ `true`
* `(set-superset-of? $set (new-set))` ➜ `true`
* `(set-superset-of? $set (new-set "b"))` ➜ `true`
* `(set-superset-of? $set (new-set "d"))` ➜ `false`

## Forms

### `(set-superset-of? base:set other:set)` ➜ `bool`

This form returns true if `base` is a superset of `other`.
