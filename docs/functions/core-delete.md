# delete

`delete` removes a key from an object or an item from a vector (shrinking that
vector by 1 element). It must be used with a path expression and can be used
on expressions that support path expressions (e.g. `(delete (do-something).gone)`).

Note that `delete` is, like all functions in Rudi, stateless and so does not
modify the given argument directly. To make your changes "stick", use the bang
modifier: `(delete! $var.foo)` – this only makes sense for symbols (i.e.
variables and the global document), which is incidentally also what the bang
modifier itself already enforces. So an expression like `(delete! [1 2][1])` is
invalid.

[set](core-set.md) is another function that is most often used with the bang
modifier.

## Examples

* `(delete! $var.key)`
* `(delete! .[1])`

This function is mostly used with the bang modifier, but can sometimes also be
useful without if you only want to modify an object/vector "in transit":

* `(handle (delete (read-config).isAdmin))` to give `handle` a config object
  object without the `isAdmin` flag.

## Forms

### `(delete target)`

* `target` is any expression that supports path expressions (symbols, tuples,
  vector nodes and object nodes)

`delete` evaluates the target expression and then removes whatever the path
expression is pointing to. The return value is the remaining data (i.e. not the
removed value).

When used with the bang modifier, `target` must be a symbol with a path
expression.

## Context

`delete` evaluates the expression in its own scope, so variables defined in it
do not leak.
