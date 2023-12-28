# delete

`delete` removes a key from an object or an item from a vector (shrinking that
vector by 1 element). It must be used on any expression that has a path
expression attached (e.g. `(delete (do-something).gone)`).

Note that `delete` is, like all functions in Rudi, stateless and so does not
modify the given argument directly. That is why `delete` can be used on
non-variables/documents, like in `(delete (read-config).password)` to get rid
of individual keys from a larger data structure. To enable this, `delete` always
returns the _remaining_ datastructure, not the value that was removed.

To make your changes "stick", use the bang modifier: `(delete! $var.foo)` – this
only makes sense for symbols (i.e. variables and the global document), since
there is no "source" to be updated for e.g. literals (hence
`(delete! [1 2 3][1])` is invalid).

[`set`](set.md) is another function that is most often used with the bang
modifier and also returns the entire datastructure, not just the newly inserted
value.

## Examples

* `(delete! $var.key)`
* `(delete! .[1])`

This function is mostly used with the bang modifier, but can sometimes also be
useful without if you only want to modify an object/vector "in transit":

* `(connect (delete (read-config).isAdmin))` to give `connect` a config object
  object without the `isAdmin` flag.

## Forms

### `(delete target:pathed)` ➜ `any`

* `target` is a any expression with a path expression.

`delete` evaluates the target expression and then removes whatever the path
expression is pointing to. The return value is the remaining data (i.e. not the
removed value).

When used with the bang modifier, `target` must be a symbol with a path
expression.
