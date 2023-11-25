# set

`set` is used to define variables and to update values of variables or the global
document. Similar to [delete](core-delete.md) it is most often used with the
bang modifier (`(set! …)`), as modifications "do not stick" otherwise.

`set` allows both overwriting the entire target value (`(set! $var …)` or
`(set! . …)`) as well as setting just a sub element (for example
`(set! .foo[0] "new-value")`).

Variables defined by `set` are scoped and only valid for all following sibling
expressions, never the parent. For for example `(if true (set! $x 42)) $x` is
invalid, as `$x` only exists in the positive branch of that `if` tuple.

`set` returns the value that was set.

## Examples

* `(set! $foo 42)` -> `42`
* `(set! $foo 42) $foo` -> `42`
* `(set! .global[0].document "new-value")` -> `"new-value"`

Without the bang modifier, `set` is less useful:

* `(set! $foo 42) (set $foo "new-value") $foo` -> `42`

## Forms

### `(set target value)`

* `target` is either a symbol, but can also be a vector/object/tuple with
  path expression, as long as the tuple would evaluate to something where the
  path expression fits.
* `value` is any expression.

`set` evaluates the value and then applies it to the `target`.

If `target` is a variable with no path expression, its value is simply overwritten.
Likewise, a `.` will make `set` overwrite the entire global document.

If a path expression is present, `set` will only set value deeply nested in the
target value (e.g. `(set $foo.bar 42)` will update the value at the key "bar"
in `$foo`, which is hopefully an object). Even in this case, the _return value_
of `set` is still the _set_ value (42 in the previous example), not the combined
value.

Also note that without the bang modifier, all of variable and document changes
are only valid inside the `set` tuple itself and will disappear / not leak
outside.

If the `value` or the `target` return an error while evaluating, the error is
returned.

## Context

`set` evaluates all expressions as their own scopes, so variables from the
`value` expression do not influence the `target` expression's evaluation.
