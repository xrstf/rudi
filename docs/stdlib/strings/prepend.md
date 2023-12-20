# prepend

`prepend` adds additional items to a vector or concatenates a string at the end
of another one, depending on the given arguments. See also this functions
cousin, [`append`](append.md).

`prepend` will prepend to a vector if the first argument coalesces to a vector;
otherwise string coalescing is attempted. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

When prepending multiple elements, they are all prepended as a single list, not
individually, so `(prepend [1] 2 3)` yields `[2 3 1]`, not `[3 2 1]`.

## Examples

* `(prepend ["foo"] "bar" "x" 3)` ➜ `["x" 3 "foo" "bar"]`
* `(prepend "foo" "bar" "baz")` ➜ `"barbazfoo"`
* `(prepend "foo" 2)` ➜ `"2foo"` with humane coalescing, error otherwise
* `(prepend 2 3 4)` ➜ `"342"` with humane coalescing, error otherwise
* `(prepend null 1)` ➜ `[1]` because `null` can turn into empty vectors

## Forms

### `(prepend base:vector prepends:any…)` ➜ `vector`

* `base` is an arbitrary expression.
* `prepends` are one or more additional arbitrary expressions.

If the `base` coalesces to a vector, all further arguments will be prepended to
the vector. Additional items in the vector can be of any type. The result is
a copy of the base vector with the newly added elements prepended to it.

### `(prepend base:string prepends:string…)` ➜ `string`

* `base` is an arbitrary expression.
* `prepends` are one or more additional arbitrary expressions.

If `base` is a string, all further `prepends` must also be strings. Each is added
to the base string without any separator.

## Context

`prepend` executes all expressions in their own contexts, so nothing is shared.
