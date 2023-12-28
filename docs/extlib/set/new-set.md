# new-set

This function returns a new string set containing all the given values. Values
are coalesced to strings, vectors are supported but only one level deep (see
examples).

## Examples

* `(new-set)` ➜ `set{}`
* `(new-set "test")` ➜ `set{"test"}`
* `(new-set "a" "b" "c" "b" "A")` ➜ `set{"a", "b", "c", "A"}`
* `(new-set "a" ["b" "c"] "d")` ➜ `set{"a", "b", "c", "d"}`
* `(new-set "a" ["b" "c" ["f"]] "d")` ➜ error

## Forms

### `(new-set)` ➜ `set`

This form returns a new, empty set.

### `(new-set value:any+)` ➜ `set`

This form coalesces all values as either string or vector. Vectors are unpacked
to one level deep (i.e. they can contain things that coalesce into a string, but
nothing else). Duplicate values can be given and will simply be dropped from the
set.
