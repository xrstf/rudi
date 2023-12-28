# set-size

This function returns the number of values in the given set.

## Examples

* `(set-size (new-set "a" "b"))` ➜ `2`
* `(set-size (set-delete (new-set "a") "a"))` ➜ `0`

## Forms

### `(set-size set:set)` ➜ `int`

This form returns the number of values in the set.
