# map

`map` creates a copy of the given vector/object and then applies a
function/expression to every element in that copy.

## Examples

* `(map ["foo" "bar"] to-upper)` ➜ `["FOO" "BAR"]`
* `(map [1 2 3] [v] (+ $v 3))` ➜ `[4 5 6]`
* `(map ["a" "b" "c" "d"] [idx v] $idx)` ➜ `[0 1 2 3]`

## Forms

### `(map source:expression func:identifier)` ➜ `any`

* `source` is an arbitrary expression.
* `func` is a function identifier.

`map` evaluates the source argument and coalesces it to a vector or object,
with vectors being preferred. If either of these operations fail, an error is
returned. Afterwards `map` will replace the value of each item with the result
of applying the function `func` to it.

`func` must be a function that allows being called with exactly 1 argument. If
more arguments are needed, use the other form of `map`.

### `(map source:expression params:vector expr:expression)` ➜ `any`

* `source` is an arbitrary expression.
* `params` is a vector describing the desired loop variable name(s).
* `expr` is an arbitrary expression.

When evaluating more complex conditions, this form can be used. Instead of
anonymously applying a function to each argument, this form allows to set the
index/value (for vectors) or key/value (for objects) as variables, which can
then be used in arbitrary expressions, for example:

* `(map .data [v] (+ $v 3))`
* `(map .data [idx v] $idx)`

`params` must be a vector containing one or two identifiers. If a single
identifier is given, it's the variable name for the value. If two identifiers
are given, the first is used for the index/key, the second is used for the value.

`expr` can then be any expression. Just like the other form, `source` is
evaluated and coalesced to vector/object and the expression is then applied to
each element.

## Context

`map` evaluates all expressions using a shared context, so it's possible for
the map functions to share variables.
