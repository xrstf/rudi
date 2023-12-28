# set-has-any?

This function returns true if at least one of the given values occur in the
given set. See also [`set-has?`](set-has.md) for checking if _all_ of the given
values occur in the set.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-has-any? $set "")` ➜ `false`
* `(set-has-any? $set "b")` ➜ `true`
* `(set-has-any? $set "b" "c")` ➜ `true`
* `(set-has-any? $set "b" "d")` ➜ `true`
* `(set-has-any? $set "d" ["x" "y"])` ➜ `false`

## Forms

### `(set-has-any? base:set value:any+)` ➜ `bool`

This form returns true if the set contains at least one of the given values.
