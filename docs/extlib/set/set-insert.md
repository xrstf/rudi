# set-insert

This function adds the given values to the set, returning a new set. For
adding values in-place, using `set-insert!`. Just like when constructing a new
set with `new-set`, values must be either directly coalescable to strings, or
be vectors that contain only strings.

## Examples

All of the examples assume that `$set` is a set with `{"a", "b", "c"}`.

* `(set-insert $set "a")` ➜ `set{"a", "b", "c"}`
* `(set-insert $set "d")` ➜ `set{"a", "b", "c", "d"}`
* `(set-insert $set "x" ["y"])` ➜ `set{"a", "b", "c", "x", "y"}`

## Forms

### `(set-insert set:set value:any+)` ➜ `set`

This form returns a copy of the set, with all the values listed being added to
the set. Values that already exist in the set are not duplicated, of course.
