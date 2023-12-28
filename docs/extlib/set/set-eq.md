# set-eq?

This function returns true if two sets are identical, i.e. contain the exact
same values.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-eq? $set (new-set))` ➜ `false`
* `(set-eq? $set (new-set "b"))` ➜ `false`
* `(set-eq? $set (new-set "a" "c" "b"))` ➜ `true`
* `(set-eq? $set (new-set "a" "c" "b" "e"))` ➜ `false`

## Forms

### `(set-eq? base:set other:set)` ➜ `bool`

This form returns true if both sets contain the same values.
