# append

`append` adds additional items to a vector or concatenates a string at the end
of another one, depending on the given arguments. See also this functions
cousin, [`prepend`](prepend.md).

`append` will append to a vector if the first argument coalesces to a vector;
otherwise string coalescing is attempted. This means depending on the coalescer
the behaviour of this function can change when the first argument is neither
vector nor string.

## Examples

* `(append ["foo"] "bar" "x" 3)` ➜ `["foo" "bar" "x" 3]`
* `(append "foo" "bar" "baz")` ➜ `"foobarbaz"`
* `(append "foo" 2)` ➜ `"foo2"` with humane coalescing, error otherwise
* `(append 2 3 4)` ➜ `"234"` with humane coalescing, error otherwise
* `(append null 1)` ➜ `[1]` because `null` can turn into empty vectors

## Forms

### `(append base:vector appends:any…)` ➜ `vector`

* `base` is an arbitrary expression.
* `appends` are one or more additional arbitrary expressions.

If the `base` coalesces to a vector, all further arguments will be appended to
the vector. Additional items in the vector can be of any type. The result is
a copy of the base vector with the newly added elements appended to it.

### `(append base:string appends:string…)` ➜ `string`

* `base` is an arbitrary expression.
* `appends` are one or more additional arbitrary expressions.

If `base` is a string, all further `appends` must also be strings. Each is added
to the base string without any separator.

## Context

`append` executes all expressions in their own contexts, so nothing is shared.
