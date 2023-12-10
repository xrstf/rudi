# range

`range` evaluates an expression for each element of a vector/object, effectively
allowing to loop over them. Note that this is not intended to actually modify
the source vector/object, for those tasks [`map`](map.md) and
[`filter`](filter.md) should be used instead. `range` is less often used,
especially with over functions that would have side effects or to build up
counters/sums.

`range` returns the value of the last evalauted expression. Note that since
objects are iterated in effectively random order, one should not rely on the
return value of `range` when ranging over objects.

## Examples

* `(range ["foo" "bar"] [value] (print $value))` -> `nil`
* `(range {a "b" c "d"} [key value] (print $key))` -> `nil`

## Forms

<!-- ### `(range source func)`

* `source` is an arbitrary expression.
* `func` is a function identifier.

`func` must be a function that allows being called with exactly 1 argument. If
more arguments are needed, use the other form of `range`. -->

### `(range source namingvec expr)`

* `source` is an arbitrary expression.
* `naming` is a vector describing the desired loop variable name(s).
* `expr` is an arbitrary expression.

`range` evaluates the source argument and coalesces it to a vector or object,
with vectors being preferred. If either of these operations fail, an error is
returned. The naming vector `namingvec` then allows to set the index/value (for
vectors) or key/value (for objects) as variables, which can then be used in the
expression `expr`, for example:

* `(range .data [v] (set! .data.users[$v] = "foo"))`
* `(range .data [v] (set! $var (+ (try $var 0) 1)))`

`namingvec` must be a vector containing one or two identifiers. If a single
identifier is given, it's the variable name for the value. If two identifiers
are given, the first is used for the index/key, the second is used for the value.

`expr` can then be any expression. Just like the other form, `source` is
evaluated and coalesced to vector/object and the expression is then applied to
each element.

## Context

`range` evaluates all expressions using a shared context, so it's possible for
the expressions to share variables.
