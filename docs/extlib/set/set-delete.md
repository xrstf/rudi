# set-delete

This function removes the given values from the set, returning a new set. For
removing values in-place, using `set-delete!`. Just like when constructing a new
set with `new-set`, values must be either directly coalescable to strings, or
be vectors that contain only strings.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-delete $set "a")` ➜ `set{"b", "c"}`
* `(set-delete $set "a" ["b" "d"] "")` ➜ `set{"c"}`

## Forms

### `(set-delete set:set value:any+)` ➜ `set`

This form returns a copy of the set, with all the values listed being removed
from the set. Values that do not occur in the set are ignored.
