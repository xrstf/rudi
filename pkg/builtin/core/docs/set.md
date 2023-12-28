# set

`set` is used to define variables and to update values of variables or the global
document. Similar to [`delete`](delete.md) it is most often used with the
bang modifier (`(set! …)`), as modifications "do not stick" otherwise. However
`set` can also be used to modify a datastructure "in transit", like in
`(set! .user.settings (set $defaultSettings.isAdmin false))`, where the inner
function will change a subfield (`isAdmin`) but still return the entire
default settings object.

`set` allows both overwriting the entire target value (`(set! $var …)` or
`(set! . …)`) as well as setting just a sub element (for example
`(set! .foo[0] "new-value")`).

Variables defined by `set` are scoped and only valid for all following sibling
expressions, never the parent. For for example `(if true (set! $x 42)) $x` is
invalid, as `$x` only exists in the positive branch of that `if` tuple.

`set` returns the entire target data structure, as if no path expression was
given. For example in `(set (read-config).isAdmin false)`, the entire
configuration would be returned, not just `false`. This is slightly different
semantics from most other functions, which would return only the resulting
value (e.g. `(append $foo.list 2)` would not return the entire `$foo` variable,
but only the `list` vector). In that sense, `set` works like `delete`, which
returns the _remaining_ data, not whatever was removed.

## Examples

* `(set! $foo 42)` ➜ `42`
* `(set! $foo 42) $foo` ➜ `42`
* `(set! .global[0].document "new-value")` ➜ `"new-value"`
* `(set! $config {a "yes" b "yes"}) (set $config.a "no")` ➜ `{a "yes" b "no"}`

## Forms

### `(set target:pathed value:any)` ➜ `any`

* `target` is any expression that can have a path expression.
* `value` is any expression.

`set` evaluates the `value` and then the path expression of `target`. If both
were successfully evaluated, the value is inserted into the target value at
the given path and then the entire target is returned.

Note that for variables, the path expression can be empty (e.g.
`(set $foo 42)`). For all other valid targets, a path expression must be set
(e.g. `(set (read-config).field 42)`) because there is no source that could be
overwritten (like with a variable or the global document).

Also note that without the bang modifier, all of variable and document changes
are only returned, the underlying value is not modified in-place.

`set!` can only be used with variables and bare path expressions (i.e. the
global document), because there is no logical way to modify the result of a
function call in-place.
