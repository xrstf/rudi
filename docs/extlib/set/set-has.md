# set-has?

This function returns true if all of the given values occur in the given set.
See also [`set-has-any?`](set-has-any.md) for checking if _at least one_ of the
given values occurs in the set.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-has? $set "")` ➜ `false`
* `(set-has? $set "b")` ➜ `true`
* `(set-has? $set "b" "c")` ➜ `true`
* `(set-has? $set "b" "d")` ➜ `false`
* `(set-has? $set "b" ["a" "c"])` ➜ `true`

## Forms

### `(set-has? base:set value:any+)` ➜ `bool`

This form returns true if the set contains all of the given values.
