# new-key-set

This function takes an object and returns a new set that contains all the keys
from the object. Duplicate keys (for badly written object literals in Rudi code)
are allowed and will simply be merged into one within the set.

## Examples

* `(new-key-set "test")` ➜ `error`
* `(new-key-set {})` ➜ (empty set)
* `(new-key-set {a 1 b 2})` ➜ `set{"a", "b"}`
* `(new-key-set {a 1 b 2 b 3})` ➜ `set{"a", "b"}`

## Forms

### `(new-key-set obj:object)` ➜ `set`

This is the only form of this function. It coalesces the given value to an
object and then lists all the keys into a set. Returns an error if the given
value cannot be coalesced into an object.
