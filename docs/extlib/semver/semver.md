# semver

This function will parse a string as a semantic version. The parsing is "relaxed",
allowing for a leading `"v"` and the least significant parts can be left out
when they are zero (e.g. `"v1.0.0"` is just as valid as `"1"`).

Parsed semvers are a custom type (not a string, not a vector). They can be
directly compared to each other and to strings (i.e. they can be coalesced to
a string, depending on the coalescer).

## Examples

* `(semver "v1.2")` ➜ semver object
* `(semver "foo")` ➜ error
* `(eq? (to-string (semver "v1.0")) "1.0.0")` ➜ `true`
* `(eq? (semver "v1.0") "1.0.0")` ➜ `true` (with human coalescing)
* `(gt? (semver "v1.0") (semver "v1.0.1"))` ➜ `false`
