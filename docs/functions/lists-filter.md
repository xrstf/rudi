# filter

`filter` returns a subset of the given vector/object, with elements filtered by
an expression. Only those items for which the given condition expression is
`true` will be in the resulting vector/object.

## Examples

* `(filter [0 "foo" ""] empty?)` -> `[0 ""]`
* `(filter [0 "foo" ""] [v] (not (empty? $v)))` -> `["foo"]`
* `(filter ["a" "b" "c" "d"] [idx v] (gt? $idx 1))` -> `["c" "d"]`

## Forms

### `(filter source func)`

* `source` is an arbitrary expression.
* `func` is a function identifier.

`filter` evaluates the source argument and coalesces it to a vector or object,
with vectors being preferred. If either of these operations fail, an error is
returned. Afterwards `filter` will apply the given `func` to each element
(for objects the function is applied to the values) and coalesces the result to
a boolean. If this yields `true`, the element is kept in the result, otherwise
it's discarded.

`func` must be a function that allows being called with exactly 1 argument. If
more arguments are needed, use the other form of `filter`.

### `(filter source namingvec expr)`

* `source` is an arbitrary expression.
* `naming` is a vector describing the desired loop variable name(s).
* `expr` is an arbitrary expression.

When evaluating more complex conditions, this form can be used. Instead of
anonymously applying a function to each argument, this form allows to set the
index/value (for vectors) or key/value (for objects) as variables, which can
then be used in arbitrary expressions, for example:

* `(filter .data [value] (not (empty? $value)))`
* `(filter .data [idx value] (gt? $idx 1))`

`namingvec` must be a vector containing one or two identifiers. If a single
identifier is given, it's the variable name for the value. If two identifiers
are given, the first is used for the index/key, the second is used for the value.

`expr` can then be any expression. Just like the other form, `source` is
evaluated and coalesced to vector/object and the expression is then applied to
each element, keeping only those for which the expression yields a boolish value.

## Context

`filter` evaluates all expressions using a shared context, so it's possible for
the filter functions to share variables.
